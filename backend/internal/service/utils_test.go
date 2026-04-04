package service

import (
	"backend/internal/testutil"
	"fmt"
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name           string
		lat1, lon1     float64
		lat2, lon2     float64
		expectedMetres float64
		tolerance      float64
	}{
		{
			name: "same point",
			lat1: 0, lon1: 0,
			lat2: 0, lon2: 0,
			expectedMetres: 0,
			tolerance:      0,
		},
		{
			name: "London to Paris",
			lat1: 51.5074, lon1: -0.1278,
			lat2: 48.8566, lon2: 2.3522,
			expectedMetres: 343_556,
			tolerance:      500,
		},
		{
			name: "New York to Los Angeles",
			lat1: 40.7128, lon1: -74.0060,
			lat2: 34.0522, lon2: -118.2437,
			expectedMetres: 3_935_746,
			tolerance:      5000,
		},
		{
			name: "antipodal points",
			lat1: 0, lon1: 0,
			lat2: 0, lon2: 180,
			expectedMetres: math.Pi * 6371000,
			tolerance:      1,
		},
		{
			name: "short distance ~100m",
			lat1: -36.8485, lon1: 174.7633,
			lat2: -36.8476, lon2: 174.7633,
			expectedMetres: 100,
			tolerance:      5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := haversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			testutil.True(t, math.Abs(got-tt.expectedMetres) <= tt.tolerance,
				"haversineDistance() = "+fmt.Sprintf("%.1f m, want %.1f m (±%.0f)", got, tt.expectedMetres, tt.tolerance))
		})
	}
}

func TestHaversineDistance_Symmetry(t *testing.T) {
	d1 := haversineDistance(51.5074, -0.1278, 48.8566, 2.3522)
	d2 := haversineDistance(48.8566, 2.3522, 51.5074, -0.1278)
	testutil.Equal(t, d1, d2)
}
