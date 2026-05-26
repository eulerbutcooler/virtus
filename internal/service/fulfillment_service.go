package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
)

type FulfillmentService struct {
	fulfillments domain.FulfillmentRepository
	requests     domain.RequestRepository
	pool         *PoolService
}

func NewFulfillmentService(fulfillments domain.FulfillmentRepository, requests domain.RequestRepository, pool *PoolService) *FulfillmentService {
	return &FulfillmentService{
		fulfillments: fulfillments,
		requests:     requests,
		pool:         pool,
	}
}

// This is what admin sends when setting vendor + cost info
type UpdateDetailsInput struct {
	VendorName *string
	VendorRef  *string
	ActualCost *float64
	Status     *domain.FulfillmentStatus
	Notes      *string
	ProcuredAt *time.Time
}

// Creates a fulfillment record for a funded request and moves the request to 'procuring' status.
// Called by admin when they start sourcing the item
func (s *FulfillmentService) Begin(ctx context.Context, requestID uuid.UUID) (*domain.Fulfillment, error) {
	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.Begin: fetch request: %w", err)
	}
	if req.Status != domain.RequestFunded {
		return nil, fmt.Errorf("%w: procurement can only begin on a funded request (got %s)", domain.ErrInvalidState, req.Status)
	}
	existing, err := s.fulfillments.GetByRequestID(ctx, requestID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("fulfillmentService.Begin: check existing: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: fulfillment already exists for this request", domain.ErrConflict)
	}
	fulfillment, err := s.fulfillments.Create(ctx, domain.CreateFulfillmentParams{
		RequestID: requestID,
	})
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.Begin: create: %w", err)
	}

	if err := s.requests.UpdateStatus(ctx, requestID, domain.RequestProcuring, nil); err != nil {
		return nil, fmt.Errorf("fulfillmentService.Begin: update request status: %w", err)
	}
	return fulfillment, nil
}

// Returns a fulfillment by its own ID
func (s *FulfillmentService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Fulfillment, error) {
	f, err := s.fulfillments.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.GetByID: %w", err)
	}
	return f, nil
}

// Returns a fulfillment for a given request
// Used by the member dashboard to show procurement status
func (s *FulfillmentService) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.Fulfillment, error) {
	f, err := s.fulfillments.GetByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.GetByRequestID: %w", err)
	}
	return f, nil
}

// Lets admin update vendor info, actual cost, and status as procurement progresses
// Uses COALESCE in SQL so only non-nil fields are changed
func (s *FulfillmentService) UpdateDetails(ctx context.Context, id uuid.UUID, in UpdateDetailsInput) (*domain.Fulfillment, error) {
	f, err := s.fulfillments.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.UpdateDetails: fetch: %w", err)
	}
	if f.ProcurementStatus == domain.FulfillmentDelivered ||
		f.ProcurementStatus == domain.FulfillmentCancelled {
		return nil, fmt.Errorf("%w: cannot update a %s fulfillment",
			domain.ErrInvalidState, f.ProcurementStatus)
	}
	updated, err := s.fulfillments.Update(ctx, id, domain.UpdateFulfillmentParams{
		VendorName: in.VendorName,
		VendorRef:  in.VendorRef,
		ActualCost: in.ActualCost,
		Status:     in.Status,
		Notes:      in.Notes,
		ProcuredAt: in.ProcuredAt,
	})
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.UpdateDetails: %w", err)
	}

	if in.ActualCost != nil && f.ActualCost == nil {
		if err := s.pool.Debit(ctx, domain.GlobalPoolID, *in.ActualCost); err != nil {
			_ = err
		}
	}
	return updated, nil
}

// Sets the fulfillment status to 'delivered' and marks the associated request as 'delivered'
// Called by DeliveryService after a delivery is verified
func (s *FulfillmentService) MarkDelivered(ctx context.Context, fulfillmentID uuid.UUID) error {
	f, err := s.fulfillments.GetByID(ctx, fulfillmentID)
	if err != nil {
		return fmt.Errorf("fulfillmentService.MarkDelivered: fetch: %w", err)
	}
	if f.ProcurementStatus == domain.FulfillmentDelivered {
		return nil
	}
	delivered := domain.FulfillmentDelivered
	now := time.Now()
	if _, err := s.fulfillments.Update(ctx, fulfillmentID, domain.UpdateFulfillmentParams{
		Status:     &delivered,
		ProcuredAt: &now,
	}); err != nil {
		return fmt.Errorf("fulfillmentService.MarkDelivered: update fulfillment: %w", err)
	}
	// Cascade to request status.
	if err := s.requests.UpdateStatus(ctx, f.RequestID, domain.RequestDelivered, nil); err != nil {
		return fmt.Errorf("fulfillmentService.MarkDelivered: update request: %w", err)
	}
	return nil
}

// Cancels a fulfillment. Only allowed while still pending or vendor_selected.
func (s *FulfillmentService) Cancel(ctx context.Context, id uuid.UUID) (*domain.Fulfillment, error) {
	f, err := s.fulfillments.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.Cancel: fetch: %w", err)
	}
	cancellable := map[domain.FulfillmentStatus]bool{
		domain.FulfillmentPending:        true,
		domain.FulfillmentVendorSelected: true,
	}
	if !cancellable[f.ProcurementStatus] {
		return nil, fmt.Errorf("%w: cannot cancel a %s fulfillment",
			domain.ErrInvalidState, f.ProcurementStatus)
	}
	cancelled := domain.FulfillmentCancelled
	updated, err := s.fulfillments.Update(ctx, id, domain.UpdateFulfillmentParams{
		Status: &cancelled,
	})
	if err != nil {
		return nil, fmt.Errorf("fulfillmentService.Cancel: %w", err)
	}
	if f.ActualCost != nil {
		if err := s.pool.Credit(ctx, domain.GlobalPoolID, *f.ActualCost); err != nil {
			_ = err
		}
	}
	if err := s.requests.UpdateStatus(ctx, f.RequestID, domain.RequestFunded, nil); err != nil {
		return nil, fmt.Errorf("fulfillmentService.Cancel: revert request status: %w", err)
	}
	return updated, nil
}

// Returns a paginated list of all fulfillments. Admin only.
func (s *FulfillmentService) List(ctx context.Context, page pagination.Page) ([]*domain.Fulfillment, pagination.Page, error) {
	items, total, err := s.fulfillments.List(ctx, page.Limit, page.Offset())
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("fulfillmentService.List: %w", err)
	}
	return items, page.WithTotal(total), nil
}
