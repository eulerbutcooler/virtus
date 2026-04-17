-- name: CreateDelivery :one
INSERT INTO deliveries (fulfillment_id, tracking_number, carrier)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDeliveryByID :one
SELECT * FROM deliveries WHERE id = $1 LIMIT 1;

-- name: GetDeliveryByFulfillmentID :one
SELECT * FROM deliveries WHERE fulfillment_id = $1 LIMIT 1;

-- name: VerifyDelivery :one
UPDATE deliveries
SET
    proof_photo_url = $2,
    status          = 'delivered',
    delivered_at    = $3,
    verified_at     = NOW(),
    verified_by     = $4,
    updated_at      = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateDeliveryStatus :exec
UPDATE deliveries
SET status = $2, updated_at = NOW()
WHERE id = $1;
