package models

import (
	"backend/pkg/apierror"
	"time"
)

type RegisterDeviceRequest struct {
	DeviceInfo       string `json:"device_info"`
	RegistrationCode string `json:"registration_code"`
}

func (r RegisterDeviceRequest) Validate() *apierror.APIError {
	if r.RegistrationCode == "" {
		return apierror.BadRequest("registration_code is required")
	}

	return nil
}

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
