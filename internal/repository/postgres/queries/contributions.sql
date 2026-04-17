-- name: CreateContribution :one
INSERT INTO contributions (user_id, pool_id, amount, currency, stripe_pi_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetContributionByID :one
SELECT * FROM contributions WHERE id = $1 LIMIT 1;

-- name: GetContributionByStripePI :one
SELECT * FROM contributions WHERE stripe_pi_id = $1 LIMIT 1;

-- name: UpdateContributionStatus :exec
UPDATE contributions
SET
    status      = $2,
    payment_ref = $3,
    updated_at  = NOW()
WHERE id = $1;

-- name: ListContributionsByUser :many
SELECT * FROM contributions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountContributionsByUser :one
SELECT COUNT(*) FROM contributions WHERE user_id = $1;

-- name: SumCompletedContributionsByUser :one
SELECT COALESCE(SUM(amount), 0)::float8
FROM contributions
WHERE user_id = $1 AND status = 'completed';
