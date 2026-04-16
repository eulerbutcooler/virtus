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
	PriorityScore        float64    `json:"priority_score"`
	FundingProgress      float64    `json:"funding_progress"`
	EstimatedFulfillment *time.Time `json:"estimated_fulfillment,omitempty"`
	EnteredAt            time.Time  `json:"entered_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// Configurable weights for the priority formula.
type QueueScoreWeights struct {
	Urgency      float64
	WaitTime     float64
	Contribution float64
	Community    float64
}

var DefaultWeights = QueueScoreWeights{
	Urgency:      0.35,
	WaitTime:     0.30,
	Contribution: 0.20,
	Community:    0.15,
}

type QueueRepository interface {
	Enqueue(ctx context.Context, requestID uuid.UUID, score float64) (*QueueEntry, error)
	GetByRequestID(ctx context.Context, requestID uuid.UUID) (*QueueEntry, error)
	GetPosition(ctx context.Context, requestID uuid.UUID) (int, error)
	UpdateScore(ctx context.Context, id uuid.UUID, score float64, position int) error
	UpdateFunding(ctx context.Context, id uuid.UUID, progress float64) error
	Dequeue(ctx context.Context, requestID uuid.UUID) error
	ListAll(ctx context.Context, limit, offset int) ([]*QueueEntry, int, error)
}
