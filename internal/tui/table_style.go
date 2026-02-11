package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	Purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	HeaderStyle  = lipgloss.NewStyle().Foreground(Purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1)
	OddRowStyle  = cellStyle.Foreground(gray)
	EvenRowStyle = cellStyle.Foreground(lightGray)
)

func TableStyleFunc(row, col int) lipgloss.Style {
	switch {
	case row == table.HeaderRow:
		return HeaderStyle
	case row%2 == 0:
		return EvenRowStyle
	default:
		return OddRowStyle
	}
}
