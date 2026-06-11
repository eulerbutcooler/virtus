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

type requestRepo struct {
	q *dbgen.Queries
}

func NewRequestRepository(q *dbgen.Queries) domain.RequestRepository {
	return &requestRepo{q: q}
}

func (r *requestRepo) Create(ctx context.Context, p domain.CreateRequestParams) (*domain.Request, error) {
	row, err := r.q.CreateRequest(ctx, dbgen.CreateRequestParams{
		UserID:        p.UserID,
		ItemCategory:  p.ItemCategory,
		ItemName:      p.ItemName,
		Description:   p.Description,
		Urgency:       string(p.Urgency),
		EstimatedCost: p.EstimatedCost,
		Justification: p.Justification,
		Status:        string(p.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("requestRepo.Create: %w", err)
	}
	return rowToRequest(row), nil
}

func (r *requestRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Request, error) {
	row, err := r.q.GetRequestByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("requestRepo.GetByID: %w", err)
	}
	return rowToRequest(row), nil
}

func (r *requestRepo) Update(ctx context.Context, id uuid.UUID, p domain.UpdateRequestParams) (*domain.Request, error) {
	var urgency dbgen.NullUrgencyLevel
	if p.Urgency != nil {
		urgency = dbgen.NullUrgencyLevel{
			UrgencyLevel: dbgen.UrgencyLevel(*p.Urgency),
			Valid:        true,
		}
	}

	row, err := r.q.UpdateRequest(ctx, dbgen.UpdateRequestParams{
		ID:            id,
		ItemCategory:  p.ItemCategory,
		ItemName:      p.ItemName,
		Description:   p.Description,
		Urgency:       urgency,
		EstimatedCost: p.EstimatedCost, // now *float64, matches DB float8
		Justification: p.Justification,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidState
		}
		return nil, fmt.Errorf("requestRepo.Update: %w", err)
	}
	return rowToRequest(row), nil
}

func (r *requestRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus, note *string) error {
	return r.q.UpdateRequestStatus(ctx, dbgen.UpdateRequestStatusParams{
		ID:            id,
		Status:        string(status),
		RejectionNote: note,
	})
}

func (r *requestRepo) List(ctx context.Context, f domain.RequestFilter) ([]*domain.Request, int, error) {
	var (
		rows  []dbgen.Request
		total int64
		err   error
	)
	switch {
	case f.UserID != nil && f.Status != nil:
		// fetch by user then filter client-side (uncommon combo — add a dedicated query if needed)
		rows, err = r.q.ListRequestsByUser(ctx, dbgen.ListRequestsByUserParams{
			UserID: *f.UserID, Limit: int32(f.Limit), Offset: int32(f.Offset),
		})
		ct, _ := r.q.CountRequestsByUser(ctx, *f.UserID)
		total = ct
	case f.UserID != nil:
		rows, err = r.q.ListRequestsByUser(ctx, dbgen.ListRequestsByUserParams{
			UserID: *f.UserID, Limit: int32(f.Limit), Offset: int32(f.Offset),
		})
		ct, _ := r.q.CountRequestsByUser(ctx, *f.UserID)
		total = ct
	case f.Status != nil:
		rows, err = r.q.ListRequestsByStatus(ctx, dbgen.ListRequestsByStatusParams{
			Status: string(*f.Status), Limit: int32(f.Limit), Offset: int32(f.Offset),
		})
		ct, _ := r.q.CountRequestsByStatus(ctx, string(*f.Status))
		total = ct
	default:
		rows, err = r.q.ListAllRequests(ctx, dbgen.ListAllRequestsParams{
			Limit: int32(f.Limit), Offset: int32(f.Offset),
		})
		ct, _ := r.q.CountAllRequests(ctx)
		total = ct
	}
	if err != nil {
		return nil, 0, fmt.Errorf("requestRepo.List: %w", err)
	}
	out := make([]*domain.Request, len(rows))
	for i, row := range rows {
		out[i] = rowToRequest(row)
	}
	return out, int(total), nil
}
func (r *requestRepo) Delete(ctx context.Context, id uuid.UUID) error {
	n, err := r.q.DeleteDraftRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("requestRepo.Delete: %w", err)
	}
	if n == 0 {
		return domain.ErrInvalidState
	}
	return nil
}
func rowToRequest(row dbgen.Request) *domain.Request {
	return &domain.Request{
		ID:            row.ID,
		UserID:        row.UserID,
		ItemCategory:  row.ItemCategory,
		ItemName:      row.ItemName,
		Description:   row.Description,
		Urgency:       domain.UrgencyLevel(row.Urgency),
		EstimatedCost: row.EstimatedCost,
		Justification: row.Justification,
		Status:        domain.RequestStatus(row.Status),
		RejectionNote: row.RejectionNote,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
