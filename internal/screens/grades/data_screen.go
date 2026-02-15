package grades

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"UniGrades/internal/api"
	"UniGrades/internal/tui"
	"UniGrades/internal/university"
)

// DataScreenModel defines the interface for the data screen model
type DataScreenModel interface {
	GetMongoClient() *mongo.Client
	GetSelectedUniversity() string
	GetTermWidth() int
	GetTableStr() string
	GetAvgStr() string
	GetAvgPerYearStr() string
	GetAvgECTSPerYearStr() string
	GetEctsStr() string
	GetTextInputView() string
	GetStatusMessage() string
	SetStatusMessage(msg string)
	RefreshCourses()
	RefreshTableStr(lipgloss.Color)
	RefreshAvgStr(lipgloss.Color)
	RefreshAvgPerYearStr(lipgloss.Color)
	RefreshAvgECTSPerYearStr(lipgloss.Color)
	RefreshEctsStr(lipgloss.Color)
	GetTextInputValue() string
	SetTextInputValue(string)
}

// HandleDataScreenInput processes text input on the DataScreen
func HandleDataScreenInput(m DataScreenModel) {
	input := m.GetTextInputValue()
	if strings.HasPrefix(input, "/add ") {
		ProcessAddCommand(m, input)
		m.SetTextInputValue("")
	} else if strings.HasPrefix(input, "/delete ") {
		ProcessDeleteCommand(m, input)
		m.SetTextInputValue("")
	} else if strings.HasPrefix(input, "/edit ") {
		ProcessEditCommand(m, input)
		m.SetTextInputValue("")
	}
}

// ProcessAddCommand parses and executes the /add command
func ProcessAddCommand(m DataScreenModel, input string) {
	// Parse: /add Name Year Grade ECTS
	parts := strings.Fields(input)
	if len(parts) < 5 {
		m.SetStatusMessage("Invalid format. Use: /add Name Year Grade ECTS")
		return
	}

	name := parts[1]
	year, errYear := strconv.Atoi(parts[2])
	grade, errGrade := strconv.ParseFloat(parts[3], 64)
	ects, errEcts := strconv.Atoi(parts[4])

	if errYear != nil || errGrade != nil || errEcts != nil {
		m.SetStatusMessage("Error: Year and ECTS must be integers, Grade must be a number")
		return
	}

	// Create course and add to database
	course := api.Course{
		Name:  name,
		Year:  year,
		Grade: grade,
		ECTS:  ects,
	}

	id, err := api.AddCourse(m.GetMongoClient(), course)
	if err != nil {
		m.SetStatusMessage(fmt.Sprintf("Error adding course: %v", err))
		return
	}

	// Refresh course list
	m.RefreshCourses()
	RefreshCharts(m)

	m.SetStatusMessage(fmt.Sprintf("✓ Course '%s' added successfully (ID: %s)", name, id))
}

// ProcessDeleteCommand parses and executes the /delete command
func ProcessDeleteCommand(m DataScreenModel, input string) {
	// Parse: /delete CourseName
	parts := strings.Fields(input)
	if len(parts) < 2 {
		m.SetStatusMessage("Invalid format. Use: /delete CourseName")
		return
	}

	courseName := parts[1]

	err := api.DeleteCourse(m.GetMongoClient(), courseName)
	if err != nil {
		m.SetStatusMessage(fmt.Sprintf("Error deleting course: %v", err))
		return
	}

	// Refresh course list
	m.RefreshCourses()
	RefreshCharts(m)

	m.SetStatusMessage(fmt.Sprintf("✓ Course '%s' deleted successfully", courseName))
}

// ProcessEditCommand parses and executes the /edit command
// Format: /edit CourseName Field NewValue
// Example: /edit Calculus Grade 9.5
func ProcessEditCommand(m DataScreenModel, input string) {
	// Parse: /edit CourseName Field NewValue
	parts := strings.Fields(input)
	if len(parts) < 4 {
		m.SetStatusMessage("Invalid format. Use: /edit CourseName Field NewValue (e.g., /edit Applied_Math Grade 9)")
		return
	}

	courseName := parts[1]
	field := parts[2]
	newValue := parts[3]

	// Validate field name
	validFields := map[string]bool{"Name": true, "Year": true, "Grade": true, "ECTS": true}
	if !validFields[field] {
		m.SetStatusMessage("Invalid field. Valid fields are: Name, Year, Grade, ECTS")
		return
	}

	err := api.UpdateCourse(m.GetMongoClient(), courseName, field, newValue)
	if err != nil {
		m.SetStatusMessage(fmt.Sprintf("Error updating course: %v", err))
		return
	}

	// Refresh course list
	m.RefreshCourses()
	RefreshCharts(m)

	m.SetStatusMessage(fmt.Sprintf("✓ Course '%s' field '%s' updated to '%v'", courseName, field, newValue))
}

