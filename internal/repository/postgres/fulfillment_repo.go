package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type fulfillmentRepo struct {
	q *dbgen.Queries
}

func NewFulfillmentRepository(q *dbgen.Queries) domain.FulfillmentRepository {
	return &fulfillmentRepo{q: q}
}

func (r *fulfillmentRepo) Create(ctx context.Context, p domain.CreateFulfillmentParams) (*domain.Fulfillment, error) {
	row, err := r.q.CreateFulfillment(ctx, p.RequestID)
	if err != nil {
		return nil, fmt.Errorf("fulfillmentRepo.Create: %w", err)
	}
	return rowToFulfillment(row), nil
}

func (r *fulfillmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Fulfillment, error) {
	row, err := r.q.GetFulfillmentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("fulfillmentRepo.GetByID: %w", err)
	}
	return rowToFulfillment(row), nil
}

func (r *fulfillmentRepo) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.Fulfillment, error) {
	row, err := r.q.GetFulfillmentByRequestID(ctx, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("fulfillmentRepo.GetByRequestID: %w", err)
	}
	return rowToFulfillment(row), nil
}

func (r *fulfillmentRepo) Update(ctx context.Context, id uuid.UUID, p domain.UpdateFulfillmentParams) (*domain.Fulfillment, error) {
	var actualCost pgtype.Numeric
	if p.ActualCost != nil {
		if err := actualCost.Scan(*p.ActualCost); err != nil {
			return nil, fmt.Errorf("fulfillmentRepo.Update: scanning actual_cost: %w", err)
		}
	}

	var status dbgen.NullFulfillmentStatus
	if p.Status != nil {
		status = dbgen.NullFulfillmentStatus{
			FulfillmentStatus: dbgen.FulfillmentStatus(*p.Status),
			Valid:             true,
		}
	}

	var procuredAt pgtype.Timestamptz
	if p.ProcuredAt != nil {
		procuredAt = pgtype.Timestamptz{Time: *p.ProcuredAt, Valid: true}
	}

	row, err := r.q.UpdateFulfillment(ctx, dbgen.UpdateFulfillmentParams{
		ID:                id,
		VendorName:        p.VendorName,
		VendorRef:         p.VendorRef,
		ActualCost:        actualCost,
		ProcurementStatus: status,
		Notes:             p.Notes,
		ProcuredAt:        procuredAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("fulfillmentRepo.Update: %w", err)
	}
	return rowToFulfillment(row), nil
}

func (r *fulfillmentRepo) List(ctx context.Context, limit, offset int) ([]*domain.Fulfillment, int, error) {
	rows, err := r.q.ListFulfillments(ctx, dbgen.ListFulfillmentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fulfillmentRepo.List: %w", err)
	}
	total, err := r.q.CountFulfillments(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("fulfillmentRepo.List count: %w", err)
	}
	out := make([]*domain.Fulfillment, len(rows))
	for i, row := range rows {
		out[i] = rowToFulfillment(row)
	}
	return out, int(total), nil
}

func rowToFulfillment(row dbgen.Fulfillment) *domain.Fulfillment {
	f := &domain.Fulfillment{
		ID:                row.ID,
		RequestID:         row.RequestID,
		VendorName:        row.VendorName,
		VendorRef:         row.VendorRef,
		ProcurementStatus: domain.FulfillmentStatus(row.ProcurementStatus),
		Notes:             row.Notes,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
	}
	if row.ActualCost.Valid {
		var v float64
		if err := row.ActualCost.Scan(&v); err == nil {
			f.ActualCost = &v
		}
	}
	if row.ProcuredAt.Valid {
		t := row.ProcuredAt.Time
		f.ProcuredAt = &t
	}
	return f
}
