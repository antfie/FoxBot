CREATE TABLE bayes_model (
    feed_group TEXT NOT NULL,
    word       TEXT NOT NULL,
    relevant   INTEGER DEFAULT 0,
    irrelevant INTEGER DEFAULT 0,
    PRIMARY KEY (feed_group, word)
);

CREATE TABLE bayes_article (
    hash       TEXT PRIMARY KEY,
    feed_group TEXT NOT NULL,
    title      TEXT NOT NULL,
    created    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bayes_stats (
    feed_group TEXT PRIMARY KEY,
    relevant   INTEGER DEFAULT 0,
    irrelevant INTEGER DEFAULT 0
);

CREATE TABLE telegram_state (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