// RefreshCharts updates all the chart displays
func RefreshCharts(m DataScreenModel) {
	selectedUni := m.GetSelectedUniversity()
	color := tui.DefaultColor
	if selectedUni != "" {
		color = university.ColorMap()[selectedUni]
	}

	if selectedUni != "TUD" && selectedUni != "TUM" {
		m.RefreshTableStr(color)
		m.RefreshAvgStr(color)
		m.RefreshAvgPerYearStr(color)
		m.RefreshAvgECTSPerYearStr(color)
		m.RefreshEctsStr(color)
	}
}

// RenderDataScreen renders the data/grades screen
func RenderDataScreen(m DataScreenModel) string {
	// Get selected university
	selectedUni := m.GetSelectedUniversity()

	// Get selected university color
	uniColor := tui.DefaultColor
	if selectedUni != "" {
		uniColor = university.ColorMap()[selectedUni]
	}

	// Check if data is unavailable for this university
	if selectedUni == "TUD" || selectedUni == "TUM" {
		message := fmt.Sprintf("Data unavailable: Studies at %s have not started yet.", selectedUni)
		msgBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(uniColor).
			Padding(1, 2).
			Foreground(ColorDimText).
			Render(message)

		contentStyle := lipgloss.NewStyle().
			Width(m.GetTermWidth()).
			Align(lipgloss.Center)

		s := "\n" + contentStyle.Render(msgBox) + "\n"

		footerStyle := lipgloss.NewStyle().
			Width(m.GetTermWidth()).
			Align(lipgloss.Center)
		s += "\n" + footerStyle.Render("Press Ctrl + Q to go back, Ctrl + C to quit.") + "\n"

		return s
	}

	gap := "   "

	// Get rendered strings from model
	tableStr := m.GetTableStr()
	avgStr := m.GetAvgStr()
	avgPerYearStr := m.GetAvgPerYearStr()
	avgECTSPerYearStr := m.GetAvgECTSPerYearStr()
	ectsStr := m.GetEctsStr()

	// Second column: average grades, average grades per year chart, total ECTS bar beneath
	col2 := lipgloss.JoinVertical(lipgloss.Left, avgStr, avgPerYearStr, "")

	// Third column: total ECTS per year chart
	col3 := lipgloss.JoinVertical(lipgloss.Left, avgECTSPerYearStr, "", ectsStr)

	// Fourth column: help sections (commands table and errors table stacked)
	helpCommands := RenderCommandsHelp(uniColor)
	helpErrors := RenderErrorsExplanation(uniColor)
	helpSection := lipgloss.JoinVertical(lipgloss.Left, helpCommands, "", helpErrors)

	// Full layout: course table | stats + avg chart + ECTS bar | ECTS/year chart | help
	grid := lipgloss.JoinHorizontal(lipgloss.Top, tableStr, gap, col2, gap, col3, gap, helpSection)

	// Create text input box with width matching the grid
	gridWidth := lipgloss.Width(grid)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1).
		Width(gridWidth - 2) // Account for padding

	textInputBox := inputStyle.Render(m.GetTextInputView())

	// Center everything horizontally
	centeredInput := lipgloss.NewStyle().
		Width(m.GetTermWidth()).
		Align(lipgloss.Center).
		Render(textInputBox)

	s := "\n" + centeredInput + "\n"

	// Show status message if available
	statusMsg := m.GetStatusMessage()
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			Width(m.GetTermWidth()).
			Align(lipgloss.Center)
		s += statusStyle.Render(statusMsg) + "\n\n"
	} else {
		s += "\n"
	}

	// Center the grid horizontally
	centeredGrid := lipgloss.NewStyle().
		Width(m.GetTermWidth()).
		Align(lipgloss.Center).
		Render(grid)

	s += centeredGrid + "\n"

	footerStyle := lipgloss.NewStyle().
		Width(m.GetTermWidth()).
		Align(lipgloss.Center)
	s += "\n" + footerStyle.Render("Press Ctrl + Q to go back, Ctrl + C to quit.") + "\n"

	return s
}

// RenderCommandsHelp renders the commands help table
func RenderCommandsHelp(uniColor lipgloss.Color) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Foreground(tui.DefaultColor).Align(lipgloss.Center)
			case row%2 == 0:
				return lipgloss.NewStyle().Foreground(ColorTableAlternate1).Padding(0, 1)
			default:
				return lipgloss.NewStyle().Foreground(ColorTableAlternate2).Padding(0, 1)
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

// RenderErrorsExplanation renders the errors explanation table
func RenderErrorsExplanation(uniColor lipgloss.Color) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Foreground(tui.DefaultColor).Align(lipgloss.Center)
			case row%2 == 0:
				return lipgloss.NewStyle().Foreground(ColorTableAlternate1).Padding(0, 1)
			default:
				return lipgloss.NewStyle().Foreground(ColorTableAlternate2).Padding(0, 1)
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
