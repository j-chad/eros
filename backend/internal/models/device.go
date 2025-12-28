package models

import "time"

type Device struct {
	ID           string    `json:"id"`
	DeviceInfo   string    `json:"device_info"`
	RegisteredAt time.Time `json:"registered_at"`
	LastSeenAt   time.Time `json:"last_seen_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (d Device) IsExpired() bool {
	return d.ExpiresAt.Before(time.Now())
}
