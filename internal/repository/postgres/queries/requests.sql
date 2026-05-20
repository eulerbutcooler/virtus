-- name: CreateRequest :one
INSERT INTO requests (user_id, item_category, item_name, description, urgency, estimated_cost, justification, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetRequestByID :one
SELECT * FROM requests WHERE id = $1 LIMIT 1;

-- name: UpdateRequest :one
UPDATE requests
SET
    item_category  = COALESCE(sqlc.narg('item_category'), item_category),
    item_name      = COALESCE(sqlc.narg('item_name'), item_name),
    description    = COALESCE(sqlc.narg('description'), description),
    urgency        = COALESCE(sqlc.narg('urgency'), urgency),
    estimated_cost = COALESCE(sqlc.narg('estimated_cost'), estimated_cost),
    justification  = COALESCE(sqlc.narg('justification'), justification),
    updated_at     = NOW()
WHERE id = sqlc.arg('id') AND status IN ('draft', 'submitted')
RETURNING *;

-- name: UpdateRequestStatus :exec
UPDATE requests
SET status = $2, rejection_note = $3, updated_at = NOW()
WHERE id = $1;

-- name: ListRequestsByUser :many
SELECT * FROM requests
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountRequestsByUser :one
SELECT COUNT(*) FROM requests WHERE user_id = $1;

-- name: ListRequestsByStatus :many
SELECT * FROM requests
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountRequestsByStatus :one
SELECT COUNT(*) FROM requests WHERE status = $1;

-- name: ListAllRequests :many
SELECT * FROM requests
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllRequests :one
SELECT COUNT(*) FROM requests;

-- name: DeleteDraftRequest :execrows
DELETE FROM requests WHERE id = $1 AND status = 'draft';
