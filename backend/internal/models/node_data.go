package models

import "time"

type LocationData struct {
	LocationArea
	ShowHint bool          `json:"show_hint"`
	Hint     *LocationArea `json:"hint,omitempty"`
}

// ClientLocationData is the sanitized version sent to clients.
// It omits the exact coordinates to prevent cheating.
type ClientLocationData struct {
	RadiusM int           `json:"radius_m"`
	Hint    *LocationArea `json:"hint,omitempty"`
}

type LocationArea struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusM   int     `json:"radius_m"`
}

type CodeData struct {
	Codes []string `json:"codes"`
}

type ManualData struct {
	Instructions string     `json:"instructions"`
	UnlockedAt   *time.Time `json:"unlocked_at"`
}

type TimeData struct {
	UnlockAt time.Time `json:"unlock_at"`
}

type RewardData struct {
	RewardType  string `json:"reward_type"`
	Payload     string `json:"payload"`
	GiveFavours int    `json:"give_favours"`
}
