package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"

	"UniGrades/internal/tui"
	"UniGrades/internal/university"
)

type Screen int

const (
	PickerScreen Screen = iota
	DataScreen
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
	termWidth         int
	termHeight        int
	screen            Screen
}

func InitialModel(headers []string, courses []bson.M) Model {
	tableStr := tui.RenderTable(tui.DefaultColor, headers, courses)
	avgStr := tui.RenderAverageGrades(tui.DefaultColor, courses)
	avgPerYearStr := tui.RenderAverageGradesPerYear(tui.DefaultColor, courses)
	avgECTSPerYearStr := tui.RenderTotalECTSPerYear(tui.DefaultColor, courses)
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
		termWidth:         80,
		termHeight:        24,
		screen:            PickerScreen,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+q":
			if m.screen == DataScreen {
				m.screen = PickerScreen
				m.selected = make(map[int]struct{})
				return m, nil
			}

		case "up", "k":
			if m.screen == PickerScreen {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.choices) - 1
				}
			}

		case "down", "j":
			if m.screen == PickerScreen {
				m.cursor++
				if m.cursor >= len(m.choices) {
					m.cursor = 0
				}
			}

		case "enter", " ":
			if m.screen == PickerScreen {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
					m.tableStr = tui.RenderTable(tui.DefaultColor, m.headers, m.courses)
					m.avgStr = tui.RenderAverageGrades(tui.DefaultColor, m.courses)
					m.avgPerYearStr = tui.RenderAverageGradesPerYear(tui.DefaultColor, m.courses)
					m.avgECTSPerYearStr = tui.RenderTotalECTSPerYear(tui.DefaultColor, m.courses)
					m.ectsStr = tui.RenderECTS(tui.DefaultColor, m.courses)
				} else {
					m.selected = map[int]struct{}{m.cursor: {}}
					color := uniColors[m.choices[m.cursor]]
					m.tableStr = tui.RenderTable(color, m.headers, m.courses)
					m.avgStr = tui.RenderAverageGrades(color, m.courses)
					m.avgPerYearStr = tui.RenderAverageGradesPerYear(color, m.courses)
					m.avgECTSPerYearStr = tui.RenderTotalECTSPerYear(color, m.courses)
					m.ectsStr = tui.RenderECTS(color, m.courses)
					m.screen = DataScreen
				}
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
	if m.screen == PickerScreen {
		return m.renderPickerScreen()
	}
	return m.renderDataScreen()
}

func (m Model) renderPickerScreen() string {
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

	s += "\nPress Ctrl + C to quit.\n"
	return s
}

func (m Model) renderDataScreen() string {
	// Get selected university
	selectedUni := ""
	for i := range m.selected {
		selectedUni = m.choices[i]
	}

	// Get selected university color
	uniColor := tui.DefaultColor
	for i := range m.selected {
		uniColor = uniColors[m.choices[i]]
	}

	// Check if data is unavailable for this university
	if selectedUni == "TUD" || selectedUni == "TUM" {
		message := fmt.Sprintf("Data unavailable: Studies at %s have not started yet", selectedUni)
		msgBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(uniColor).
			Padding(1, 2).
			Foreground(lipgloss.Color("243")).
			Render(message)

		s := "\n" + msgBox + "\n"
		s += "\nPress Ctrl + Q to go back, Ctrl + C to quit.\n"
		return s
	}

	gap := "   "

	// Second column: average grades, average grades per year chart, total ECTS bar beneath
	col2 := lipgloss.JoinVertical(lipgloss.Left, m.avgStr, m.avgPerYearStr, "")

	// Third column: total ECTS per year chart
	col3 := lipgloss.JoinVertical(lipgloss.Left, m.avgECTSPerYearStr, "", m.ectsStr)

	// Full layout: course table | stats + avg chart + ECTS bar | ECTS/year chart
	grid := lipgloss.JoinHorizontal(lipgloss.Top, m.tableStr, gap, col2, gap, col3)

	// Create text box with width matching the grid
	gridWidth := lipgloss.Width(grid)
	textBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1).
		Width(gridWidth).
		Render("Type your notes here...")

	s := "\n" + textBox + "\n"
	s += "\n" + grid + "\n"
	s += "\nPress Ctrl + Q to go back, Ctrl + C to quit.\n"

	return s
}
