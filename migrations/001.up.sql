PRAGMA journal_mode=WAL;

BEGIN;

CREATE TABLE IF NOT EXISTS users (
    buck_id TEXT PRIMARY KEY,
    discord_id INTEGER,
    name_num TEXT NOT NULL,
    display_name TEXT NOT NULL,
    last_seen_timestamp INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    last_attended_timestamp INTEGER,
    
    added_to_mailinglist INTEGER NOT NULL DEFAULT (FALSE),

    -- 0 or 1 depending on if the user has the affiliation
    student INTEGER NOT NULL,
    alum INTEGER NOT NULL,
    employee INTEGER NOT NULL,
    faculty INTEGER NOT NULL
) STRICT, WITHOUT ROWID;

CREATE INDEX IF NOT EXISTS discord_id ON users (discord_id);
CREATE INDEX IF NOT EXISTS buck_id ON users (buck_id);
CREATE INDEX IF NOT EXISTS student ON users (student);

CREATE TABLE IF NOT EXISTS attendance_records (
    user_id INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    -- 0 for in person, 1 for online
    kind INTEGER NOT NULL,
    PRIMARY KEY (user_id, timestamp),
    FOREIGN KEY (user_id) REFERENCES users(buck_id)
) STRICT, WITHOUT ROWID;

CREATE TABLE IF NOT EXISTS elections (
    election_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,

    -- 0 means not yet started
    -- 1 means started and in progress (should only be exactly one election at a time)
    -- >1 means ended at that unix timestamp
    timestamp INTEGER NOT NULL DEFAULT (0)
) STRICT;

CREATE INDEX IF NOT EXISTS timestamp ON elections (timestamp);

CREATE TABLE IF NOT EXISTS candidates (
    candidate_id INTEGER PRIMARY KEY,
    election_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    votes INTEGER NOT NULL DEFAULT 0
) STRICT;

CREATE INDEX IF NOT EXISTS election_id ON candidates (election_id);

CREATE TABLE IF NOT EXISTS votes (
    election_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (election_id, user_id)
) STRICT, WITHOUT ROWID;

COMMIT;
