package models

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusM   int     `json:"radius_m"`
}

type CodeData struct {
	Code string `json:"code"`
}

type RewardData struct {
	RewardType  string `json:"reward_type"`
	Content     string `json:"content"`
	MediaType   string `json:"media_type"`
	GiveFavours int    `json:"give_favours"`
}

type StartData struct {
	StartingAt   string   `json:"starting_at"`
	ViewportX    *float64 `json:"viewport_x,omitempty"`
	ViewportY    *float64 `json:"viewport_y,omitempty"`
	ViewportZoom *float64 `json:"viewport_zoom,omitempty"`
}
