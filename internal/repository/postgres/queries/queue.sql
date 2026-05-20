-- name: MaxQueuePosition :one
SELECT COALESCE(MAX(position), 0)::int4 AS max_position FROM queue_entries;

-- name: EnqueueRequest :one
INSERT INTO queue_entries (request_id, position)
VALUES ($1, $2)
RETURNING *;

-- name: GetQueueEntryByRequestID :one
SELECT * FROM queue_entries WHERE request_id = $1 LIMIT 1;

-- name: GetQueuePosition :one
SELECT position FROM queue_entries WHERE request_id = $1 LIMIT 1;

-- name: UpdateQueuePosition :exec
UPDATE queue_entries
SET position = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateQueueFunding :exec
UPDATE queue_entries
SET funding_progress = $2, updated_at = NOW()
WHERE id = $1;

-- name: DequeueRequest :exec
DELETE FROM queue_entries WHERE request_id = $1;

-- name: ListQueueEntries :many
SELECT * FROM queue_entries
ORDER BY entered_at ASC
LIMIT $1 OFFSET $2;

-- name: CountQueueEntries :one
SELECT COUNT(*) FROM queue_entries;

-- name: ListAllQueueEntries :many
SELECT * FROM queue_entries
ORDER BY entered_at ASC;
