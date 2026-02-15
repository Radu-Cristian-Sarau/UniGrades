package picker

import (
	"fmt"
	"strconv"
	"strings"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"UniGrades/internal/api"
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
	textInput         textinput.Model
	mongoClient       *mongo.Client
	statusMessage     string
}

func InitialModel(headers []string, courses []bson.M, client *mongo.Client) Model {
	ti := textinput.New()
	ti.Placeholder = "Type /add Name Year Grade ECTS or your notes..."
	ti.Focus()

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
		textInput:         ti,
		mongoClient:       client,
		statusMessage:     "",
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
				m.textInput.SetValue("")
				m.statusMessage = ""
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

		case "enter":
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
			} else if m.screen == DataScreen {
				// Handle enter in text input - process /add command
				input := m.textInput.Value()
				if strings.HasPrefix(input, "/add ") {
					m.processAddCommand(input)
					m.textInput.SetValue("")
				}
			}
		}

		// Handle text input when on DataScreen
		if m.screen == DataScreen {
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// processAddCommand parses and executes the /add command
func (m *Model) processAddCommand(input string) {
	// Parse: /add Name Year Grade ECTS
	parts := strings.Fields(input)
	if len(parts) < 5 {
		m.statusMessage = "Invalid format. Use: /add Name Year Grade ECTS"
		return
	}

	name := parts[1]
	year, errYear := strconv.Atoi(parts[2])
	grade, errGrade := strconv.ParseFloat(parts[3], 64)
	ects, errEcts := strconv.Atoi(parts[4])

	if errYear != nil || errGrade != nil || errEcts != nil {
		m.statusMessage = "Error: Year and ECTS must be integers, Grade must be a number"
		return
	}

	// Create course and add to database
	course := api.Course{
		Name:  name,
		Year:  year,
		Grade: grade,
		ECTS:  ects,
	}

	id, err := api.AddCourse(m.mongoClient, course)
	if err != nil {
		m.statusMessage = fmt.Sprintf("Error adding course: %v", err)
		return
	}

	// Refresh course list
	m.courses = api.GetAllCourses(m.mongoClient)
	m.refreshCharts()

	m.statusMessage = fmt.Sprintf("✓ Course '%s' added successfully (ID: %s)", name, id)
}

// refreshCharts updates all the chart displays
func (m *Model) refreshCharts() {
	selectedUni := ""
	for i := range m.selected {
		selectedUni = m.choices[i]
	}

	color := tui.DefaultColor
	for i := range m.selected {
		color = uniColors[m.choices[i]]
	}

	if selectedUni != "TUD" && selectedUni != "TUM" {
		m.tableStr = tui.RenderTable(color, m.headers, m.courses)
		m.avgStr = tui.RenderAverageGrades(color, m.courses)
		m.avgPerYearStr = tui.RenderAverageGradesPerYear(color, m.courses)
		m.avgECTSPerYearStr = tui.RenderTotalECTSPerYear(color, m.courses)
		m.ectsStr = tui.RenderECTS(color, m.courses)
	}
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
		message := fmt.Sprintf("Data unavailable: Studies at %s have not started yet.", selectedUni)
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

	// Create text input box with width matching the grid
	gridWidth := lipgloss.Width(grid)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1).
		Width(gridWidth - 2) // Account for padding

	textInputBox := inputStyle.Render(m.textInput.View())

	s := "\n" + textInputBox + "\n"

	// Show status message if available
	if m.statusMessage != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)
		s += statusStyle.Render(m.statusMessage) + "\n\n"
	} else {
		s += "\n"
	}

	s += grid + "\n"
	s += "\nPress Ctrl + Q to go back, Ctrl + C to quit.\n"

	return s
}
