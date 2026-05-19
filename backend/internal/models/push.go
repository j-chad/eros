package models

import "time"

type PushSubscription struct {
	Endpoint string   `json:"endpoint"`
	Keys     PushKeys `json:"keys"`
}

type PushKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type PushUrgency string

const (
	// UrgencyVeryLow requires device state: on power and Wi-Fi
	UrgencyVeryLow PushUrgency = "very-low"
	// UrgencyLow requires device state: on either power or Wi-Fi
	UrgencyLow PushUrgency = "low"
	// UrgencyNormal excludes device state: low battery
	UrgencyNormal PushUrgency = "normal"
	// UrgencyHigh admits device state: low battery
	UrgencyHigh PushUrgency = "high"
)

type PushRequest struct {
	Message PushMessage   `json:"message"`
	Topic   string        `json:"topic"`
	TTL     time.Duration `json:"ttl"`
	Urgency PushUrgency   `json:"urgency"`
}

type PushMessage struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tag   string   `json:"tag"`
	Data  PushData `json:"data"`
}

type PushData struct {
	URL string `json:"url"`
}
