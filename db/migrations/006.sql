CREATE TABLE IF NOT EXISTS http_cache (
    url            TEXT PRIMARY KEY,
    etag           TEXT NOT NULL DEFAULT '',
    last_modified  TEXT NOT NULL DEFAULT '',
    fail_count     INTEGER NOT NULL DEFAULT 0
);
