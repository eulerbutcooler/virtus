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

type impactRepo struct {
	q *dbgen.Queries
}

func NewImpactRepository(q *dbgen.Queries) domain.ImpactRepository {
	return &impactRepo{q: q}
}

func (r *impactRepo) Create(ctx context.Context, p domain.CreateImpactParams) (*domain.ImpactRecord, error) {
	var score *int16
	if p.SatisfactionScore != nil {
		s := int16(*p.SatisfactionScore)
		score = &s
	}

	row, err := r.q.CreateImpactRecord(ctx, dbgen.CreateImpactRecordParams{
		DeliveryID:         p.DeliveryID,
		UserID:             p.UserID,
		IntervalLabel:      p.IntervalLabel,
		OutcomeDescription: p.OutcomeDescription,
		SatisfactionScore:  score,
		Metrics:            p.Metrics,
	})
	if err != nil {
		return nil, fmt.Errorf("impactRepo.Create: %w", err)
	}
	return rowToImpact(row), nil
}

func (r *impactRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ImpactRecord, error) {
	row, err := r.q.GetImpactRecordByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("impactRepo.GetByID: %w", err)
	}
	return rowToImpact(row), nil
}

func (r *impactRepo) ListByDelivery(ctx context.Context, deliveryID uuid.UUID) ([]*domain.ImpactRecord, error) {
	rows, err := r.q.ListImpactByDelivery(ctx, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("impactRepo.ListByDelivery: %w", err)
	}
	out := make([]*domain.ImpactRecord, len(rows))
	for i, row := range rows {
		out[i] = rowToImpact(row)
	}
	return out, nil
}

func (r *impactRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ImpactRecord, int, error) {
	rows, err := r.q.ListImpactByUser(ctx, dbgen.ListImpactByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("impactRepo.ListByUser: %w", err)
	}
	total, err := r.q.CountImpactByUser(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("impactRepo.ListByUser count: %w", err)
	}
	out := make([]*domain.ImpactRecord, len(rows))
	for i, row := range rows {
		out[i] = rowToImpact(row)
	}
	return out, int(total), nil
}

func rowToImpact(row dbgen.ImpactRecord) *domain.ImpactRecord {
	rec := &domain.ImpactRecord{
		ID:                 row.ID,
		DeliveryID:         row.DeliveryID,
		UserID:             row.UserID,
		IntervalLabel:      row.IntervalLabel,
		OutcomeDescription: row.OutcomeDescription,
		Metrics:            row.Metrics,
		RecordedAt:         row.RecordedAt.Time,
	}
	if row.SatisfactionScore != nil {
		s := int(*row.SatisfactionScore)
		rec.SatisfactionScore = &s
	}
	return rec
}
