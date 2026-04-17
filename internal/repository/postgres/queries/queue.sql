-- name: MaxQueuePosition :one
SELECT COALESCE(MAX(position), 0)::int4 AS max_position FROM queue_entries;

-- name: EnqueueRequest :one
INSERT INTO queue_entries (request_id, priority_score, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetQueueEntryByRequestID :one
SELECT * FROM queue_entries WHERE request_id = $1 LIMIT 1;

-- name: GetQueuePosition :one
SELECT position FROM queue_entries WHERE request_id = $1 LIMIT 1;

-- name: UpdateQueueScore :exec
UPDATE queue_entries
SET priority_score = $2, position = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpdateQueueFunding :exec
UPDATE queue_entries
SET funding_progress = $2, updated_at = NOW()
WHERE id = $1;

-- name: DequeueRequest :exec
DELETE FROM queue_entries WHERE request_id = $1;

-- name: ListQueueEntries :many
SELECT * FROM queue_entries
ORDER BY priority_score DESC, entered_at ASC
LIMIT $1 OFFSET $2;

-- name: CountQueueEntries :one
SELECT COUNT(*) FROM queue_entries;
