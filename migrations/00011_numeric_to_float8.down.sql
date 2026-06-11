ALTER TABLE contributions
    ALTER COLUMN amount TYPE NUMERIC(15,2) USING amount::numeric;

ALTER TABLE pools
    ALTER COLUMN balance   TYPE NUMERIC(15,2) USING balance::numeric,
    ALTER COLUMN total_in  TYPE NUMERIC(15,2) USING total_in::numeric,
    ALTER COLUMN total_out TYPE NUMERIC(15,2) USING total_out::numeric;

ALTER TABLE queue_entries
    ALTER COLUMN funding_progress TYPE NUMERIC(15,2) USING funding_progress::numeric;

ALTER TABLE requests
    ALTER COLUMN estimated_cost TYPE NUMERIC(15,2) USING estimated_cost::numeric;
