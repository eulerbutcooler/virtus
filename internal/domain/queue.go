package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type QueueEntry struct {
	ID                   uuid.UUID  `json:"id"`
	RequestID            uuid.UUID  `json:"request_id"`
	Position             int        `json:"position"`
	FundingProgress      float64    `json:"funding_progress"`
	EstimatedFulfillment *time.Time `json:"estimated_fulfillment,omitempty"`
	EnteredAt            time.Time  `json:"entered_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type QueueRepository interface {
	Enqueue(ctx context.Context, requestID uuid.UUID) (*QueueEntry, error)
	GetByRequestID(ctx context.Context, requestID uuid.UUID) (*QueueEntry, error)
	GetPosition(ctx context.Context, requestID uuid.UUID) (int, error)
	UpdateFunding(ctx context.Context, id uuid.UUID, progress float64) error
	UpdatePosition(ctx context.Context, id uuid.UUID, position int) error
	Dequeue(ctx context.Context, requestID uuid.UUID) error
	ListAll(ctx context.Context, limit, offset int) ([]*QueueEntry, int, error)
	ListAllOrdered(ctx context.Context) ([]*QueueEntry, error)
}
