package tasks

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ImpactReminderTask finds deliveries that were verified at least 30 days ago
// and have no impact record on file. It logs each one so an admin or a future
// notification system can prompt the member for a follow-up survey.
//
// Intervals checked: 30 days, 90 days, 180 days.
// A separate impact record per interval is expected; this task fires whenever
// ANY interval is overdue and no record for that interval exists yet.
type ImpactReminderTask struct {
	db *pgxpool.Pool
}

func NewImpactReminderTask(db *pgxpool.Pool) *ImpactReminderTask {
	return &ImpactReminderTask{db: db}
}

// intervalDays lists the follow-up checkpoints, in days since delivery.
var intervalDays = []struct {
	label string
	days  int
}{
	{"30_day", 30},
	{"90_day", 90},
	{"180_day", 180},
}

type pendingReminder struct {
	DeliveryID uuid.UUID
	UserID     uuid.UUID
	VerifiedAt time.Time
	Interval   string
}

func (t *ImpactReminderTask) Run(ctx context.Context) error {
	reminders, err := t.findPending(ctx)
	if err != nil {
		return fmt.Errorf("impactReminder: %w", err)
	}
	if len(reminders) == 0 {
		return nil
	}
	for _, r := range reminders {
		slog.Info("impactReminder: follow-up needed",
			"delivery_id", r.DeliveryID,
			"user_id", r.UserID,
			"interval", r.Interval,
			"verified_at", r.VerifiedAt.Format(time.DateOnly),
		)
	}
	slog.Info("impactReminder: cycle complete", "pending_count", len(reminders))
	return nil
}

// findPending returns deliveries that are overdue for a specific interval survey
// and have no existing impact record for that interval label.
func (t *ImpactReminderTask) findPending(ctx context.Context) ([]pendingReminder, error) {
	// Pull all verified deliveries joined with their fulfillment → request → user.
	// We then check in Go which intervals are missing rather than doing a complex
	// NOT EXISTS per-interval in SQL.
	const deliveriesQuery = `
		SELECT
			d.id          AS delivery_id,
			r.user_id,
			d.verified_at
		FROM deliveries d
		JOIN fulfillments f ON f.id = d.fulfillment_id
		JOIN requests     r ON r.id = f.request_id
		WHERE d.status      = 'delivered'
		  AND d.verified_at IS NOT NULL
		  AND d.verified_at < NOW() - INTERVAL '30 days'
		ORDER BY d.verified_at ASC
	`

	rows, err := t.db.Query(ctx, deliveriesQuery)
	if err != nil {
		return nil, fmt.Errorf("query deliveries: %w", err)
	}
	defer rows.Close()

	type deliveryRow struct {
		deliveryID uuid.UUID
		userID     uuid.UUID
		verifiedAt time.Time
	}
	var deliveries []deliveryRow

	for rows.Next() {
		var dr deliveryRow
		if err := rows.Scan(&dr.deliveryID, &dr.userID, &dr.verifiedAt); err != nil {
			return nil, fmt.Errorf("scan delivery row: %w", err)
		}
		deliveries = append(deliveries, dr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate deliveries: %w", err)
	}

	if len(deliveries) == 0 {
		return nil, nil
	}

	// For each delivery, fetch existing impact interval labels.
	const labelsQuery = `
		SELECT interval_label FROM impact_records WHERE delivery_id = $1
	`

	var reminders []pendingReminder
	now := time.Now()

	for _, d := range deliveries {
		labelRows, err := t.db.Query(ctx, labelsQuery, d.deliveryID)
		if err != nil {
			slog.Warn("impactReminder: fetch labels", "delivery_id", d.deliveryID, "err", err)
			continue
		}

		existing := make(map[string]bool)
		for labelRows.Next() {
			var label string
			if err := labelRows.Scan(&label); err != nil {
				labelRows.Close()
				break
			}
			existing[label] = true
		}
		labelRows.Close()

		daysSince := int(now.Sub(d.verifiedAt).Hours() / 24)

		for _, iv := range intervalDays {
			if daysSince >= iv.days && !existing[iv.label] {
				reminders = append(reminders, pendingReminder{
					DeliveryID: d.deliveryID,
					UserID:     d.userID,
					VerifiedAt: d.verifiedAt,
					Interval:   iv.label,
				})
			}
		}
	}

	return reminders, nil
}
