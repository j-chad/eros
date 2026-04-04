package service

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"testing"
	"time"
)

func TestValidateUnlockPayload_CodeGate_Correct(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: models.CodeData{Code: "SECRET"},
	}
	testutil.NilErr(t, validateUnlockPayload(node, "SECRET"))
}

func TestValidateUnlockPayload_CodeGate_Wrong(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: models.CodeData{Code: "SECRET"},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, "WRONG"), NodeUnlockIncorrect)
}

func TestValidateUnlockPayload_CodeGate_Empty(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: models.CodeData{Code: "SECRET"},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, ""), NodeUnlockIncorrect)
}

func TestValidateUnlockPayload_CodeGate_CaseSensitive(t *testing.T) {
	node := &models.Node{
		ID:   "n1",
		Type: models.CodeGateNode,
		Data: models.CodeData{Code: "Secret"},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, "secret"), NodeUnlockIncorrect)
}

func TestValidateUnlockPayload_LocationGate_WithinRadius(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: models.LocationData{Latitude: -36.8485, Longitude: 174.7633, RadiusM: 100},
	}
	testutil.NilErr(t, validateUnlockPayload(node, "-36.8484,174.7633"))
}

func TestValidateUnlockPayload_LocationGate_OutsideRadius(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: models.LocationData{Latitude: -36.8485, Longitude: 174.7633, RadiusM: 10},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, "-36.8400,174.7633"), NodeUnlockIncorrect)
}

func TestValidateUnlockPayload_LocationGate_DefaultRadius(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: models.LocationData{Latitude: -36.8485, Longitude: 174.7633, RadiusM: 0},
	}
	// ~5m away — within default 10m
	testutil.NilErr(t, validateUnlockPayload(node, "-36.84846,174.7633"))
}

func TestValidateUnlockPayload_LocationGate_InvalidPayload(t *testing.T) {
	node := &models.Node{
		ID:   "n2",
		Type: models.LocationGateNode,
		Data: models.LocationData{Latitude: 0, Longitude: 0, RadiusM: 10},
	}

	tests := []struct {
		name    string
		payload string
	}{
		{"empty", ""},
		{"no comma", "123.456"},
		{"letters", "abc,def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUnlockPayload(node, tt.payload)
			testutil.NotNilErr(t, err)
			testutil.ErrorIsNot(t, err, NodeUnlockIncorrect)
		})
	}
}

func TestValidateUnlockPayload_ManualGate_Unlocked(t *testing.T) {
	now := time.Now()
	node := &models.Node{
		ID:   "n3",
		Type: models.ManualNode,
		Data: models.ManualData{Instructions: "do the thing", UnlockedAt: &now},
	}
	testutil.NilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_ManualGate_NotUnlocked(t *testing.T) {
	node := &models.Node{
		ID:   "n3",
		Type: models.ManualNode,
		Data: models.ManualData{Instructions: "do the thing", UnlockedAt: nil},
	}
	testutil.ErrorIs(t, validateUnlockPayload(node, ""), NodeUnlockIncorrect)
}

func TestValidateUnlockPayload_UnsupportedType(t *testing.T) {
	node := &models.Node{ID: "n4", Type: models.StartNode}
	testutil.NotNilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_RewardType(t *testing.T) {
	node := &models.Node{ID: "n5", Type: models.RewardNode}
	testutil.NotNilErr(t, validateUnlockPayload(node, ""))
}

func TestValidateUnlockPayload_WrongDataType(t *testing.T) {
	node := &models.Node{
		ID:   "n6",
		Type: models.CodeGateNode,
		Data: models.LocationData{Latitude: 0, Longitude: 0},
	}
	testutil.NotNilErr(t, validateUnlockPayload(node, "anything"))
}
