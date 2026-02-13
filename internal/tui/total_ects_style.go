package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	ECTSMaxValue  = 180.0
	ECTSBarWidth  = 44.0
	ECTSBarHeight = 1
)

var RemainingColor = lipgloss.Color("238")

func ECTSBarStyle(uniColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(uniColor).
		Background(uniColor)
}

func ECTSRemainingStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(RemainingColor).
		Background(RemainingColor)
}
