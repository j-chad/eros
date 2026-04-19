package service

import "math"

// haversineDistance calculates the distance in metres between two points specified by their latitude and longitude using the Haversine formula.
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // metres

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	return earthRadius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// setDifferenceBy returns elements in A whose key (extracted by keyFn) is not
// present in B. When the same key appears in both, the element from A is kept
// with its latest values from A.
func setDifferenceBy[T any, K comparable](a, b []T, keyFn func(T) K) []T {
	bKeys := make(map[K]struct{}, len(b))
	for _, v := range b {
		bKeys[keyFn(v)] = struct{}{}
	}

	result := make([]T, 0)
	for _, v := range a {
		if _, ok := bKeys[keyFn(v)]; !ok {
			result = append(result, v)
		}
	}

	return result
}
