package models

import "time"

type NodeType string

const (
	StartNode        NodeType = "start"
	ChoiceNode       NodeType = "choice"
	CodeGateNode     NodeType = "code_gate"
	LocationGateNode NodeType = "location_gate"
	RewardNode       NodeType = "reward"
)

type Node struct {
	ID   int64    `json:"id"`
	Type NodeType `json:"type"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Edges []Edge `json:"edges,omitempty"`

	// Only one of the following will be populated based on the NodeType
	Location *LocationData `json:"location,omitempty"`
	Code     *CodeData     `json:"code,omitempty"`
	Reward   *RewardData   `json:"reward,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}

type Edge struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`

	ChoiceLabel string `json:"choice_label,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
