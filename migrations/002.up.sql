BEGIN;
PRAGMA user_version=2;
ALTER TABLE users ADD COLUMN last_alum_check_timestamp INTEGER;
COMMIT;