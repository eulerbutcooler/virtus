-- name: CreateImpactRecord :one
INSERT INTO impact_records
    (delivery_id, user_id, interval_label, outcome_description, satisfaction_score, metrics)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetImpactRecordByID :one
SELECT * FROM impact_records WHERE id = $1 LIMIT 1;

-- name: ListImpactByDelivery :many
SELECT * FROM impact_records
WHERE delivery_id = $1
ORDER BY recorded_at ASC;

-- name: ListImpactByUser :many
SELECT * FROM impact_records
WHERE user_id = $1
ORDER BY recorded_at DESC
LIMIT $2 OFFSET $3;

-- name: CountImpactByUser :one
SELECT COUNT(*) FROM impact_records WHERE user_id = $1;
