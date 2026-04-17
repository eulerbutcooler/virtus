-- name: GetPool :one
SELECT * FROM pools WHERE id = sqlc.arg(id) LIMIT 1;

-- name: CreditPool :exec
UPDATE pools
SET
    balance    = balance + sqlc.arg(amount),
    total_in   = total_in + sqlc.arg(amount),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: DebitPool :execrows
UPDATE pools
SET
    balance    = balance - sqlc.arg(amount),
    total_out  = total_out + sqlc.arg(amount),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND balance >= sqlc.arg(amount);
