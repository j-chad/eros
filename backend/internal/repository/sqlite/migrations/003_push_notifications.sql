CREATE TABLE push_subscription (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    expires_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (device_id) REFERENCES device(id) ON DELETE CASCADE
);

ALTER TABLE graph ADD COLUMN notified_at DATETIME; -- when the user was last notified about this graph, used to avoid sending multiple notifications for the same graph

-- +down
DROP TABLE push_subscription;
ALTER TABLE graph DROP COLUMN notified_at;