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

func RenderTotalECTSPerYear(uniColor lipgloss.Color, courses []bson.M) string {
	ects, years := computations.ParseECTSAndYears(courses)

	totalPerYear := computations.TotalECTSPerYear(ects, years)

	// Sort years so bars appear in order
	sortedYears := make([]int, 0, len(totalPerYear))
	for y := range totalPerYear {
		sortedYears = append(sortedYears, y)
	}
	sort.Ints(sortedYears)

	barData := BuildBarDataPerYear(sortedYears, totalPerYear)

	bc := barchart.New(TotalECTSChartWidth, TotalECTSChartHeight,
		barchart.WithMaxValue(TotalECTSChartMax),
		barchart.WithNoAutoMaxValue(),
		barchart.WithStyles(BarAxisStyle, BarLabelStyle),
		barchart.WithDataSet(barData),
	)
	bc.Draw()

	// Build header with per-year totals
	header := "Total ECTS Per Year\n"
	for _, y := range sortedYears {
		if y > 1 {
			header += fmt.Sprintf("  Year %d: %.0f", y, totalPerYear[y])
		} else {
			header += fmt.Sprintf("Year %d: %.0f", y, totalPerYear[y])
		}

	}
	header += "\n"

	content := header + bc.View()

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1)

	return box.Render(content)
}
