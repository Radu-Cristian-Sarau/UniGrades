package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"

	"UniGrades/internal/tui"
	"UniGrades/internal/university"
)

type Model struct {
	choices           []string
	cursor            int
	selected          map[int]struct{}
	tableStr          string
	avgStr            string
	avgPerYearStr     string
	avgECTSPerYearStr string
	ectsStr           string
	headers           []string
	courses           []bson.M
}

func InitialModel(headers []string, courses []bson.M) Model {
	tableStr := tui.RenderTable(tui.DefaultColor, headers, courses)
	avgStr := tui.RenderAverageGrades(tui.DefaultColor, courses)
	avgPerYearStr := tui.RenderAverageGradesPerYear(courses)
	avgECTSPerYearStr := tui.RenderTotalECTSPerYear(courses)
	ectsStr := tui.RenderECTS(tui.DefaultColor, courses)
	return Model{
		choices:           university.Names(),
		selected:          make(map[int]struct{}),
		tableStr:          tableStr,
		avgStr:            avgStr,
		avgPerYearStr:     avgPerYearStr,
		avgECTSPerYearStr: avgECTSPerYearStr,
		ectsStr:           ectsStr,
		headers:           headers,
		courses:           courses,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
				m.tableStr = tui.RenderTable(tui.DefaultColor, m.headers, m.courses)
				m.avgStr = tui.RenderAverageGrades(tui.DefaultColor, m.courses)
				m.ectsStr = tui.RenderECTS(tui.DefaultColor, m.courses)
			} else {
				m.selected = map[int]struct{}{m.cursor: {}}
				color := uniColors[m.choices[m.cursor]]
				m.tableStr = tui.RenderTable(color, m.headers, m.courses)
				m.avgStr = tui.RenderAverageGrades(color, m.courses)
				m.ectsStr = tui.RenderECTS(color, m.courses)
			}
		}
	}

	return m, nil
}

// SelectedUniversity returns the name of the selected university, or "" if none.
func (m Model) SelectedUniversity() string {
	for i := range m.selected {
		return m.choices[i]
	}
	return ""
}

var uniColors = university.ColorMap()

func (m Model) View() string {
	s := "\n\nSelect university: \n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		style := lipgloss.NewStyle().Foreground(uniColors[choice])
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, style.Render(choice))
	}
	s += "\n" + m.tableStr + "\n"
	s += "\n" + m.avgStr + "\n"
	s += "\n" + m.avgPerYearStr + "\n"
	s += "\n" + m.avgECTSPerYearStr + "\n"
	s += "\n" + m.ectsStr + "\n"

	s += "\nPress Ctrl + C to quit.\n"

	return s
}
