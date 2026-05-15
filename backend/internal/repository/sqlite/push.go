package sqlite

import (
	"backend/internal/models"
	"context"
)

func (s *sqliteDB) CreatePushSubscription(ctx context.Context, deviceID string, sub models.PushSubscription) error {
	_, err := s.executor().ExecContext(ctx, `
		INSERT INTO push_subscription (device_id, endpoint, p256dh, auth)
		VALUES (?, ?, ?, ?)
	`, deviceID, sub.Endpoint, sub.P256DH, sub.Auth)
	return err
}

func (s *sqliteDB) DeletePushSubscriptions(ctx context.Context, deviceID string) error {
	_, err := s.executor().ExecContext(ctx, `DELETE FROM push_subscription WHERE device_id = ?`, deviceID)
	return err
}

func (s *sqliteDB) GetPushSubscriptions(ctx context.Context) ([]models.PushSubscription, error) {
	rows, err := s.executor().QueryContext(ctx, `SELECT endpoint, p256dh, auth FROM push_subscription`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.PushSubscription
	for rows.Next() {
		var sub models.PushSubscription
		if err := rows.Scan(&sub.Endpoint, &sub.P256DH, &sub.Auth); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, rows.Err()
}
