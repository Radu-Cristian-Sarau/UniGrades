package tui

import (
	"fmt"
	"strings"

	"UniGrades/internal/computations"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RenderECTS(uniColor lipgloss.Color, courses []bson.M) string {
	totalECTS := computations.TotalECTS(computations.ParseECTS(courses))

	remaining := ECTSMaxValue - totalECTS
	if remaining < 0 {
		remaining = 0
	}

	d1 := barchart.BarData{
		Label: "Total ECTS Obtained",
		Values: []barchart.BarValue{
			{"ECTS", totalECTS, ECTSBarStyle(uniColor)},
			{"Remaining", remaining, ECTSRemainingStyle()},
		},
	}

	bc := barchart.New(ECTSBarWidth, ECTSBarHeight, barchart.WithHorizontalBars(), barchart.WithMaxValue(ECTSMaxValue), barchart.WithNoAxis())
	bc.PushAll([]barchart.BarData{d1})
	bc.Draw()

	scaleLine := buildScaleLine(totalECTS)

	header := fmt.Sprintf("Total ECTS: %.0f / %.0f", totalECTS, ECTSMaxValue)
	return header + "\n" + bc.View() + "\n" + scaleLine
}

func buildScaleLine(totalECTS float64) string {
	ectsPos := int(totalECTS / ECTSMaxValue * float64(ECTSBarWidth))
	if ectsPos < 0 {
		ectsPos = 0
	}
	if ectsPos > ECTSBarWidth {
		ectsPos = ECTSBarWidth
	}

	ectsLabel := fmt.Sprintf("%.0f", totalECTS)
	maxLabel := fmt.Sprintf("%.0f", ECTSMaxValue)

	scale := make([]byte, ECTSBarWidth)
	for i := range scale {
		scale[i] = ' '
	}

	if ectsPos+len(ectsLabel) <= ECTSBarWidth {
		copy(scale[ectsPos:], ectsLabel)
	}

	maxPos := ECTSBarWidth - len(maxLabel)
	if maxPos > ectsPos+len(ectsLabel)+1 {
		copy(scale[maxPos:], maxLabel)
	}

	if ectsPos > 2 {
		scale[0] = '0'
	}

	return strings.TrimRight(string(scale), " ")
}
