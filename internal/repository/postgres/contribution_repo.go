package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type contributionRepo struct {
	q *dbgen.Queries
}

func NewContributionRepository(q *dbgen.Queries) domain.ContributionRepository {
	return &contributionRepo{q: q}
}

func (r *contributionRepo) Create(ctx context.Context, p domain.CreateContributionParams) (*domain.Contribution, error) {
	row, err := r.q.CreateContribution(ctx, dbgen.CreateContributionParams{
		UserID:     p.UserID,
		PoolID:     p.PoolID,
		Amount:     p.Amount,
		Currency:   p.Currency,
		StripePiID: p.StripePIID,
	})
	if err != nil {
		return nil, fmt.Errorf("contributionRepo.Create: %w", err)
	}
	return rowToContribution(row), nil
}

func (r *contributionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Contribution, error) {
	row, err := r.q.GetContributionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("contributionRepo.GetByID: %w", err)
	}
	return rowToContribution(row), nil
}

func (r *contributionRepo) GetByStripePI(ctx context.Context, piID string) (*domain.Contribution, error) {
	row, err := r.q.GetContributionByStripePI(ctx, &piID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("contributionRepo.GetByStripePI: %w", err)
	}
	return rowToContribution(row), nil
}
func (r *contributionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ContributionStatus, paymentRef *string) error {
	return r.q.UpdateContributionStatus(ctx, dbgen.UpdateContributionStatusParams{
		ID:         id,
		Status:     string(status),
		PaymentRef: paymentRef,
	})
}
func (r *contributionRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Contribution, int, error) {
	rows, err := r.q.ListContributionsByUser(ctx, dbgen.ListContributionsByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("contributionRepo.ListByUser: %w", err)
	}
	total, err := r.q.CountContributionsByUser(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("contributionRepo.ListByUser count: %w", err)
	}
	out := make([]*domain.Contribution, len(rows))
	for i, row := range rows {
		out[i] = rowToContribution(row)
	}
	return out, int(total), nil
}
func (r *contributionRepo) SumByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	sum, err := r.q.SumCompletedContributionsByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("contributionRepo.SumByUser: %w", err)
	}
	return sum, nil
}
func rowToContribution(row dbgen.Contribution) *domain.Contribution {
	return &domain.Contribution{
		ID:         row.ID,
		UserID:     row.UserID,
		PoolID:     row.PoolID,
		Amount:     row.Amount,
		Currency:   row.Currency,
		Status:     domain.ContributionStatus(row.Status),
		PaymentRef: row.PaymentRef,
		StripePIID: row.StripePiID,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}
