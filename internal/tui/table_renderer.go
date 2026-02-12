package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RenderTable(uniColor lipgloss.Color, headers []string, courses []bson.M) string {

	rows := make([][]string, 0, len(courses))
	for _, course := range courses {
		row := make([]string, 0, len(headers))
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", course[h]))
		}
		rows = append(rows, row)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(TableStyleFunc(uniColor)).
		Headers(headers...).
		Rows(rows...)

	// You can also add tables row-by-row
	// t.Row("English", "You look absolutely fabulous.", "How's it going?")

	return t.Render()
}
