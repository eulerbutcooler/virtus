package service

import (
	"context"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
)

// Manages the fulfillment queue - enqueuing requests and exposing the public queue listing.
type QueueService struct {
	queue domain.QueueRepository
}

func NewQueueService(queue domain.QueueRepository) *QueueService {
	return &QueueService{
		queue: queue}
}

// Enqueue places a verified request at the back of the queue
func (s *QueueService) Enqueue(ctx context.Context, req *domain.Request) (*domain.QueueEntry, error) {
	entry, err := s.queue.Enqueue(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("queueService.Enqueue: %w", err)
	}
	return entry, nil
}

// Returns the queue entry for a request
func (s *QueueService) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.QueueEntry, error) {
	entry, err := s.queue.GetByRequestID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("queueService.GetByRequestID: %w", err)
	}
	return entry, nil
}

// Returns just the queue position number for a request.
func (s *QueueService) GetPosition(ctx context.Context, requestID uuid.UUID) (int, error) {
	pos, err := s.queue.GetPosition(ctx, requestID)
	if err != nil {
		return 0, fmt.Errorf("queueService.GetPosition: %w", err)
	}
	return pos, nil
}

// Returns a paginated public view of the queue.
func (s *QueueService) ListAll(ctx context.Context, page pagination.Page) ([]*domain.QueueEntry, pagination.Page, error) {
	entries, total, err := s.queue.ListAll(ctx, page.Limit, page.Offset())
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("queueService.ListAll: %w", err)
	}
	return entries, page.WithTotal(total), nil
}

// Updates how much of the pool has been allocated toward a request.
// When progress reaches 100% the fulfillment service marks it funded.
func (s *QueueService) UpdateFundingProgress(ctx context.Context, requestID uuid.UUID, progress float64) error {
	entry, err := s.queue.GetByRequestID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("queueService.UpdateFundingProgress: %w", err)
	}
	if err := s.queue.UpdateFunding(ctx, entry.ID, progress); err != nil {
		return fmt.Errorf("queueService.UpdateFundingProgress: %w", err)
	}
	return nil
}

// Removes a request from the queue.
// Called when a request is funded and moves to procuring, or is rejected/cancelled.
func (s *QueueService) Dequeue(ctx context.Context, requestID uuid.UUID) error {
	if err := s.queue.Dequeue(ctx, requestID); err != nil {
		return fmt.Errorf("queueService.Dequeue: %w", err)
	}
	return nil
}

// ListAllOrdered returns all queue entries sorted by position ascending.
// Used by the worker's fund-queue task.
func (s *QueueService) ListAllOrdered(ctx context.Context) ([]*domain.QueueEntry, error) {
	entries, err := s.queue.ListAllOrdered(ctx)
	if err != nil {
		return nil, fmt.Errorf("queueService.ListAllOrdered: %w", err)
	}
	return entries, nil
}

func (s *QueueService) RecalculateAll(ctx context.Context) error {
	entries, err := s.queue.ListAllOrdered(ctx)
	if err != nil {
		return fmt.Errorf("queueService.RecalculateAll: %w", err)
	}
	for i, e := range entries {
		newPos := i + 1
		if e.Position == newPos {
			continue
		}
		if err := s.queue.UpdatePosition(ctx, e.ID, newPos); err != nil {
			return fmt.Errorf("queueService.RecalculateAll: update position %d: %w", newPos, err)
		}
	}
	return nil
}
