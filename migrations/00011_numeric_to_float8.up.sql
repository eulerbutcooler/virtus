-- Change NUMERIC columns to DOUBLE PRECISION (float8) so pgx/sqlc
-- can scan them natively as float64 without pgtype.Numeric wrappers.
-- NUMERIC(15,2) → DOUBLE PRECISION is acceptable for this MVP;
-- if sub-cent precision ever matters, use a BIGINT (cents) approach instead.

ALTER TABLE requests
    ALTER COLUMN estimated_cost TYPE DOUBLE PRECISION USING estimated_cost::float8;

ALTER TABLE queue_entries
    ALTER COLUMN funding_progress TYPE DOUBLE PRECISION USING funding_progress::float8;

ALTER TABLE pools
    ALTER COLUMN balance   TYPE DOUBLE PRECISION USING balance::float8,
    ALTER COLUMN total_in  TYPE DOUBLE PRECISION USING total_in::float8,
    ALTER COLUMN total_out TYPE DOUBLE PRECISION USING total_out::float8;

ALTER TABLE contributions
    ALTER COLUMN amount TYPE DOUBLE PRECISION USING amount::float8;
