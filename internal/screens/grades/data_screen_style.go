// Package grades provides the data/grades screen for displaying course information and statistics.
package grades

import "github.com/charmbracelet/lipgloss"

// Color constants for data screen styling.
const (
	// ColorDimText is a dim grey color used for unavailable message text
	ColorDimText = lipgloss.Color("243")
	// ColorTableAlternate1 is light grey for alternating table rows (even rows)
	ColorTableAlternate1 = lipgloss.Color("245")
	// ColorTableAlternate2 is dark grey for alternating table rows (odd rows)
	ColorTableAlternate2 = lipgloss.Color("241")
	// ColorSuccess is green used for success messages
	ColorSuccess = lipgloss.Color("42")
)
