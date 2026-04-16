CREATE TYPE institution_type AS ENUM ('corporation', 'ngo', 'government', 'foundation', 'other');

CREATE TABLE institutions (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT,
    name           VARCHAR(255) NOT NULL,
    type           institution_type NOT NULL DEFAULT 'other',
    contact_email  VARCHAR(255) NOT NULL,
    website        VARCHAR(500),
    esg_goals      JSONB NOT NULL DEFAULT '{}',
    verified       BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE institutional_contributions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id  UUID NOT NULL REFERENCES institutions(id) ON DELETE RESTRICT,
    pool_id         UUID NOT NULL REFERENCES pools(id) ON DELETE RESTRICT,
    amount          NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL DEFAULT 'USD',
    status          contribution_status NOT NULL DEFAULT 'pending',
    payment_ref     VARCHAR(255),
    category_tag    VARCHAR(100),   -- 'mobility_aids', 'education'
    region_tag      VARCHAR(100),   -- 'UK', 'South Asia'
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_institutions_user_id ON institutions(user_id);
CREATE INDEX idx_inst_contributions_institution_id ON institutional_contributions(institution_id);
