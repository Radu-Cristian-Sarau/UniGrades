// Package computations provides mathematical calculations for course grades and ECTS credits.
package computations

import (
	// Standard library imports
	"fmt"     // Formatted conversion of values to strings
	"strconv" // String conversion utilities

	// MongoDB BSON types
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ParseECTS extracts ECTS (European Credit Transfer System) credits from a slice of courses.
// It safely handles type conversions and skips courses with invalid data.
func ParseECTS(courses []bson.M) []float64 {
	var ects []float64
	for _, course := range courses {
		// Convert ECTS field to string, then parse as float64
		ectsStr := fmt.Sprintf("%v", course["ECTS"])
		val, err := strconv.ParseFloat(ectsStr, 64)
		if err != nil {
			continue // Skip courses with invalid ECTS
		}
		ects = append(ects, val)
	}
	return ects
}

// TotalECTS calculates the sum of all ECTS credits.
func TotalECTS(ects []float64) float64 {
	totalECTS := 0.0
	for _, ectsValue := range ects {
		totalECTS += ectsValue
	}
	return totalECTS
}

// ParseECTSAndYears extracts ECTS credits and corresponding years from course documents.
// It safely handles type conversions and skips courses with invalid data.
func ParseECTSAndYears(courses []bson.M) ([]float64, []int) {
	var ects []float64
	var years []int

	for _, course := range courses {
		// Convert ECTS field to string, then parse as float64
		ectsStr := fmt.Sprintf("%v", course["ECTS"])
		// Convert Year field to string, then parse as int
		yearStr := fmt.Sprintf("%v", course["Year"])

		e, err := strconv.ParseFloat(ectsStr, 64)
		if err != nil {
			continue // Skip courses with invalid ECTS
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			continue // Skip courses with invalid year
		}

		ects = append(ects, e)
		years = append(years, year)
	}

	return ects, years
}

// TotalECTSPerYear groups ECTS credits by year and calculates the total for each year.
// Returns a map where keys are years and values are the total ECTS for that year.
func TotalECTSPerYear(ects []float64, years []int) map[int]float64 {
	// Group ECTS by year
	yearECTS := make(map[int][]float64)
	for i, ects := range ects {
		year := years[i]
		yearECTS[year] = append(yearECTS[year], ects)
	}
	// Calculate total for each year
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
