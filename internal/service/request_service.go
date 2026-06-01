package service

import (
	"context"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
)

// RequestService manages the full lifecycle of a member's item request —
// from submission through admin verification and into the queue.
type RequestService struct {
	requests domain.RequestRepository
	queue    *QueueService // injected so Verify can immediately enqueue
}

func NewRequestService(requests domain.RequestRepository, queue *QueueService) *RequestService {
	return &RequestService{requests: requests, queue: queue}
}

// ── Input types ───────────────────────────────────────────────────────────────

// SubmitInput is what a member sends when they submit a new request.
// Members submit directly — no draft-save step in the API for MVP.
type SubmitInput struct {
	UserID        uuid.UUID
	ItemCategory  string
	ItemName      string
	Description   string
	Urgency       domain.UrgencyLevel
	EstimatedCost float64
	Justification string
}

// UpdateRequestInput contains the fields a member may edit on their own request.
// Only allowed while status is 'draft' or 'submitted' (enforced by the SQL query).
type UpdateRequestInput struct {
	ItemCategory  *string
	ItemName      *string
	Description   *string
	Urgency       *domain.UrgencyLevel
	EstimatedCost *float64
	Justification *string
}

// ── Methods ───────────────────────────────────────────────────────────────────

// Submit creates a new request in 'submitted' status.
// The member fills the form and submits — we write it directly to 'submitted'
// so it joins the admin verification queue immediately.
func (s *RequestService) Submit(ctx context.Context, in SubmitInput) (*domain.Request, error) {
	if in.EstimatedCost <= 0 {
		return nil, fmt.Errorf("%w: estimated_cost must be positive", domain.ErrInvalidInput)
	}
	if in.ItemName == "" || in.Description == "" || in.Justification == "" {
		return nil, fmt.Errorf("%w: item_name, description, and justification are required", domain.ErrInvalidInput)
	}

	req, err := s.requests.Create(ctx, domain.CreateRequestParams{
		UserID:        in.UserID,
		ItemCategory:  in.ItemCategory,
		ItemName:      in.ItemName,
		Description:   in.Description,
		Urgency:       in.Urgency,
		EstimatedCost: in.EstimatedCost,
		Justification: in.Justification,
		Status:        domain.RequestSubmitted, // skip draft — submit directly
	})
	if err != nil {
		return nil, fmt.Errorf("requestService.Submit: %w", err)
	}
	return req, nil
}

// GetByID fetches a request by ID.
// Used by both members (their own requests) and admins (any request).
// Ownership checking is the handler's responsibility.
func (s *RequestService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Request, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("requestService.GetByID: %w", err)
	}
	return req, nil
}

// Update lets a member edit their request while it's still in draft or submitted.
// Once an admin has verified it, editing is locked (enforced at the DB query level:
// the UpdateRequest SQL has WHERE status IN ('draft', 'submitted')).
func (s *RequestService) Update(ctx context.Context, id uuid.UUID, in UpdateRequestInput) (*domain.Request, error) {
	if in.EstimatedCost != nil && *in.EstimatedCost <= 0 {
		return nil, fmt.Errorf("%w: estimated_cost must be positive", domain.ErrInvalidInput)
	}

	req, err := s.requests.Update(ctx, id, domain.UpdateRequestParams{
		ItemCategory:  in.ItemCategory,
		ItemName:      in.ItemName,
		Description:   in.Description,
		Urgency:       in.Urgency,
		EstimatedCost: in.EstimatedCost,
		Justification: in.Justification,
	})
	if err != nil {
		return nil, fmt.Errorf("requestService.Update: %w", err)
	}
	return req, nil
}

