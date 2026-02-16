// Package picker provides the university selection screen for the UniGrades application.
// This screen allows users to select a university and view their course data.
package picker

import (
	// Standard library imports
	"fmt" // Formatted I/O

	// Bubble Tea components
	textinput "github.com/charmbracelet/bubbles/textinput" // Text input component
	tea "github.com/charmbracelet/bubbletea"               // TUI framework
	"github.com/charmbracelet/lipgloss"                    // Styling and layout

	// MongoDB types
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	// Internal packages
	"UniGrades/internal/api"            // Database operations
	"UniGrades/internal/screens/grades" // Data screen
	"UniGrades/internal/tui"            // UI rendering
	"UniGrades/internal/university"     // University data
)

// Screen represents the current screen being displayed.
type Screen int

const (
	// PickerScreen shows the university selection menu
	PickerScreen Screen = iota
	// DataScreen shows the course data and statistics
	DataScreen
)

// Model represents the application state for the TUI.
type Model struct {
	// Universe selection state
	Choices  []string         // Available university names
	Cursor   int              // Currently selected university index
	Selected map[int]struct{} // Map of selected university indices

	// Rendered display strings (cached for performance)
	TableStr          string // Rendered course table
	AvgStr            string // Rendered average grades table
	AvgPerYearStr     string // Rendered grades per year chart
	AvgECTSPerYearStr string // Rendered ECTS per year chart
	EctsStr           string // Rendered total ECTS bar

	// Course data
	Headers []string // Column headers for the table
	Courses []bson.M // Course documents from database

	// Terminal state
	TermWidth  int // Width of the terminal
	TermHeight int // Height of the terminal

	// Screen management
	Screen    Screen          // Current screen (Picker or Data)
	TextInput textinput.Model // Text input component for commands

	// External resources
	MongoClient   *mongo.Client // MongoDB connection
	StatusMessage string        // User feedback message
}

// InitialModel creates and returns a new Model with initial state.
func InitialModel(headers []string, courses []bson.M, client *mongo.Client) Model {
	// Initialize text input for commands
	ti := textinput.New()
	ti.Placeholder = "Commands: /add Name Year Grade ECTS | /edit Name Field Value | /delete Name"
	ti.Focus()

	// Render all initial displays
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

// Init initializes the model. Called by Bubble Tea on startup.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates model state accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update terminal dimensions when window is resized
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// Quit application
			return m, tea.Quit

		case "ctrl+q":
			// Return to picker screen from data screen
			if m.Screen == DataScreen {
				m.Screen = PickerScreen
				m.Selected = make(map[int]struct{})
				m.TextInput.SetValue("")
				m.StatusMessage = ""
				return m, nil
			}

		case "up", "k":
			// Navigate up in picker screen
			if m.Screen == PickerScreen {
				m.Cursor--
				if m.Cursor < 0 {
					m.Cursor = len(m.Choices) - 1
				}
			}

		case "down", "j":
			// Navigate down in picker screen
			if m.Screen == PickerScreen {
				m.Cursor++
				if m.Cursor >= len(m.Choices) {
					m.Cursor = 0
				}
			}

		case "enter":
			// Toggle selection in picker screen, or handle commands in data screen
			if m.Screen == PickerScreen {
				_, ok := m.Selected[m.Cursor]
				if ok {
					// Deselect the university
					delete(m.Selected, m.Cursor)
					m.TableStr = tui.RenderTable(tui.DefaultColor, m.Headers, m.Courses)
					m.AvgStr = tui.RenderAverageGrades(tui.DefaultColor, m.Courses)
					m.AvgPerYearStr = tui.RenderAverageGradesPerYear(tui.DefaultColor, m.Courses)
					m.AvgECTSPerYearStr = tui.RenderTotalECTSPerYear(tui.DefaultColor, m.Courses)
					m.EctsStr = tui.RenderECTS(tui.DefaultColor, m.Courses)
				} else {
					// Select the university
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
				// Handle text input commands on data screen
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

// SelectedUniversity returns the name of the selected university, or empty string if none.
func (m Model) SelectedUniversity() string {
	for i := range m.Selected {
		return m.Choices[i]
	}
	return ""
}

// Interface implementation methods for grades.DataScreenModel
// These methods allow the data screen to access model state without circular imports.

// GetMongoClient returns the MongoDB client.
func (m Model) GetMongoClient() *mongo.Client {
	return m.MongoClient
}

// GetSelectedUniversity returns the selected university name.
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

// GetStatusMessage returns the current status message.
func (m Model) GetStatusMessage() string {
	return m.StatusMessage
}

// SetStatusMessage sets the status message to display to the user.
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

// GetTextInputValue returns the current text input value.
func (m Model) GetTextInputValue() string {
	return m.TextInput.Value()
}

// SetTextInputValue sets the text input value.
func (m *Model) SetTextInputValue(value string) {
	m.TextInput.SetValue(value)
}

// Module-level variable holding the university color map.
var uniColors = university.ColorMap()

// View renders the current screen based on the model state.
func (m Model) View() string {
	if m.Screen == PickerScreen {
		return m.renderPickerScreen()
	}
	return grades.RenderDataScreen(&m)
}

// renderPickerScreen renders the university selection menu.
func (m Model) renderPickerScreen() string {
	s := "\n\n\nSelect university:\n\n"

	// Render each university option with selection cursor and checkbox
	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		// Apply university color styling to the choice
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
