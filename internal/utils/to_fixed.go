package utils

import "math"

// Rounds a float to 2 decimals
func ToFixed(num float64) float64 {
	return math.Round((num)*100) / 100
}
