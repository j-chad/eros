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

// setDifference returns a set with the result of A - B
func setDifference[T comparable](a, b []T) []T {
	resultSet := make(map[T]struct{}, len(a))

	for _, v := range a {
		resultSet[v] = struct{}{}
	}

	for _, v := range b {
		delete(resultSet, v)
	}

	result := make([]T, 0, len(resultSet))
	for v := range resultSet {
		result = append(result, v)
	}

	return result
}
