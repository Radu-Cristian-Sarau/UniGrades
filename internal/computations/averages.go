// Package computations provides mathematical calculations for course grades and ECTS credits.
// It includes functions for parsing course data, calculating averages, and breaking down metrics by year.
package computations

import (
	// Standard library imports
	"fmt"     // Formatted conversion of values to strings
	"strconv" // String conversion utilities

	// MongoDB BSON types
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ParseGradesAndYears extracts grade and year data from a slice of course documents.
// It safely handles type conversions and skips courses with invalid data.
func ParseGradesAndYears(courses []bson.M) ([]float64, []int) {
	var grades []float64
	var years []int

	for _, course := range courses {
		// Convert Grade field to string, then parse as float64
		gradeStr := fmt.Sprintf("%v", course["Grade"])
		// Convert Year field to string, then parse as int
		yearStr := fmt.Sprintf("%v", course["Year"])

		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			continue // Skip courses with invalid grade
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			continue // Skip courses with invalid year
		}

		grades = append(grades, grade)
		years = append(years, year)
	}

	return grades, years
}

// Average calculates the simple arithmetic mean of a slice of numbers.
// Returns 0 if the slice is empty.
func Average(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

// WeightedAverage calculates the weighted arithmetic mean of two slices.
// Returns 0 if the slices are empty, have different lengths, or if weights sum to zero.
func WeightedAverage(nums []float64, weights []float64) float64 {
	if len(nums) == 0 || len(weights) == 0 || len(nums) != len(weights) {
		return 0
	}
	sum := 0.0
	weightSum := 0.0
	for i, num := range nums {
		sum += num * weights[i]
		weightSum += weights[i]
	}
	if weightSum == 0 {
		return 0
	}
	return sum / weightSum
}

// AverageGradePerYear groups grades by year and calculates the average for each year.
// Returns a map where keys are years and values are the average grade for that year.
func AverageGradePerYear(grades []float64, years []int) map[int]float64 {
	// Group grades by year
	yearGrades := make(map[int][]float64)
	for i, grade := range grades {
		year := years[i]
		yearGrades[year] = append(yearGrades[year], grade)
	}
	// Calculate average for each year
	avgPerYear := make(map[int]float64)
	for year, grades := range yearGrades {
		avgPerYear[year] = Average(grades)
	}
	return avgPerYear
}
