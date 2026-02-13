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
	AvgGradesChartWidth  = 40
	AvgGradesChartHeight = 15
	AvgGradesMaxValue    = 10.0
)

var yearBarColors = []lipgloss.Color{
	lipgloss.Color("#1a80bb"), // blue
	lipgloss.Color("#ea801c"), // orange
	lipgloss.Color("#17b118"), // green
}

func barStyleForYear(index int) lipgloss.Style {
	c := yearBarColors[index%len(yearBarColors)]
	return lipgloss.NewStyle().Foreground(c).Background(c)
}

func RenderAverageGradesPerYear(uniColor lipgloss.Color, courses []bson.M) string {
	grades, years := computations.ParseGradesAndYears(courses)

	avgPerYear := computations.AverageGradePerYear(grades, years)

	// Sort years so bars appear in order
	sortedYears := make([]int, 0, len(avgPerYear))
	for y := range avgPerYear {
		sortedYears = append(sortedYears, y)
	}
	sort.Ints(sortedYears)

	barData := BuildBarDataPerYear(sortedYears, avgPerYear)

	bc := barchart.New(AvgGradesChartWidth, AvgGradesChartHeight,
		barchart.WithMaxValue(AvgGradesMaxValue),
		barchart.WithNoAutoMaxValue(),
		barchart.WithStyles(BarAxisStyle, BarLabelStyle),
		barchart.WithDataSet(barData),
	)
	bc.Draw()

	// Build header with per-year averages
	header := "Average Grades Per Year\n"
	for _, y := range sortedYears {
		if y > 1 {
			header += fmt.Sprintf("  Year %d: %.2f", y, avgPerYear[y])
		} else {
			header += fmt.Sprintf("Year %d: %.2f", y, avgPerYear[y])
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
