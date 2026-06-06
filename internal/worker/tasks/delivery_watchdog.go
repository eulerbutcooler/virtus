package tasks

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DeliveryWatchdogTask finds deliveries that have been stuck in 'in_transit'
// for longer than the stale threshold and logs them. An admin can then
// investigate and either mark them failed or follow up with the carrier.
//
// This task does NOT automatically fail deliveries — that requires human judgement.
type DeliveryWatchdogTask struct {
	db             *pgxpool.Pool
	staleThreshold time.Duration
}

func NewDeliveryWatchdogTask(db *pgxpool.Pool) *DeliveryWatchdogTask {
	return &DeliveryWatchdogTask{
		db:             db,
		staleThreshold: 14 * 24 * time.Hour, // 14 days
	}
}

type staleDelivery struct {
	ID             uuid.UUID
	FulfillmentID  uuid.UUID
	TrackingNumber *string
	Carrier        *string
	DaysSince      int
}

func (t *DeliveryWatchdogTask) Run(ctx context.Context) error {
	stale, err := t.findStale(ctx)
	if err != nil {
		return fmt.Errorf("deliveryWatchdog: %w", err)
	}
	if len(stale) == 0 {
		return nil
	}
	for _, d := range stale {
		tracking := "(no tracking number)"
		if d.TrackingNumber != nil {
			tracking = *d.TrackingNumber
		}
		carrier := "(unknown carrier)"
		if d.Carrier != nil {
			carrier = *d.Carrier
		}
		slog.Warn("deliveryWatchdog: stale in-transit delivery",
			"delivery_id", d.ID,
			"fulfillment_id", d.FulfillmentID,
			"tracking", tracking,
			"carrier", carrier,
			"days_in_transit", d.DaysSince,
		)
	}
	slog.Warn("deliveryWatchdog: cycle complete", "stale_count", len(stale))
	return nil
}

func (t *DeliveryWatchdogTask) findStale(ctx context.Context) ([]staleDelivery, error) {
	cutoff := time.Now().Add(-t.staleThreshold)

	const query = `
		SELECT id, fulfillment_id, tracking_number, carrier, updated_at
		FROM deliveries
		WHERE status     = 'in_transit'
		  AND updated_at < $1
		ORDER BY updated_at ASC
	`

	rows, err := t.db.Query(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("query stale deliveries: %w", err)
	}
	defer rows.Close()

	var results []staleDelivery
	now := time.Now()

	for rows.Next() {
		var (
			id             uuid.UUID
			fulfillmentID  uuid.UUID
			trackingNumber *string
			carrier        *string
			updatedAt      time.Time
		)
		if err := rows.Scan(&id, &fulfillmentID, &trackingNumber, &carrier, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan stale delivery: %w", err)
		}
		results = append(results, staleDelivery{
			ID:             id,
			FulfillmentID:  fulfillmentID,
			TrackingNumber: trackingNumber,
			Carrier:        carrier,
			DaysSince:      int(now.Sub(updatedAt).Hours() / 24),
		})
	}

	return results, rows.Err()
}
