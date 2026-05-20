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

type queueRepo struct {
	q *dbgen.Queries
}

func NewQueueRepository(q *dbgen.Queries) domain.QueueRepository {
	return &queueRepo{q: q}
}

// Inserts a request into the queue at the next available position
func (r *queueRepo) Enqueue(ctx context.Context, requestID uuid.UUID) (*domain.QueueEntry, error) {
	maxPos, err := r.q.MaxQueuePosition(ctx)
	if err != nil {
		return nil, fmt.Errorf("queueRepo.Enqueue maxPos: %w", err)
	}
	row, err := r.q.EnqueueRequest(ctx, dbgen.EnqueueRequestParams{
		RequestID: requestID,
		Position:  maxPos + 1,
	})
	if err != nil {
		return nil, fmt.Errorf("queueRepo.Enqueue: %w", err)
	}
	return rowToQueueEntry(row), nil
}

func (r *queueRepo) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.QueueEntry, error) {
	row, err := r.q.GetQueueEntryByRequestID(ctx, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("queueRepo.GetByRequestID: %w", err)
	}
	return rowToQueueEntry(row), nil
}

func (r *queueRepo) GetPosition(ctx context.Context, requestID uuid.UUID) (int, error) {
	pos, err := r.q.GetQueuePosition(ctx, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, fmt.Errorf("queueRepo.GetPosition: %w", err)
	}
	return int(pos), nil
}

func (r *queueRepo) UpdateFunding(ctx context.Context, id uuid.UUID, progress float64) error {
	return r.q.UpdateQueueFunding(ctx, dbgen.UpdateQueueFundingParams{
		ID:              id,
		FundingProgress: progress,
	})
}

func (r *queueRepo) UpdatePosition(ctx context.Context, id uuid.UUID, position int) error {
	return r.q.UpdateQueuePosition(ctx, dbgen.UpdateQueuePositionParams{
		ID:       id,
		Position: int32(position),
	})
}

func (r *queueRepo) Dequeue(ctx context.Context, requestID uuid.UUID) error {
	return r.q.DequeueRequest(ctx, requestID)
}

func (r *queueRepo) ListAll(ctx context.Context, limit, offset int) ([]*domain.QueueEntry, int, error) {
	rows, err := r.q.ListQueueEntries(ctx, dbgen.ListQueueEntriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("queueRepo.ListAll: %w", err)
	}
	total, _ := r.q.CountQueueEntries(ctx)
	out := make([]*domain.QueueEntry, len(rows))
	for i, row := range rows {
		out[i] = rowToQueueEntry(row)
	}
	return out, int(total), nil
}

func (r *queueRepo) ListAllOrdered(ctx context.Context) ([]*domain.QueueEntry, error) {
	rows, err := r.q.ListAllQueueEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("queueRepo.ListAllOrdered: %w", err)
	}
	out := make([]*domain.QueueEntry, len(rows))
	for i, row := range rows {
		out[i] = rowToQueueEntry(row)
	}
	return out, nil
}

func rowToQueueEntry(row dbgen.QueueEntry) *domain.QueueEntry {
	e := &domain.QueueEntry{
		ID:              row.ID,
		RequestID:       row.RequestID,
		Position:        int(row.Position),
		FundingProgress: row.FundingProgress,
		EnteredAt:       row.EnteredAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
	if row.EstimatedFulfillment.Valid {
		t := row.EstimatedFulfillment.Time
		e.EstimatedFulfillment = &t
	}
	return e
}
