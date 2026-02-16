// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// TUI libraries
	"github.com/charmbracelet/lipgloss"
)

// ECTS constants for the total ECTS bar visualization.
const (
	// ECTSMaxValue is the maximum ECTS value (180 credits for most bachelor's programs)
	ECTSMaxValue = 180.0
	// ECTSBarWidth is the width of the horizontal ECTS progress bar in characters
	ECTSBarWidth = 44.0
	// ECTSBarHeight is the height of the ECTS bar (always 1 for horizontal bars)
	ECTSBarHeight = 1
)

// RemainingColor is the grey color used for the "remaining ECTS" portion of the bar
var RemainingColor = lipgloss.Color("238")

// ECTSBarStyle returns a styled component for the earned ECTS part of the bar.
// Colors it with the university's brand color.
func ECTSBarStyle(uniColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(uniColor).
		Background(uniColor)
}

// ECTSRemainingStyle returns a styled component for the remaining ECTS part of the bar.
// Colors it with a dark grey to indicate unfilled progress.
func ECTSRemainingStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(RemainingColor).
		Background(RemainingColor)
}