// Verify is an admin action. It:
//  1. Checks the request is in 'submitted' status
//  2. Moves it to 'verified'
//  3. Immediately enqueues it in the priority queue
//
// These two steps are not wrapped in a DB transaction intentionally:
// if the status update succeeds but the enqueue fails, the request sits at
// 'verified' and the queue worker / admin can retry enqueuing it manually.
// A partial failure here is recoverable; a rollback would lose the verification.
func (s *RequestService) Verify(ctx context.Context, id uuid.UUID) (*domain.Request, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("requestService.Verify: fetch: %w", err)
	}
	if req.Status != domain.RequestSubmitted {
		return nil, fmt.Errorf("%w: only submitted requests can be verified (got %s)", domain.ErrInvalidState, req.Status)
	}

	// Step 1 — update status.
	if err := s.requests.UpdateStatus(ctx, id, domain.RequestVerified, nil); err != nil {
		return nil, fmt.Errorf("requestService.Verify: update status: %w", err)
	}
	req.Status = domain.RequestVerified

	// Step 2 — enqueue (uses req.Urgency and req.UserID for initial scoring).
	if _, err := s.queue.Enqueue(ctx, req); err != nil {
		return nil, fmt.Errorf("requestService.Verify: enqueue: %w", err)
	}

	// Step 3 — mark as queued so the member sees the right status.
	if err := s.requests.UpdateStatus(ctx, id, domain.RequestQueued, nil); err != nil {
		return nil, fmt.Errorf("requestService.Verify: mark queued: %w", err)
	}
	req.Status = domain.RequestQueued

	return req, nil
}

// Reject is an admin action. It moves the request to 'rejected' and records
// the reason. The rejection note is shown to the member so they understand why.
// A rejected request cannot be re-submitted — the member must open a new one.
func (s *RequestService) Reject(ctx context.Context, id uuid.UUID, note string) (*domain.Request, error) {
	req, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("requestService.Reject: fetch: %w", err)
	}
	if req.Status != domain.RequestSubmitted && req.Status != domain.RequestVerified {
		return nil, fmt.Errorf("%w: cannot reject a request in %s status", domain.ErrInvalidState, req.Status)
	}

	wasQueued := req.Status == domain.RequestVerified || req.Status == domain.RequestQueued

	if err := s.requests.UpdateStatus(ctx, id, domain.RequestRejected, &note); err != nil {
		return nil, fmt.Errorf("requestService.Reject: %w", err)
	}
	req.Status = domain.RequestRejected
	req.RejectionNote = &note

	// If it was already in the queue (verified/queued), remove it.
	if wasQueued {
		_ = s.queue.Dequeue(ctx, id) // best-effort; no queue entry if not yet enqueued
	}

	return req, nil
}

// MyRequests returns a paginated list of requests belonging to the calling user.
func (s *RequestService) MyRequests(ctx context.Context, userID uuid.UUID, page pagination.Page) ([]*domain.Request, pagination.Page, error) {
	reqs, total, err := s.requests.List(ctx, domain.RequestFilter{
		UserID: &userID,
		Limit:  page.Limit,
		Offset: page.Offset(),
	})
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("requestService.MyRequests: %w", err)
	}
	return reqs, page.WithTotal(total), nil
}

// AdminList returns a paginated list of requests filtered by status.
// Pass nil status to get all requests.
func (s *RequestService) AdminList(ctx context.Context, status *domain.RequestStatus, page pagination.Page) ([]*domain.Request, pagination.Page, error) {
	reqs, total, err := s.requests.List(ctx, domain.RequestFilter{
		Status: status,
		Limit:  page.Limit,
		Offset: page.Offset(),
	})
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("requestService.AdminList: %w", err)
	}
	return reqs, page.WithTotal(total), nil
}

// Delete removes a request. Only succeeds if the request is in 'draft' status
// (enforced at the SQL level via DeleteDraftRequest).
// A member can only delete their own drafts — ownership checked in the handler.
func (s *RequestService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.requests.Delete(ctx, id); err != nil {
		return fmt.Errorf("requestService.Delete: %w", err)
	}
	return nil
}

// MarkFunded updates the request status to 'funded' when the pool has allocated
// enough funds. Called by FulfillmentService when procurement begins.
func (s *RequestService) MarkFunded(ctx context.Context, id uuid.UUID) error {
	if err := s.requests.UpdateStatus(ctx, id, domain.RequestFunded, nil); err != nil {
		return fmt.Errorf("requestService.MarkFunded: %w", err)
	}
	_ = s.queue.Dequeue(ctx, id) // remove from queue now it's funded
	return nil
}
