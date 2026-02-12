package tui

import (
	"UniGrades/internal/computations"
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RenderAverageGrades(uniColor lipgloss.Color, courses []bson.M) string {
	var grades []float64
	var ects []float64

	for _, course := range courses {
		gradeStr := fmt.Sprintf("%v", course["Grade"])
		ectsStr := fmt.Sprintf("%v", course["ECTS"])

		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			continue
		}
		credit, err := strconv.ParseFloat(ectsStr, 64)
		if err != nil {
			continue
		}

		grades = append(grades, grade)
		ects = append(ects, credit)
	}

	avg := computations.Average(grades)
	weightedAvg := computations.WeightedAverage(grades, ects)

	headers := []string{"Metric", "Value"}
	rows := [][]string{
		{"Average Grade", fmt.Sprintf("%.2f", avg)},
		{"Weighted Average (ECTS)", fmt.Sprintf("%.2f", weightedAvg)},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(TableStyleFunc(uniColor)).
		Headers(headers...).
		Rows(rows...)

	return t.Render()
}
