package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var GlobalPoolID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type Pool struct {
	ID        uuid.UUID `json:"id"`
	Balance   float64   `json:"balance"`
	TotalIn   float64   `json:"total_in"`
	TotalOut  float64   `json:"total_out"`
	Currency  string    `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PoolRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*Pool, error)
	Credit(ctx context.Context, poolID uuid.UUID, amount float64) error
	Debit(ctx context.Context, poolID uuid.UUID, amount float64) error
}
