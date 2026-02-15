package grades

import "github.com/charmbracelet/lipgloss"

// Color constants for data screen styling
const (
	// Text colors
	ColorDimText         = lipgloss.Color("243") // Dim grey for unavailable message
	ColorTableAlternate1 = lipgloss.Color("245") // Light grey for alternating table rows
	ColorTableAlternate2 = lipgloss.Color("241") // Dark grey for alternating table rows
	ColorSuccess         = lipgloss.Color("42")  // Green for success messages
)
