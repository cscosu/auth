PRAGMA journal_mode=WAL;
PRAGMA user_version=2;
BEGIN;
ALTER TABLE users ADD COLUMN last_alum_check_timestamp INTEGER;
COMMIT;