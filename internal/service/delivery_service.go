package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/google/uuid"
)

type DeliveryService struct {
	deliveries   domain.DeliveryRepository
	fulfillments *FulfillmentService
}

func NewDeliveryService(
	deliveries domain.DeliveryRepository,
	fulfillments *FulfillmentService,
) *DeliveryService {
	return &DeliveryService{
		deliveries:   deliveries,
		fulfillments: fulfillments,
	}
}

// This is what admin provides when marking an item as shipped.
type ShipInput struct {
	FulfillmentID  uuid.UUID
	TrackingNumber *string
	Carrier        *string
}

// Creates a delivery record when an item has been dispatched.
// This represents the moment the item leaves the vendor i.e. the fulfillment
// must be in 'ordered' or 'shipped' status (in case it was already created).
func (s *DeliveryService) Ship(ctx context.Context, in ShipInput) (*domain.Delivery, error) {
	f, err := s.fulfillments.GetByID(ctx, in.FulfillmentID)
	if err != nil {
		return nil, fmt.Errorf("deliveryService.Ship: fetch fulfillment: %w", err)
	}

	allowed := map[domain.FulfillmentStatus]bool{
		domain.FulfillmentOrdered: true,
		domain.FulfillmentShipped: true,
	}
	if !allowed[f.ProcurementStatus] {
		return nil, fmt.Errorf("%w: can only create delivery for an ordered/shipped fulfillment (got %s)",
			domain.ErrInvalidState, f.ProcurementStatus)
	}

	existing, err := s.deliveries.GetByFulfillmentID(ctx, in.FulfillmentID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("deliveryService.Ship: check existing: %w", err)
	}
	if existing != nil {
		return existing, nil
	}

	delivery, err := s.deliveries.Create(ctx, domain.CreateDeliveryParams{
		FulfillmentID:  in.FulfillmentID,
		TrackingNumber: in.TrackingNumber,
		Carrier:        in.Carrier,
	})
	if err != nil {
		return nil, fmt.Errorf("deliveryService.Ship: create: %w", err)
	}

	shipped := domain.FulfillmentShipped
	if _, err := s.fulfillments.UpdateDetails(ctx, in.FulfillmentID, UpdateDetailsInput{
		Status: &shipped,
	}); err != nil {
		return nil, fmt.Errorf("deliveryService.Ship: update fulfillment status: %w", err)
	}

	return delivery, nil
}

// Returns a delivery record by its own ID.
func (s *DeliveryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error) {
	d, err := s.deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("deliveryService.GetByID: %w", err)
	}
	return d, nil
}

// Returns the delivery record for a fulfillment.
// Used by the member dashboard and admin panel.
func (s *DeliveryService) GetByFulfillmentID(ctx context.Context, fulfillmentID uuid.UUID) (*domain.Delivery, error) {
	d, err := s.deliveries.GetByFulfillmentID(ctx, fulfillmentID)
	if err != nil {
		return nil, fmt.Errorf("deliveryService.GetByFulfillmentID: %w", err)
	}
	return d, nil
}

// This is what admin/verifier provides when confirming delivery.
type VerifyInput struct {
	ProofPhotoURL string
	DeliveredAt   time.Time
	VerifiedBy    uuid.UUID
}

// Verifies if the item has reached the member.
// Uses photo proof url to verify.
func (s *DeliveryService) Verify(ctx context.Context, id uuid.UUID, in VerifyInput) (*domain.Delivery, error) {
	if in.ProofPhotoURL == "" {
		return nil, fmt.Errorf("%w: proof_photo_url is required to verify a delivery", domain.ErrInvalidInput)
	}

	delivery, err := s.deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("deliveryService.Verify: fetch: %w", err)
	}

	if delivery.Status == domain.DeliveryDelivered {
		return delivery, nil
	}
	if delivery.Status != domain.DeliveryInTransit {
		return nil, fmt.Errorf("%w: can only verify an in_transit delivery (got %s)",
			domain.ErrInvalidState, delivery.Status)
	}

	verified, err := s.deliveries.Verify(ctx, id, domain.VerifyDeliveryParams{
		ProofPhotoURL: in.ProofPhotoURL,
		DeliveredAt:   in.DeliveredAt,
		VerifiedBy:    in.VerifiedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("deliveryService.Verify: persist: %w", err)
	}

	if err := s.fulfillments.MarkDelivered(ctx, delivery.FulfillmentID); err != nil {
		return nil, fmt.Errorf("deliveryService.Verify: cascade fulfillment: %w", err)
	}

	return verified, nil
}

// Records a failed delivery (lost in transit, returned, etc.)
// and moves the fulfillment back to 'ordered' so admin can retry shipment.
func (s *DeliveryService) MarkFailed(ctx context.Context, id uuid.UUID) error {
	delivery, err := s.deliveries.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deliveryService.MarkFailed: fetch: %w", err)
	}
	if delivery.Status != domain.DeliveryInTransit {
		return fmt.Errorf("%w: can only fail an in_transit delivery (got %s)",
			domain.ErrInvalidState, delivery.Status)
	}

	if err := s.deliveries.UpdateStatus(ctx, id, domain.DeliveryFailed); err != nil {
		return fmt.Errorf("deliveryService.MarkFailed: update delivery: %w", err)
	}

	ordered := domain.FulfillmentOrdered
	if _, err := s.fulfillments.UpdateDetails(ctx, delivery.FulfillmentID, UpdateDetailsInput{
		Status: &ordered,
	}); err != nil {
		return fmt.Errorf("deliveryService.MarkFailed: revert fulfillment: %w", err)
	}

	return nil
}
