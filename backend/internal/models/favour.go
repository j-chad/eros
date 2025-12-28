package models

import "time"

type FavourChoice struct {
	ID          string    `json:"id"`
	Label       string    `json:"label"`
	Description *string   `json:"description"`
	CanMessage  bool      `json:"can_message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
