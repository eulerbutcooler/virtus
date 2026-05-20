ALTER TABLE queue_entries AND COLUMN priority_score NUMERIC(10,4) NOT NULL DEFAULT 0.0;
CREATE INDEX idx_queue_entries_priority ON queue_entries(priority_score DESC);
