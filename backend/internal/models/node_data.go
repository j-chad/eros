package models

import "time"

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusM   int     `json:"radius_m"`
}

type CodeData struct {
	Code string `json:"code"`
}

type ManualData struct {
	Instructions string     `json:"instructions"`
	UnlockedAt   *time.Time `json:"unlocked_at"`
}

type RewardData struct {
	RewardType  string `json:"reward_type"`
	Content     string `json:"content"`
	MediaType   string `json:"media_type"`
	GiveFavours int    `json:"give_favours"`
}
