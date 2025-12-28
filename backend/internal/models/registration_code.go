package models

import "time"

type RegistrationCode struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"time"`
}
