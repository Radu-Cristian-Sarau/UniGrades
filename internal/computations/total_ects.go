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

func TotalECTSPerYear(ects []float64, years []int) map[int]float64 {
	yearECTS := make(map[int][]float64)
	for i, ects := range ects {
		year := years[i]
		yearECTS[year] = append(yearECTS[year], ects)
	}
	totalPerYear := make(map[int]float64)
	for year, ects := range yearECTS {
		sum := 0.0
		for _, e := range ects {
			sum += e
		}
		totalPerYear[year] = sum
	}
	return totalPerYear
}
