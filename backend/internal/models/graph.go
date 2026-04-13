package models

import (
	"backend/internal"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type NodeType string

const (
	StartNode        NodeType = "start"
	ManualNode       NodeType = "manual"
	CodeGateNode     NodeType = "code"
	LocationGateNode NodeType = "location"
	RewardNode       NodeType = "reward"
)

var nodeDataDecoders = map[NodeType]func(json.RawMessage) (any, error){
	LocationGateNode: internal.DecodeInto[LocationData],
	CodeGateNode:     internal.DecodeInto[CodeData],
	ManualNode:       internal.DecodeInto[ManualData],
	RewardNode:       internal.DecodeInto[RewardData],
}

type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Node struct {
	ID      string   `json:"id"`
	GraphID string   `json:"-"`
	Type    NodeType `json:"type"`

	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`

	UIPosition *NodePosition `json:"ui_position,omitempty"`

	Data any `json:"data,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}

// UnmarshalJSON implements custom unmarshaling for Node to handle dynamic Data field.
// It uses the nodeDataDecoders map to decode the Data field based on the Node Type.
func (n *Node) UnmarshalJSON(raw []byte) error {
	type nodeAlias Node

	aux := struct {
		*nodeAlias
		Data json.RawMessage `json:"data,omitempty"`
	}{
		nodeAlias: (*nodeAlias)(n),
	}

	if err := json.Unmarshal(raw, &aux); err != nil {
		return err
	}

	// If no data (or null/empty), keep nil
	if len(bytes.TrimSpace(aux.Data)) == 0 || bytes.Equal(bytes.TrimSpace(aux.Data), []byte("null")) {
		n.Data = nil
		return nil
	}

	decoder, ok := nodeDataDecoders[n.Type]
	if !ok {
		return nil
	}

	decodedData, err := decoder(aux.Data)
	if err != nil {
		return fmt.Errorf("failed to decode node data for type %s: %w", n.Type, err)
	}

	n.Data = decodedData
	return nil
}

type NewGraphRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartingAt  time.Time `json:"starting_at"`
}

type Edge struct {
	ID      string `json:"id"`
	GraphID string `json:"-"`
	From    string `json:"from"`
	To      string `json:"to"`

	ChoiceLabel string `json:"choice_label,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Viewport struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}

type GraphStatus string

const (
	GraphNotStarted GraphStatus = "not_started"
	GraphInProgress GraphStatus = "in_progress"
	GraphCompleted  GraphStatus = "completed"
)

type Graph struct {
	ID string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description,omitempty"`

	Status GraphStatus `json:"status,omitempty"`

	StartingAt time.Time `json:"starting_at"`

	Viewport *Viewport `json:"viewport,omitempty"`

	Nodes *[]Node `json:"nodes,omitempty"`
	Edges *[]Edge `json:"edges,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UnlockResult struct {
	UnlockedNode Node   `json:"unlocked_node"`
	NewNodes     []Node `json:"new_nodes"`
	NewEdges     []Edge `json:"new_edges"`
}
