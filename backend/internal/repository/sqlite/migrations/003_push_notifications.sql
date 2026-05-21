CREATE TABLE push_subscription (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL UNIQUE,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    expires_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (device_id) REFERENCES device(id) ON DELETE CASCADE
);

ALTER TABLE graph ADD COLUMN notified_at DATETIME; -- when the user was last notified about this graph, used to avoid sending multiple notifications for the same graph


CREATE TRIGGER update_push_subscription_updated_at
AFTER UPDATE ON push_subscription
FOR EACH ROW
BEGIN
    UPDATE push_subscription SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- +down
DROP TABLE push_subscription;
ALTER TABLE graph DROP COLUMN notified_at;
DROP TRIGGER update_push_subscription_updated_at;