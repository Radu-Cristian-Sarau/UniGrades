package tui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func sortCoursesByYear(courses []bson.M) []bson.M {
	sorted := make([]bson.M, len(courses))
	copy(sorted, courses)
	sort.Slice(sorted, func(i, j int) bool {
		yearI := toYear(sorted[i]["Year"])
		yearJ := toYear(sorted[j]["Year"])
		return yearI < yearJ
	})
	return sorted
}

func toYear(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		year, _ := strconv.Atoi(val)
		return year
	default:
		return 0
	}
}

func RenderTable(uniColor lipgloss.Color, headers []string, courses []bson.M) string {
	courses = sortCoursesByYear(courses)

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
