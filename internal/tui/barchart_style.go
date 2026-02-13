package tui

import (
	"fmt"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

var (
	BarAxisStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	BarLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

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
