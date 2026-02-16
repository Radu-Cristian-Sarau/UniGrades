// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// TUI libraries
	"github.com/charmbracelet/lipgloss"       // Styling and layout
	"github.com/charmbracelet/lipgloss/table" // Table component
)

var (
	// DefaultColor is white, used for default text color in tables
	DefaultColor = lipgloss.Color("#ffffff")
	// gray is used for even-row text in tables
	gray = lipgloss.Color("245")
	// lightGray is used for odd-row text in tables
	lightGray = lipgloss.Color("241")

	// cellStyle provides base padding for table cells
	cellStyle = lipgloss.NewStyle().Padding(0, 1)
	// OddRowStyle applies styling for odd-numbered rows
	OddRowStyle = cellStyle.Foreground(gray)
	// EvenRowStyle applies styling for even-numbered rows
	EvenRowStyle = cellStyle.Foreground(lightGray)
)

// TableStyleFunc returns a style function for tables based on row position.
// The function applies different styles to headers and alternating rows.
// uniColor is the university brand color used for header styling.
func TableStyleFunc(uniColor lipgloss.Color) func(row, col int) lipgloss.Style {
	headerStyle := lipgloss.NewStyle().Foreground(DefaultColor).Align(lipgloss.Center)
	return func(row, col int) lipgloss.Style {
		switch {
		case row == table.HeaderRow:
			// Header row uses centered white text
			return headerStyle
		case row%2 == 0:
			// Even rows use even row style
			return EvenRowStyle
		default:
			// Odd rows use odd row style
			return OddRowStyle
		}
	}
}
