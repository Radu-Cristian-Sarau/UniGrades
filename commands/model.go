// Package main demonstrates a simple Bubble Tea example that checks
// the HTTP status of a remote server. This is example code based on Bubble Tea tutorials.
package main

import (
	// Standard library imports
	"fmt"      // Formatted I/O
	"net/http" // HTTP client functionality
	"os"       // Operating system operations
	"time"     // Time utilities

	// Bubble Tea framework for terminal UI
	tea "github.com/charmbracelet/bubbletea"
)

// URL to check for HTTP status
const url = "https://charm.sh/"

// model represents the application state.
type model struct {
	// status stores the HTTP response status code
	status int
	// err stores any error that occurred during the request
	err error
}

// checkServer makes an HTTP GET request to the configured URL and returns
// either the status code or an error wrapped in a tea.Msg.
func checkServer() tea.Msg {
	// Create an HTTP client with a 10-second timeout
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	// If the request failed, return the error wrapped in an errMsg
	if err != nil {
		return errMsg{err}
	}

	// Return the HTTP status code as a statusMsg
	return statusMsg(res.StatusCode)
}

// statusMsg is a message type that contains an HTTP status code.
type statusMsg int

// errMsg is a message type that contains an error.
type errMsg struct{ err error }

// Error implements the error interface for errMsg, allowing it to be used as an error.
func (e errMsg) Error() string { return e.err.Error() }

// Init initializes the model by starting the HTTP check.
func (m model) Init() tea.Cmd {
	return checkServer
}

// Update handles incoming messages and updates the model state accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		// Store the HTTP status and quit the program
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		// Store the error and quit the program
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		// Handle Ctrl+C to quit gracefully
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// Ignore all other messages
	return m, nil
}

// View renders the current state as a string for display in the terminal.
func (m model) View() string {
	// If an error occurred, display it and exit
	if m.err != nil {
		return fmt.Sprintf("\n We had some trouble: %v\n\n", m.err)
	}

	// Build the status message
	s := fmt.Sprintf("Checking %s ... ", url)

	// Add the status code and text if available
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n\n"
}

// main entry point for the HTTP status checker example.
func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
