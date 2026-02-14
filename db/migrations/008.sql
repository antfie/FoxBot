CREATE TABLE IF NOT EXISTS weather_notification (
    location TEXT PRIMARY KEY,
    last_notified DATE NOT NULL
);
