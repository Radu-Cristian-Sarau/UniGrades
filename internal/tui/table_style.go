package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	DefaultColor = lipgloss.Color("#ffffff")
	gray         = lipgloss.Color("245")
	lightGray    = lipgloss.Color("241")

	cellStyle    = lipgloss.NewStyle().Padding(0, 1)
	OddRowStyle  = cellStyle.Foreground(gray)
	EvenRowStyle = cellStyle.Foreground(lightGray)
)

func TableStyleFunc(uniColor lipgloss.Color) func(row, col int) lipgloss.Style {
	headerStyle := lipgloss.NewStyle().Foreground(DefaultColor).Align(lipgloss.Center)
	return func(row, col int) lipgloss.Style {
		switch {
		case row == table.HeaderRow:
			return headerStyle
		case row%2 == 0:
			return EvenRowStyle
		default:
			return OddRowStyle
		}
	}
}
