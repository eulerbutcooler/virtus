ALTER TABLE queue_entries DROP COLUMN priority_score;
DROP INDEX IF EXISTS idx_queue_entries_priority;
