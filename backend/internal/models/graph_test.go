package models

import (
	"backend/internal/testutil"
	"encoding/json"
	"strings"
	"testing"
)

func TestNodeUnmarshalJSON_CodeGate(t *testing.T) {
	raw := `{
		"id": "n1",
		"type": "code",
		"title": "Enter Code",
		"data": {"codes": ["SECRET", "PASSWORD"]},
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))

	data := testutil.TypeAssert[CodeData](t, node.Data)
	testutil.ElementsMatch(t, data.Codes, []string{"SECRET", "PASSWORD"})
}

func TestNodeUnmarshalJSON_LocationGate(t *testing.T) {
	raw := `{
		"id": "n2",
		"type": "location",
		"title": "Find Me",
		"data": {"latitude": -36.85, "longitude": 174.76, "radius_m": 50},
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))

	data := testutil.TypeAssert[LocationData](t, node.Data)
	testutil.Equal(t, data.Latitude, -36.85)
	testutil.Equal(t, data.RadiusM, 50)
}

func TestNodeUnmarshalJSON_ManualGate(t *testing.T) {
	raw := `{
		"id": "n3",
		"type": "manual",
		"title": "Wait",
		"data": {"instructions": "ask nicely"},
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))

	data := testutil.TypeAssert[ManualData](t, node.Data)
	testutil.Equal(t, data.Instructions, "ask nicely")
	testutil.Nil(t, data.UnlockedAt)
}

func TestNodeUnmarshalJSON_RewardNode(t *testing.T) {
	raw := `{
		"id": "n4",
		"type": "reward",
		"title": "Prize",
		"data": {"reward_type": "image", "payload": "img.jpg", "give_favours": 3},
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))

	data := testutil.TypeAssert[RewardData](t, node.Data)
	testutil.Equal(t, data.RewardType, "image")
	testutil.Equal(t, data.GiveFavours, 3)
}

func TestNodeUnmarshalJSON_NullData(t *testing.T) {
	raw := `{
		"id": "n5", "type": "start", "title": "Start",
		"data": null,
		"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))
	testutil.Nil(t, node.Data)
}

func TestNodeUnmarshalJSON_MissingData(t *testing.T) {
	raw := `{
		"id": "n6", "type": "start", "title": "Start",
		"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))
	testutil.Nil(t, node.Data)
}

func TestNodeUnmarshalJSON_StartNodeWithData(t *testing.T) {
	raw := `{
		"id": "n7", "type": "start", "title": "Start",
		"data": {"unexpected": "field"},
		"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))
	testutil.Nil(t, node.Data)
}

func TestNodeUnmarshalJSON_InvalidJSON(t *testing.T) {
	var node Node
	testutil.NotNilErr(t, json.Unmarshal([]byte(`not json`), &node))
}

func TestNodeUnmarshalJSON_InvalidDataForType(t *testing.T) {
	raw := `{
		"id": "n8", "type": "location", "title": "Bad",
		"data": {"latitude": "not a number"},
		"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"
	}`

	var node Node
	err := json.Unmarshal([]byte(raw), &node)
	testutil.NotNilErr(t, err)
	testutil.True(t, strings.Contains(err.Error(), "failed to decode node data"),
		"error should mention 'failed to decode node data'")
}

func TestNodeUnmarshalJSON_PreservesFields(t *testing.T) {
	raw := `{
		"id": "abc", "type": "code", "title": "My Node",
		"description": "desc", "data": {"code": "X"},
		"created_at": "2024-06-15T10:30:00Z", "updated_at": "2024-06-15T11:00:00Z"
	}`

	var node Node
	testutil.NilErr(t, json.Unmarshal([]byte(raw), &node))

	testutil.Equal(t, node.ID, "abc")
	testutil.Equal(t, node.Title, "My Node")
	testutil.NotNil(t, node.Description)
	testutil.Equal(t, *node.Description, "desc")
	testutil.Equal(t, node.Type, CodeGateNode)
}
