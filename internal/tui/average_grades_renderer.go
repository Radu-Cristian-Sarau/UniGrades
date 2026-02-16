// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// Internal packages
	"UniGrades/internal/computations"
	// Standard library imports
	"fmt"     // Formatted I/O
	"strconv" // String conversion

	// TUI libraries
	"github.com/charmbracelet/lipgloss"       // Styling and layout
	"github.com/charmbracelet/lipgloss/table" // Table component

	// MongoDB types
	"go.mongodb.org/mongo-driver/v2/bson"
)

// RenderAverageGrades displays overall grade statistics.
// Shows both simple average and ECTS-weighted average grade.
//
// Parameters:
//
//	uniColor: The university brand color for table styling
//	courses: The course documents to analyze
//
// Returns:
//
//	A formatted table string with average grade metrics
func RenderAverageGrades(uniColor lipgloss.Color, courses []bson.M) string {
	var grades []float64
	var ects []float64

	// Extract grades and ECTS from all courses
	for _, course := range courses {
		gradeStr := fmt.Sprintf("%v", course["Grade"])
		ectsStr := fmt.Sprintf("%v", course["ECTS"])

		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			continue // Skip courses with invalid grade
		}
		credit, err := strconv.ParseFloat(ectsStr, 64)
		if err != nil {
			continue // Skip courses with invalid ECTS
		}

		grades = append(grades, grade)
		ects = append(ects, credit)
	}

	// Calculate both averages
	avg := computations.Average(grades)
	weightedAvg := computations.WeightedAverage(grades, ects)

	// Create table with results
	headers := []string{"Metric", "Value"}
	rows := [][]string{
		{"Average Grade", fmt.Sprintf("%.2f", avg)},
		{"Weighted Average (ECTS)", fmt.Sprintf("%.2f", weightedAvg)},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(TableStyleFunc(uniColor)).
		Headers(headers...).
		Rows(rows...)

	return t.Render()
}
