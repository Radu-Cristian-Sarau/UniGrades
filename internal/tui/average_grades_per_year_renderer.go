// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// Internal packages
	"UniGrades/internal/computations"
	// Standard library imports
	"fmt"  // Formatted I/O
	"sort" // Sorting utilities

	// TUI libraries
	"github.com/NimbleMarkets/ntcharts/barchart" // Bar chart component
	"github.com/charmbracelet/lipgloss"          // Styling and layout

	// MongoDB types
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Constants for the average grades per year bar chart.
const (
	// AvgGradesChartWidth is the width of the chart in characters
	AvgGradesChartWidth = 40
	// AvgGradesChartHeight is the height of the chart in characters
	AvgGradesChartHeight = 15
	// AvgGradesMaxValue is the maximum value on the chart axis (10 for grades)
	AvgGradesMaxValue = 10.0
)

var (
	// yearBarColors defines the colors cycled through for each year's bar
	yearBarColors = []lipgloss.Color{
		lipgloss.Color("#1a80bb"), // Blue
		lipgloss.Color("#ea801c"), // Orange
		lipgloss.Color("#17b118"), // Green
	}
)

// barStyleForYear returns a color style for a bar at a given index.
// Cycles through predefined colors for each year.
func barStyleForYear(index int) lipgloss.Style {
	c := yearBarColors[index%len(yearBarColors)]
	return lipgloss.NewStyle().Foreground(c).Background(c)
}

// RenderAverageGradesPerYear displays a bar chart of average grades grouped by year.
// Shows both the visual representation and a summary below.
//
// Parameters:
//
//	uniColor: The university brand color for box styling
//	courses: The course documents to analyze
//
// Returns:
//
//	A formatted string with the chart and statistics
func RenderAverageGradesPerYear(uniColor lipgloss.Color, courses []bson.M) string {
	// Parse grades and years from courses
	grades, years := computations.ParseGradesAndYears(courses)

	// Calculate average grade for each year
	avgPerYear := computations.AverageGradePerYear(grades, years)

	// Sort years for consistent display order
	sortedYears := make([]int, 0, len(avgPerYear))
	for y := range avgPerYear {
		sortedYears = append(sortedYears, y)
	}
	sort.Ints(sortedYears)

	// Build bar chart data
	barData := BuildBarDataPerYear(sortedYears, avgPerYear)

	// Create and render the bar chart
	bc := barchart.New(AvgGradesChartWidth, AvgGradesChartHeight,
		barchart.WithMaxValue(AvgGradesMaxValue),
		barchart.WithNoAutoMaxValue(),
		barchart.WithStyles(BarAxisStyle, BarLabelStyle),
		barchart.WithDataSet(barData),
	)
	bc.Draw()

	// Build header with per-year averages listed
	header := "Average Grades Per Year\n"
	for _, y := range sortedYears {
		if y > 1 {
			header += fmt.Sprintf("  Year %d: %.2f", y, avgPerYear[y])
		} else {
			header += fmt.Sprintf("Year %d: %.2f", y, avgPerYear[y])
		}
	}
	header += "\n"

	// Combine header and chart
	content := header + bc.View()

	// Style the content in a bordered box
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1)

	return box.Render(content)
}
