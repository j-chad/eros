package models

import "time"

type PushSubscription struct {
	DeviceID       string     `json:"deviceID"`
	Endpoint       string     `json:"endpoint"`
	ExpirationTime *time.Time `json:"expirationTime"`
	Keys           PushKeys   `json:"keys"`
}

type PushKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type PushUrgency string

const (
	// PushUrgencyVeryLow requires device state: on power and Wi-Fi
	PushUrgencyVeryLow PushUrgency = "very-low"
	// PushUrgencyLow requires device state: on either power or Wi-Fi
	PushUrgencyLow PushUrgency = "low"
	// PushUrgencyNormal excludes device state: low battery
	PushUrgencyNormal PushUrgency = "normal"
	// PushUrgencyHigh admits device state: low battery
	PushUrgencyHigh PushUrgency = "high"
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
	Tag   string   `json:"tag,omitempty"`
	Data  PushData `json:"data,omitempty"`
}

type PushData struct {
	URL string `json:"url,omitempty"`
}

type PushSendResult struct {
	Sent    int `json:"sent"`
	Failed  int `json:"failed"`
	Cleaned int `json:"cleaned"`
}
