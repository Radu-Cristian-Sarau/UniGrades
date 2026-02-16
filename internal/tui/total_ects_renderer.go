// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// Internal packages
	"UniGrades/internal/computations"
	// Standard library imports
	"fmt"     // Formatted I/O
	"strings" // String manipulation

	// TUI libraries
	"github.com/NimbleMarkets/ntcharts/barchart" // Bar chart component
	"github.com/charmbracelet/lipgloss"          // Styling and layout

	// MongoDB types
	"go.mongodb.org/mongo-driver/v2/bson"
)

// RenderECTS displays a horizontal progress bar showing earned vs remaining ECTS credits.
// Shows current progress toward the 180 ECTS degree requirement.
//
// Parameters:
//
//	uniColor: The university brand color for bar styling
//	courses: The course documents to analyze
//
// Returns:
//
//	A formatted string with the ECTS progress bar and scale
func RenderECTS(uniColor lipgloss.Color, courses []bson.M) string {
	// Calculate total earned ECTS
	totalECTS := computations.TotalECTS(computations.ParseECTS(courses))

	// Calculate remaining ECTS to reach 180 (capped at 0 if exceeded)
	remaining := ECTSMaxValue - totalECTS
	if remaining < 0 {
		remaining = 0
	}

	// Create bar chart data with earned vs remaining segments
	d1 := barchart.BarData{
		Label: "ECTS",
		Values: []barchart.BarValue{
			{"ECTS", totalECTS, ECTSBarStyle(uniColor)},
			{"Remaining", remaining, ECTSRemainingStyle()},
		},
	}

	// Create and render the horizontal bar chart
	bc := barchart.New(ECTSBarWidth, ECTSBarHeight, barchart.WithHorizontalBars(), barchart.WithMaxValue(ECTSMaxValue), barchart.WithNoAxis())
	bc.PushAll([]barchart.BarData{d1})
	bc.Draw()

	// Build the display with header, chart, and scale line
	header := fmt.Sprintf("Total ECTS")
	scaleLine := buildScaleLine(totalECTS)
	content := header + "\n" + bc.View() + "\n" + scaleLine

	// Style in a bordered box
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1)

	return box.Render(content)
}

// buildScaleLine constructs a scale line showing the current ECTS position (0, current value, 180).
// Positions labels proportionally along the bar width.
func buildScaleLine(totalECTS float64) string {
	// Calculate position of current ECTS on the scale
	ectsPos := int(totalECTS / ECTSMaxValue * float64(ECTSBarWidth))
	if ectsPos < 0 {
		ectsPos = 0
	}
	if ectsPos > ECTSBarWidth {
		ectsPos = ECTSBarWidth
	}

	// Format value labels
	ectsLabel := fmt.Sprintf("%.0f", totalECTS)
	maxLabel := fmt.Sprintf("%.0f", ECTSMaxValue)

	// Initialize scale with spaces
	scale := make([]byte, ECTSBarWidth)
	for i := range scale {
		scale[i] = ' '
	}

	// Place "0" at the start if there's room
	if ectsPos > 2 {
		scale[0] = '0'
	}

	// Place current ECTS value at its position
	if ectsPos+len(ectsLabel) <= ECTSBarWidth {
		copy(scale[ectsPos:], ectsLabel)
	}

	// Place max value at the end
	maxPos := ECTSBarWidth - len(maxLabel)
	if maxPos > ectsPos+len(ectsLabel)+1 {
		copy(scale[maxPos:], maxLabel)
	}

	return strings.TrimRight(string(scale), " ")
}
