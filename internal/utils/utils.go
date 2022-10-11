package utils

import (
	"math"
)

func ConvertFloatToKopecks(value float64) int64 {
	integer, fraction := math.Modf(value)

	return int64(integer*100) + int64(fraction*100)
}

func ConvertKopecksToFloat(kopecks int64) float64 {
	return float64(kopecks) / 100
}
