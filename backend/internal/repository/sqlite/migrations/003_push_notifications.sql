CREATE TABLE push_subscription (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL FOREIGN KEY REFERENCES device(id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE graph ADD COLUMN notified_at DATETIME; -- when the user was last notified about this graph, used to avoid sending multiple notifications for the same graph