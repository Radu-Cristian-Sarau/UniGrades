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

// Constants for the total ECTS per year bar chart.
const (
	// TotalECTSChartWidth is the width of the chart in characters
	TotalECTSChartWidth = 40
	// TotalECTSChartHeight is the height of the chart in characters
	TotalECTSChartHeight = 15
	// TotalECTSChartMax is the maximum value on the chart axis (75 gives room for 60 ECTS)
	TotalECTSChartMax = 75.0
)

// RenderTotalECTSPerYear displays a bar chart of total ECTS credits grouped by year.
// Shows both the visual representation and a summary of ECTS per year.
//
// Parameters:
//
//	uniColor: The university brand color for box styling
//	courses: The course documents to analyze
//
// Returns:
//
//	A formatted string with the chart and statistics
func RenderTotalECTSPerYear(uniColor lipgloss.Color, courses []bson.M) string {
	// Parse ECTS and years from courses
	ects, years := computations.ParseECTSAndYears(courses)

	// Calculate total ECTS for each year
	totalPerYear := computations.TotalECTSPerYear(ects, years)

	// Sort years for consistent display order
	sortedYears := make([]int, 0, len(totalPerYear))
	for y := range totalPerYear {
		sortedYears = append(sortedYears, y)
	}
	sort.Ints(sortedYears)

	// Build bar chart data
	barData := BuildBarDataPerYear(sortedYears, totalPerYear)

	// Create and render the bar chart
	bc := barchart.New(TotalECTSChartWidth, TotalECTSChartHeight,
		barchart.WithMaxValue(TotalECTSChartMax),
		barchart.WithNoAutoMaxValue(),
		barchart.WithStyles(BarAxisStyle, BarLabelStyle),
		barchart.WithDataSet(barData),
	)
	bc.Draw()

	// Build header with per-year totals listed
	header := "Total ECTS Per Year\n"
	for _, y := range sortedYears {
		if y > 1 {
			header += fmt.Sprintf("  Year %d: %.0f", y, totalPerYear[y])
		} else {
			header += fmt.Sprintf("Year %d: %.0f", y, totalPerYear[y])
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
