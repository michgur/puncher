-- Loyalty Card Instances Table
DROP TABLE IF EXISTS CardInstances;

CREATE TABLE CardInstances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    card_id INTEGER NOT NULL,
    slots INTEGER NOT NULL
);