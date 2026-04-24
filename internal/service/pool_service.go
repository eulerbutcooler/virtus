package service

import (
	"context"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/google/uuid"
)

type PoolService struct {
	pools domain.PoolRepository
}

func NewPoolService(pools domain.PoolRepository) *PoolService {
	return &PoolService{pools: pools}
}

// Returns current state of the global pools
// Will be shown on the dashboard
func (s *PoolService) GetStatus(ctx context.Context) (*domain.Pool, error) {
	pool, err := s.pools.Get(ctx, domain.GlobalPoolID)
	if err != nil {
		return nil, fmt.Errorf("poolService.GetStatus: %w", err)
	}
	return pool, nil
}

// Adds to the pool
// Called internally by ContributionService when a payment completes.
// Not exposed to an HTTP endpoint
func (s *PoolService) Credit(ctx context.Context, poolID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("%w: credit amount must be positive", domain.ErrInvalidInput)
	}
	if err := s.pools.Credit(ctx, poolID, amount); err != nil {
		return fmt.Errorf("poolService.Credit: %w", err)
	}
	return nil
}

// Removes funds from the poolID
// Called internally by FulfillmentService when a procurement is finalized
// Returns ErrInsufficientFunds if pool balance would go negative
func (s *PoolService) Debit(ctx context.Context, poolID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("%w: debit amount must be positive", domain.ErrInvalidInput)
	}
	if err := s.pools.Debit(ctx, poolID, amount); err != nil {
		return fmt.Errorf("poolService.Debit: %w", err)
	}
	return nil
}
