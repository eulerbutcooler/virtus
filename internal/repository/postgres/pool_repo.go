package postgres

import (
	"context"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	"github.com/google/uuid"
)

type poolRepo struct {
	q *dbgen.Queries
}

func NewPoolRepository(q *dbgen.Queries) domain.PoolRepository {
	return &poolRepo{q: q}
}

func (r *poolRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Pool, error) {
	row, err := r.q.GetPool(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("poolRepo.Get: %w", err)
	}
	return &domain.Pool{
		ID:        row.ID,
		Balance:   row.Balance,
		TotalIn:   row.TotalIn,
		TotalOut:  row.TotalOut,
		Currency:  row.Currency,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *poolRepo) Credit(ctx context.Context, poolID uuid.UUID, amount float64) error {
	if err := r.q.CreditPool(ctx, dbgen.CreditPoolParams{
		ID:     poolID,
		Amount: amount,
	}); err != nil {
		return fmt.Errorf("poolRepo.Credit: %w", err)
	}
	return nil
}

func (r *poolRepo) Debit(ctx context.Context, poolID uuid.UUID, amount float64) error {
	n, err := r.q.DebitPool(ctx, dbgen.DebitPoolParams{
		ID:     poolID,
		Amount: amount,
	})
	if err != nil {
		return fmt.Errorf("poolRepo.Debit: %w", err)
	}
	if n == 0 {
		return domain.ErrInsufficientFunds
	}
	return nil
}
