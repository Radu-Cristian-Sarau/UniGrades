// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// Standard library imports
	"fmt"     // Formatted I/O and string conversion
	"sort"    // Sorting utilities
	"strconv" // String conversion utilities

	// TUI libraries
	"github.com/charmbracelet/lipgloss"       // Styling and layout
	"github.com/charmbracelet/lipgloss/table" // Table component

	// MongoDB types
	"go.mongodb.org/mongo-driver/v2/bson"
)

// sortCoursesByYear sorts a slice of course documents by their year field in ascending order.
// A copy is made to avoid modifying the original slice.
func sortCoursesByYear(courses []bson.M) []bson.M {
	sorted := make([]bson.M, len(courses))
	copy(sorted, courses)
	sort.Slice(sorted, func(i, j int) bool {
		yearI := toYear(sorted[i]["Year"])
		yearJ := toYear(sorted[j]["Year"])
		return yearI < yearJ
	})
	return sorted
}

// toYear converts various numeric types to int for year comparison.
// Handles int, int32, int64, float64, and string types.
func toYear(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		year, _ := strconv.Atoi(val)
		return year
	default:
		return 0
	}
}

// RenderTable creates a formatted table displaying courses with a border styled in the university color.
// Courses are sorted by year, and fields are displayed in the order of the provided headers.
//
// Parameters:
//
//	uniColor: The university brand color for table borders
//	headers: The column headers to display
//	courses: The course documents to display
//
// Returns:
//
//	A formatted table string
func RenderTable(uniColor lipgloss.Color, headers []string, courses []bson.M) string {
	// Sort courses by year so they appear in chronological order
	courses = sortCoursesByYear(courses)

	// Convert each course to a row of strings
	rows := make([][]string, 0, len(courses))
	for _, course := range courses {
		row := make([]string, 0, len(headers))
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", course[h]))
		}
		rows = append(rows, row)
	}

	// Create and configure the table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(TableStyleFunc(uniColor)).
		Headers(headers...).
		Rows(rows...)

	return t.Render()
}
