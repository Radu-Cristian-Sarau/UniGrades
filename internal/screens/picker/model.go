package picker

import (
	"fmt"
	"strconv"
	"strings"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
	ti.Placeholder = "Commands: /add Name Year Grade ECTS | /edit Name Field Value | /delete Name"
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
				// Handle enter in text input - process commands
				input := m.textInput.Value()
				if strings.HasPrefix(input, "/add ") {
					m.processAddCommand(input)
					m.textInput.SetValue("")
				} else if strings.HasPrefix(input, "/delete ") {
					m.processDeleteCommand(input)
					m.textInput.SetValue("")
				} else if strings.HasPrefix(input, "/edit ") {
					m.processEditCommand(input)
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

// processDeleteCommand parses and executes the /delete command
func (m *Model) processDeleteCommand(input string) {
	// Parse: /delete CourseName
	parts := strings.Fields(input)
	if len(parts) < 2 {
		m.statusMessage = "Invalid format. Use: /delete CourseName"
		return
	}

	courseName := parts[1]

	err := api.DeleteCourse(m.mongoClient, courseName)
	if err != nil {
		m.statusMessage = fmt.Sprintf("Error deleting course: %v", err)
		return
	}

	// Refresh course list
	m.courses = api.GetAllCourses(m.mongoClient)
	m.refreshCharts()

	m.statusMessage = fmt.Sprintf("✓ Course '%s' deleted successfully", courseName)
}

// processEditCommand parses and executes the /edit command
// Format: /edit CourseName Field NewValue
// Example: /edit Calculus Grade 9.5
func (m *Model) processEditCommand(input string) {
	// Parse: /edit CourseName Field NewValue
	parts := strings.Fields(input)
	if len(parts) < 4 {
		m.statusMessage = "Invalid format. Use: /edit CourseName Field NewValue (e.g., /edit Calculus Grade 9.5)"
		return
	}

	courseName := parts[1]
	field := parts[2]
	newValue := parts[3]

	// Validate field name
	validFields := map[string]bool{"Name": true, "Year": true, "Grade": true, "ECTS": true}
	if !validFields[field] {
		m.statusMessage = "Invalid field. Valid fields are: Name, Year, Grade, ECTS"
		return
	}

	err := api.UpdateCourse(m.mongoClient, courseName, field, newValue)
	if err != nil {
		m.statusMessage = fmt.Sprintf("Error updating course: %v", err)
		return
	}

	// Refresh course list
	m.courses = api.GetAllCourses(m.mongoClient)
	m.refreshCharts()

	m.statusMessage = fmt.Sprintf("✓ Course '%s' field '%s' updated to '%v'", courseName, field, newValue)
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

	// Fourth column: help sections (commands table and errors table stacked)
	helpCommands := m.renderCommandsHelp(uniColor)
	helpErrors := m.renderErrorsExplanation(uniColor)
	helpSection := lipgloss.JoinVertical(lipgloss.Left, helpCommands, "", helpErrors)

	// Full layout: course table | stats + avg chart + ECTS bar | ECTS/year chart | help
	grid := lipgloss.JoinHorizontal(lipgloss.Top, m.tableStr, gap, col2, gap, col3, gap, helpSection)

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

func (m Model) renderCommandsHelp(uniColor lipgloss.Color) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Foreground(tui.DefaultColor).Align(lipgloss.Center)
			case row%2 == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Padding(0, 1)
			default:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(0, 1)
			}
		}).
		Headers("Command", "Description", "Example").
		Rows(
			[]string{"/add", "Add new course", "/add Applied_Math 1 7 5"},
			[]string{"/edit", "Update course field", "/edit Applied_Math Grade 9"},
			[]string{"/delete", "Delete course", "/delete Applied_math"},
		)

	return t.Render()
}

func (m Model) renderErrorsExplanation(uniColor lipgloss.Color) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Foreground(tui.DefaultColor).Align(lipgloss.Center)
			case row%2 == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Padding(0, 1)
			default:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(0, 1)
			}
		}).
		Headers("Error", "Explanation").
		Rows(
			[]string{"Invalid format", "Wrong command syntax"},
			[]string{"Course not found", "Course name doesn't exist"},
			[]string{"Year not integer", "Year must be a number"},
			[]string{"Grade not number", "Grade must be decimal/int"},
			[]string{"ECTS not integer", "ECTS must be a number"},
			[]string{"Invalid field", "Field not in Name/Year/Grade/ECTS"},
		)

	return t.Render()
}
