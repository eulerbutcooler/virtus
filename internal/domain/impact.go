package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ImpactRecord struct {
	ID                 uuid.UUID       `json:"id"`
	DeliveryID         uuid.UUID       `json:"delivery_id"`
	UserID             uuid.UUID       `json:"user_id"`
	IntervalLabel      string          `json:"interval_label"`
	OutcomeDescription *string         `json:"outcome_description,omitempty"`
	SatisfactionScore  *int            `json:"satisfaction_score,omitempty"`
	Metrics            json.RawMessage `json:"metrics"`
	RecordedAt         time.Time       `json:"recorded_at"`
}

type CreateImpactParams struct {
	DeliveryID         uuid.UUID
	UserID             uuid.UUID
	IntervalLabel      string
	OutcomeDescription *string
	SatisfactionScore  *int
	Metrics            json.RawMessage
}

type ImpactRepository interface {
	Create(ctx context.Context, p CreateImpactParams) (*ImpactRecord, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ImpactRecord, error)
	ListByDelivery(ctx context.Context, deliveryID uuid.UUID) ([]*ImpactRecord, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ImpactRecord, int, error)
}
