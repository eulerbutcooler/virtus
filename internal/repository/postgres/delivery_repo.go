package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eulerbutcooler/virtus/internal/domain"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type deliveryRepo struct {
	q *dbgen.Queries
}

func NewDeliveryRepository(q *dbgen.Queries) domain.DeliveryRepository {
	return &deliveryRepo{q: q}
}

func (r *deliveryRepo) Create(ctx context.Context, p domain.CreateDeliveryParams) (*domain.Delivery, error) {
	row, err := r.q.CreateDelivery(ctx, dbgen.CreateDeliveryParams{
		FulfillmentID:  p.FulfillmentID,
		TrackingNumber: p.TrackingNumber,
		Carrier:        p.Carrier,
	})
	if err != nil {
		return nil, fmt.Errorf("deliveryRepo.Create: %w", err)
	}
	return rowToDelivery(row), nil
}

func (r *deliveryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error) {
	row, err := r.q.GetDeliveryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("deliveryRepo.GetByID: %w", err)
	}
	return rowToDelivery(row), nil
}

func (r *deliveryRepo) GetByFulfillmentID(ctx context.Context, fulfillmentID uuid.UUID) (*domain.Delivery, error) {
	row, err := r.q.GetDeliveryByFulfillmentID(ctx, fulfillmentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("deliveryRepo.GetByFulfillmentID: %w", err)
	}
	return rowToDelivery(row), nil
}

func (r *deliveryRepo) Verify(ctx context.Context, id uuid.UUID, p domain.VerifyDeliveryParams) (*domain.Delivery, error) {
	row, err := r.q.VerifyDelivery(ctx, dbgen.VerifyDeliveryParams{
		ID:            id,
		ProofPhotoUrl: &p.ProofPhotoURL,
		DeliveredAt:   pgtype.Timestamptz{Time: p.DeliveredAt, Valid: true},
		VerifiedBy:    pgtype.UUID{Bytes: p.VerifiedBy, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("deliveryRepo.Verify: %w", err)
	}
	return rowToDelivery(row), nil
}

func (r *deliveryRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus) error {
	if err := r.q.UpdateDeliveryStatus(ctx, dbgen.UpdateDeliveryStatusParams{
		ID:     id,
		Status: string(status),
	}); err != nil {
		return fmt.Errorf("deliveryRepo.UpdateStatus: %w", err)
	}
	return nil
}

func rowToDelivery(row dbgen.Delivery) *domain.Delivery {
	d := &domain.Delivery{
		ID:             row.ID,
		FulfillmentID:  row.FulfillmentID,
		TrackingNumber: row.TrackingNumber,
		Carrier:        row.Carrier,
		ProofPhotoURL:  row.ProofPhotoUrl,
		Status:         domain.DeliveryStatus(row.Status),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
	if row.DeliveredAt.Valid {
		t := row.DeliveredAt.Time
		d.DeliveredAt = &t
	}
	if row.VerifiedAt.Valid {
		t := row.VerifiedAt.Time
		d.VerifiedAt = &t
	}
	if row.VerifiedBy.Valid {
		id := uuid.UUID(row.VerifiedBy.Bytes)
		d.VerifiedBy = &id
	}
	return d
}

// For converting time.Time → pgtype.Timestamptz.
func timePtr(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
