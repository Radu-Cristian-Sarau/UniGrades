// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// Standard library imports
	"fmt" // Formatted string creation

	// TUI libraries
	"github.com/NimbleMarkets/ntcharts/barchart" // Bar chart component
	"github.com/charmbracelet/lipgloss"          // Styling and layout
)

var (
	// BarAxisStyle is the grey color used for bar chart axes
	BarAxisStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	// BarLabelStyle is the grey color used for bar chart labels
	BarLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

// BuildBarDataPerYear constructs bar chart data from a list of years and their values.
// It creates one bar per year with labels formatted as "Y1", "Y2", etc.
//
// Parameters:
//
//	sortedYears: A sorted slice of year numbers
//	valuesPerYear: A map of year -> value to display
//
// Returns:
//
//	A slice of BarData ready for use with the barchart component
func BuildBarDataPerYear(sortedYears []int, valuesPerYear map[int]float64) []barchart.BarData {
	var barData []barchart.BarData
	for i, y := range sortedYears {
		barData = append(barData, barchart.BarData{
			Label: fmt.Sprintf("Y%d", y),
			Values: []barchart.BarValue{
				{Name: fmt.Sprintf("Year %d", y), Value: valuesPerYear[y], Style: barStyleForYear(i)},
			},
		})
	}
	return barData
}
