-- name: CreateFulfillment :one
INSERT INTO fulfillments (request_id)
VALUES ($1)
RETURNING *;

-- name: GetFulfillmentByID :one
SELECT * FROM fulfillments WHERE id = $1 LIMIT 1;

-- name: GetFulfillmentByRequestID :one
SELECT * FROM fulfillments WHERE request_id = $1 LIMIT 1;

-- name: UpdateFulfillment :one
UPDATE fulfillments
SET
    vendor_name        = COALESCE(sqlc.narg('vendor_name'), vendor_name),
    vendor_ref         = COALESCE(sqlc.narg('vendor_ref'), vendor_ref),
    actual_cost        = COALESCE(sqlc.narg('actual_cost'), actual_cost),
    procurement_status = COALESCE(sqlc.narg('procurement_status'), procurement_status),
    notes              = COALESCE(sqlc.narg('notes'), notes),
    procured_at        = COALESCE(sqlc.narg('procured_at'), procured_at),
    updated_at         = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListFulfillments :many
SELECT * FROM fulfillments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountFulfillments :one
SELECT COUNT(*) FROM fulfillments;
