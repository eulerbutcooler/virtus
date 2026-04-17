-- name: CreateInstitution :one
INSERT INTO institutions (user_id, name, type, contact_email, website, esg_goals)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInstitutionByID :one
SELECT * FROM institutions WHERE id = $1 LIMIT 1;

-- name: GetInstitutionByUserID :one
SELECT * FROM institutions WHERE user_id = $1 LIMIT 1;

-- name: CreateInstitutionalContribution :one
INSERT INTO institutional_contributions
    (institution_id, pool_id, amount, currency, category_tag, region_tag)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListInstitutionalContributions :many
SELECT * FROM institutional_contributions
WHERE institution_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountInstitutionalContributions :one
SELECT COUNT(*) FROM institutional_contributions WHERE institution_id = $1;
