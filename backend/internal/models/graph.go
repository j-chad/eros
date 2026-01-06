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

type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Node struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type"`

	Title       string `json:"title"`
	Description string `json:"description"`

	UIPosition *NodePosition `json:"ui_position,omitempty"`

	Data any `json:"data"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}

type NewGraphRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartingAt  time.Time `json:"starting_at"`
}

type Edge struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`

	ChoiceLabel string `json:"choice_label,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Viewport struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}

type Graph struct {
	ID string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	StartingAt time.Time `json:"starting_at"`

	Viewport *Viewport `json:"viewport,omitempty"`

	Nodes *[]Node `json:"nodes,omitempty"`
	Edges *[]Edge `json:"edges,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
