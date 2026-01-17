package models

import "time"

type FavourChoice struct {
	ID          string    `json:"id"`
	Label       string    `json:"label"`
	Description *string   `json:"description"`
	Costs       int       `json:"costs"`
	CanMessage  bool      `json:"can_message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FavourCount struct {
	Total     int `json:"total"`
	Remaining int `json:"remaining"`
}

type FavourRequest struct {
	ID                string     `json:"id"`
	ChoiceID          string     `json:"choice_id"`
	ChoiceLabel       string     `json:"choice_label"`
	ChoiceDescription *string    `json:"choice_description"`
	Message           *string    `json:"message"`
	RequestedAt       time.Time  `json:"requested_at"`
	FulfilledAt       *time.Time `json:"fulfilled_at"`
}
