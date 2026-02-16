// Package grades provides the data/grades screen for displaying course information and statistics.
package grades

import (
	// Standard library imports
	"fmt"     // Formatted I/O and string conversion
	"strconv" // String to number conversions
	"strings" // String manipulation

	// Terminal UI libraries
	"github.com/charmbracelet/lipgloss"       // Styling and layout
	"github.com/charmbracelet/lipgloss/table" // Table rendering
	"go.mongodb.org/mongo-driver/v2/mongo"    // MongoDB client

	// Internal packages
	"UniGrades/internal/api"        // Database operations
	"UniGrades/internal/tui"        // UI rendering
	"UniGrades/internal/university" // University data
)

// DataScreenModel defines the interface for the data screen model.
// This interface allows the data screen to interact with the picker model
// without creating circular dependencies between packages.
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

// HandleDataScreenInput processes user text input on the data screen.
// It parses commands prefixed with '/' and delegates to appropriate handlers.
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

// ProcessAddCommand parses and executes the /add command.
// Format: /add Name Year Grade ECTS
// Example: /add Applied_Math 1 7 5
func ProcessAddCommand(m DataScreenModel, input string) {
	// Parse command arguments
	parts := strings.Fields(input)
	if len(parts) < 5 {
		m.SetStatusMessage("Invalid format. Use: /add Name Year Grade ECTS")
		return
	}

	name := parts[1]
	year, errYear := strconv.Atoi(parts[2])
	grade, errGrade := strconv.ParseFloat(parts[3], 64)
	ects, errEcts := strconv.Atoi(parts[4])

	// Validate all numeric conversions
	if errYear != nil || errGrade != nil || errEcts != nil {
		m.SetStatusMessage("Error: Year and ECTS must be integers, Grade must be a number")
		return
	}

	// Create course struct and insert into database
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

	// Refresh all displays to show the new course
	m.RefreshCourses()
	RefreshCharts(m)

	m.SetStatusMessage(fmt.Sprintf("✓ Course '%s' added successfully (ID: %s)", name, id))
}

// ProcessDeleteCommand parses and executes the /delete command.
// Format: /delete CourseName
// Example: /delete Applied_Math
func ProcessDeleteCommand(m DataScreenModel, input string) {
	// Parse command arguments
	parts := strings.Fields(input)
	if len(parts) < 2 {
		m.SetStatusMessage("Invalid format. Use: /delete CourseName")
		return
	}

	courseName := parts[1]

	// Delete course from database
	err := api.DeleteCourse(m.GetMongoClient(), courseName)
	if err != nil {
		m.SetStatusMessage(fmt.Sprintf("Error deleting course: %v", err))
		return
	}

	// Refresh all displays
	m.RefreshCourses()
	RefreshCharts(m)

	m.SetStatusMessage(fmt.Sprintf("✓ Course '%s' deleted successfully", courseName))
}

// ProcessEditCommand parses and executes the /edit command.
// Format: /edit CourseName Field NewValue
// Valid fields: Name, Year, Grade, ECTS
// Example: /edit Applied_Math Grade 9
func ProcessEditCommand(m DataScreenModel, input string) {
	// Parse command arguments
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

	// Update course in database
	err := api.UpdateCourse(m.GetMongoClient(), courseName, field, newValue)
	if err != nil {
		m.SetStatusMessage(fmt.Sprintf("Error updating course: %v", err))
		return
	}

	// Refresh all displays
	m.RefreshCourses()
	RefreshCharts(m)

	m.SetStatusMessage(fmt.Sprintf("✓ Course '%s' field '%s' updated to '%v'", courseName, field, newValue))
}

// RefreshCharts updates all chart and statistics displays.
// Recomputes tables and visualizations based on current university and course data.
func RefreshCharts(m DataScreenModel) {
	selectedUni := m.GetSelectedUniversity()
	color := tui.DefaultColor
	if selectedUni != "" {
		color = university.ColorMap()[selectedUni]
	}

	// Only refresh charts for universities with data
	if selectedUni != "TUD" && selectedUni != "TUM" {
		m.RefreshTableStr(color)
		m.RefreshAvgStr(color)
		m.RefreshAvgPerYearStr(color)
		m.RefreshAvgECTSPerYearStr(color)
		m.RefreshEctsStr(color)
	}
}

// RenderDataScreen renders the complete data/grades screen display.
// Shows courses table, statistics, charts, and command help.
func RenderDataScreen(m DataScreenModel) string {
	// Get selected university and its color
	selectedUni := m.GetSelectedUniversity()
	uniColor := tui.DefaultColor
	if selectedUni != "" {
		uniColor = university.ColorMap()[selectedUni]
	}

	// Check if data is unavailable for this university (e.g., future studies)
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

	// Get all rendered visualization strings from model
	tableStr := m.GetTableStr()
	avgStr := m.GetAvgStr()
	avgPerYearStr := m.GetAvgPerYearStr()
	avgECTSPerYearStr := m.GetAvgECTSPerYearStr()
	ectsStr := m.GetEctsStr()

	// Organize columns: average stats + per-year average chart
	col2 := lipgloss.JoinVertical(lipgloss.Left, avgStr, avgPerYearStr, "")

	// Third column: per-year ECTS chart + total ECTS bar
	col3 := lipgloss.JoinVertical(lipgloss.Left, avgECTSPerYearStr, "", ectsStr)

	// Fourth column: help sections with command reference and error explanations
	helpCommands := RenderCommandsHelp(uniColor)
	helpErrors := RenderErrorsExplanation(uniColor)
	helpSection := lipgloss.JoinVertical(lipgloss.Left, helpCommands, "", helpErrors)

	// Main layout: arrange all columns horizontally
	grid := lipgloss.JoinHorizontal(lipgloss.Top, tableStr, gap, col2, gap, col3, gap, helpSection)

	// Create text input box with appropriate styling
	gridWidth := lipgloss.Width(grid)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(uniColor).
		Padding(0, 1).
		Width(gridWidth - 2)

	textInputBox := inputStyle.Render(m.GetTextInputView())

	centeredInput := lipgloss.NewStyle().
		Width(m.GetTermWidth()).
		Align(lipgloss.Center).
		Render(textInputBox)

	s := "\n" + centeredInput + "\n"

	// Display status message if available
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

	// Center the main content grid
	centeredGrid := lipgloss.NewStyle().
		Width(m.GetTermWidth()).
		Align(lipgloss.Center).
		Render(grid)

	s += centeredGrid + "\n"

	// Add footer with keyboard shortcuts
	footerStyle := lipgloss.NewStyle().
		Width(m.GetTermWidth()).
		Align(lipgloss.Center)
	s += "\n" + footerStyle.Render("Press Ctrl + Q to go back, Ctrl + C to quit.") + "\n"

	return s
}

// RenderCommandsHelp renders a table showing available commands and their usage.
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

// RenderErrorsExplanation renders a table explaining common error messages.
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
