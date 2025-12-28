package models

import "time"

type RegistrationCode struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (r RegistrationCode) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}
