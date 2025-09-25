PRAGMA journal_mode=WAL;

BEGIN;

ALTER TABLE users ADD COLUMN last_alum_check_timestamp INTEGER;

COMMIT;