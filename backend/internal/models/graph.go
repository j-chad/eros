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
	ID   string   `json:"id"`
	Type NodeType `json:"type"`

	Title       string `json:"title"`
	Description string `json:"description"`

	// Only one of the following will be populated based on the NodeType
	Start    *StartData    `json:"start,omitempty"`
	Location *LocationData `json:"location,omitempty"`
	Code     *CodeData     `json:"code,omitempty"`
	Reward   *RewardData   `json:"reward,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}

type NodeWithEdges struct {
	Node  Node   `json:"node"`
	Edges []Edge `json:"edges"`
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

type Graph struct {
	StartNode NodeWithEdges `json:"start_node"`
	Nodes     []Node        `json:"nodes"`
	Edges     []Edge        `json:"edges"`
}
