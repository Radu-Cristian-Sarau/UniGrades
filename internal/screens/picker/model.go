package picker

import (
	"fmt"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"UniGrades/internal/api"
	"UniGrades/internal/screens/grades"
	"UniGrades/internal/tui"
	"UniGrades/internal/university"
)

type Screen int

const (
	PickerScreen Screen = iota
	DataScreen
)

type Model struct {
	Choices           []string
	Cursor            int
	Selected          map[int]struct{}
	TableStr          string
	AvgStr            string
	AvgPerYearStr     string
	AvgECTSPerYearStr string
	EctsStr           string
	Headers           []string
	Courses           []bson.M
	TermWidth         int
	TermHeight        int
	Screen            Screen
	TextInput         textinput.Model
	MongoClient       *mongo.Client
	StatusMessage     string
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
		Choices:           university.Names(),
		Selected:          make(map[int]struct{}),
		TableStr:          tableStr,
		AvgStr:            avgStr,
		AvgPerYearStr:     avgPerYearStr,
		AvgECTSPerYearStr: avgECTSPerYearStr,
		EctsStr:           ectsStr,
		Headers:           headers,
		Courses:           courses,
		TermWidth:         80,
		TermHeight:        24,
		Screen:            PickerScreen,
		TextInput:         ti,
		MongoClient:       client,
		StatusMessage:     "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+q":
			if m.Screen == DataScreen {
				m.Screen = PickerScreen
				m.Selected = make(map[int]struct{})
				m.TextInput.SetValue("")
				m.StatusMessage = ""
				return m, nil
			}

		case "up", "k":
			if m.Screen == PickerScreen {
				m.Cursor--
				if m.Cursor < 0 {
					m.Cursor = len(m.Choices) - 1
				}
			}

		case "down", "j":
			if m.Screen == PickerScreen {
				m.Cursor++
				if m.Cursor >= len(m.Choices) {
					m.Cursor = 0
				}
			}

		case "enter":
			if m.Screen == PickerScreen {
				_, ok := m.Selected[m.Cursor]
				if ok {
					delete(m.Selected, m.Cursor)
					m.TableStr = tui.RenderTable(tui.DefaultColor, m.Headers, m.Courses)
					m.AvgStr = tui.RenderAverageGrades(tui.DefaultColor, m.Courses)
					m.AvgPerYearStr = tui.RenderAverageGradesPerYear(tui.DefaultColor, m.Courses)
					m.AvgECTSPerYearStr = tui.RenderTotalECTSPerYear(tui.DefaultColor, m.Courses)
					m.EctsStr = tui.RenderECTS(tui.DefaultColor, m.Courses)
				} else {
					m.Selected = map[int]struct{}{m.Cursor: {}}
					color := uniColors[m.Choices[m.Cursor]]
					m.TableStr = tui.RenderTable(color, m.Headers, m.Courses)
					m.AvgStr = tui.RenderAverageGrades(color, m.Courses)
					m.AvgPerYearStr = tui.RenderAverageGradesPerYear(color, m.Courses)
					m.AvgECTSPerYearStr = tui.RenderTotalECTSPerYear(color, m.Courses)
					m.EctsStr = tui.RenderECTS(color, m.Courses)
					m.Screen = DataScreen
				}
			} else if m.Screen == DataScreen {
				grades.HandleDataScreenInput(&m)
			}
		}

		// Handle text input when on DataScreen
		if m.Screen == DataScreen {
			var cmd tea.Cmd
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// SelectedUniversity returns the name of the selected university, or "" if none.
func (m Model) SelectedUniversity() string {
	for i := range m.Selected {
		return m.Choices[i]
	}
	return ""
}

// Interface methods for grades package - needed to avoid circular imports

// GetMongoClient returns the MongoDB client.
func (m Model) GetMongoClient() *mongo.Client {
	return m.MongoClient
}

// GetSelectedUniversity returns the selected university.
func (m Model) GetSelectedUniversity() string {
	return m.SelectedUniversity()
}

// GetTermWidth returns the terminal width.
func (m Model) GetTermWidth() int {
	return m.TermWidth
}

// GetTableStr returns the rendered table string.
func (m Model) GetTableStr() string {
	return m.TableStr
}

// GetAvgStr returns the rendered average grades string.
func (m Model) GetAvgStr() string {
	return m.AvgStr
}

// GetAvgPerYearStr returns the rendered average grades per year string.
func (m Model) GetAvgPerYearStr() string {
	return m.AvgPerYearStr
}

// GetAvgECTSPerYearStr returns the rendered average ECTS per year string.
func (m Model) GetAvgECTSPerYearStr() string {
	return m.AvgECTSPerYearStr
}

// GetEctsStr returns the rendered ECTS string.
func (m Model) GetEctsStr() string {
	return m.EctsStr
}

// GetTextInputView returns the text input view.
func (m Model) GetTextInputView() string {
	return m.TextInput.View()
}

// GetStatusMessage returns the status message.
func (m Model) GetStatusMessage() string {
	return m.StatusMessage
}

// SetStatusMessage sets the status message.
func (m *Model) SetStatusMessage(msg string) {
	m.StatusMessage = msg
}

// RefreshCourses refreshes the courses from the database.
func (m *Model) RefreshCourses() {
	m.Courses = api.GetAllCourses(m.MongoClient)
}

// RefreshTableStr refreshes the table string with the given color.
func (m *Model) RefreshTableStr(color lipgloss.Color) {
	m.TableStr = tui.RenderTable(color, m.Headers, m.Courses)
}

// RefreshAvgStr refreshes the average grades string with the given color.
func (m *Model) RefreshAvgStr(color lipgloss.Color) {
	m.AvgStr = tui.RenderAverageGrades(color, m.Courses)
}

// RefreshAvgPerYearStr refreshes the average grades per year string with the given color.
func (m *Model) RefreshAvgPerYearStr(color lipgloss.Color) {
	m.AvgPerYearStr = tui.RenderAverageGradesPerYear(color, m.Courses)
}

// RefreshAvgECTSPerYearStr refreshes the average ECTS per year string with the given color.
func (m *Model) RefreshAvgECTSPerYearStr(color lipgloss.Color) {
	m.AvgECTSPerYearStr = tui.RenderTotalECTSPerYear(color, m.Courses)
}

// RefreshEctsStr refreshes the ECTS string with the given color.
func (m *Model) RefreshEctsStr(color lipgloss.Color) {
	m.EctsStr = tui.RenderECTS(color, m.Courses)
}

// GetTextInputValue returns the text input value.
func (m Model) GetTextInputValue() string {
	return m.TextInput.Value()
}

// SetTextInputValue sets the text input value.
func (m *Model) SetTextInputValue(value string) {
	m.TextInput.SetValue(value)
}

var uniColors = university.ColorMap()

func (m Model) View() string {
	if m.Screen == PickerScreen {
		return m.renderPickerScreen()
	}
	return grades.RenderDataScreen(&m)
}

func (m Model) renderPickerScreen() string {
	s := "\n\n\nSelect university:\n\n"

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		style := lipgloss.NewStyle().Foreground(uniColors[choice])
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, style.Render(choice))
	}

	s += "\nPress Ctrl + C to quit."

	// Center the entire picker screen
	centeredContent := lipgloss.NewStyle().
		Width(m.TermWidth).
		Align(lipgloss.Center).
		Render(s)

	return "\n" + centeredContent + "\n"
}
