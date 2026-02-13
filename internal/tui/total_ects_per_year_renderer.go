package tui

import (
	"UniGrades/internal/computations"
	"fmt"
	"sort"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	TotalECTSChartWidth  = 40
	TotalECTSChartHeight = 15
	TotalECTSChartMax    = 75.0 // Max is actually 60, but this gives the bars some room
)

func RenderTotalECTSPerYear(courses []bson.M) string {
	ects, years := computations.ParseECTSAndYears(courses)

	totalPerYear := computations.TotalECTSPerYear(ects, years)

	// Sort years so bars appear in order
	sortedYears := make([]int, 0, len(totalPerYear))
	for y := range totalPerYear {
		sortedYears = append(sortedYears, y)
	}
	sort.Ints(sortedYears)

	// Build bar data for each year with fixed colors
	axisStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	var barData []barchart.BarData
	for i, y := range sortedYears {
		total := totalPerYear[y]
		barData = append(barData, barchart.BarData{
			Label: fmt.Sprintf("Y%d", y),
			Values: []barchart.BarValue{
				{Name: fmt.Sprintf("Year %d", y), Value: total, Style: barStyleForYear(i)},
			},
		})
	}

	bc := barchart.New(TotalECTSChartWidth, TotalECTSChartHeight,
		barchart.WithMaxValue(TotalECTSChartMax),
		barchart.WithNoAutoMaxValue(),
		barchart.WithStyles(axisStyle, labelStyle),
		barchart.WithDataSet(barData),
	)
	bc.Draw()

	// Build header with per-year totals
	header := "Total ECTS Per Year\n"
	for _, y := range sortedYears {
		header += fmt.Sprintf("  Year %d: %.0f", y, totalPerYear[y])
	}
	header += "\n"

	return header + bc.View()
}
