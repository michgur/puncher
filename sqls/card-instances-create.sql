-- Loyalty Card Instances Table
CREATE TABLE IF NOT EXISTS CardInstances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    card_id VARCHAR NOT NULL,
    slots INTEGER NOT NULL
);