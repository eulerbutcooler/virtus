package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ContributionStatus string

const (
	ContributionPending   ContributionStatus = "pending"
	ContributionCompleted ContributionStatus = "completed"
	ContributionFailed    ContributionStatus = "failed"
	ContributionRefunded  ContributionStatus = "refunded"
)

type Contribution struct {
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	PoolID     uuid.UUID          `json:"pool_id"`
	Amount     float64            `json:"amount"`
	Currency   string             `json:"currency"`
	Status     ContributionStatus `json:"status"`
	PaymentRef *string            `json:"payment_ref,omitempty"`
	StripePIID *string            `json:"stripe_pi_id,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type CreateContributionParams struct {
	UserID     uuid.UUID
	PoolID     uuid.UUID
	Amount     float64
	Currency   string
	StripePIID *string
}

type ContributionRepository interface {
	Create(ctx context.Context, p CreateContributionParams) (*Contribution, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Contribution, error)
	GetByStripePI(ctx context.Context, piID string) (*Contribution, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status ContributionStatus, paymentRef *string) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Contribution, int, error)
	SumByUser(ctx context.Context, userID uuid.UUID) (float64, error)
}
