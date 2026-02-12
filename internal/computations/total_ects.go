package computations

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ParseECTS(courses []bson.M) []float64 {
	var ects []float64
	for _, course := range courses {
		ectsStr := fmt.Sprintf("%v", course["ECTS"])
		val, err := strconv.ParseFloat(ectsStr, 64)
		if err != nil {
			continue
		}
		ects = append(ects, val)
	}
	return ects
}

func TotalECTS(ects []float64) float64 {
	totalECTS := 0.0
	for _, ectsValue := range ects {
		totalECTS += ectsValue
	}
	return totalECTS
}
