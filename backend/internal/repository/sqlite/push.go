package sqlite

import (
	"backend/internal/models"
	"context"
)

func (s *sqliteDB) CreatePushSubscription(ctx context.Context, deviceID string, sub models.PushSubscription) error {
	_, err := s.executor().ExecContext(ctx, `
		INSERT INTO push_subscription (device_id, endpoint, p256dh, auth, expires_at) VALUES (?, ?, ?, ?, ?) 
		ON CONFLICT(device_id) 
		    DO UPDATE SET 
		    	endpoint = excluded.endpoint, 
		    	p256dh = excluded.p256dh, 
		    	auth = excluded.auth, 
		    	expires_at = excluded.expires_at
	`, deviceID, sub.Endpoint, sub.Keys.P256dh, sub.Keys.Auth, sub.ExpirationTime)
	return err
}

func (s *sqliteDB) DeletePushSubscriptions(ctx context.Context, deviceID string) error {
	_, err := s.executor().ExecContext(ctx, `DELETE FROM push_subscription WHERE device_id = ?`, deviceID)
	return err
}

func (s *sqliteDB) GetPushSubscriptions(ctx context.Context) ([]models.PushSubscription, error) {
	//language=sqlite
	rows, err := s.executor().QueryContext(ctx, `
		SELECT ps.device_id, d.device_info, ps.endpoint, ps.p256dh, ps.auth, ps.created_at, ps.updated_at
		FROM push_subscription ps
		JOIN device d ON d.id = ps.device_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.PushSubscription
	for rows.Next() {
		var sub models.PushSubscription
		if err := rows.Scan(&sub.DeviceID, &sub.DeviceName, &sub.Endpoint, &sub.Keys.P256dh, &sub.Keys.Auth, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, rows.Err()
}
