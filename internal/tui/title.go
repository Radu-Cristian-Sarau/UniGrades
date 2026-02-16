// Package tui provides terminal user interface rendering components for UniGrades.
package tui

import (
	// Standard library imports
	"os" // Operating system functionality

	// TUI libraries
	"github.com/charmbracelet/lipgloss" // Styling and layout
	"github.com/charmbracelet/x/term"   // Terminal utilities
)

// ASCII art banner for the application title
const title = `
███    █▄  ███▄▄▄▄    ▄█     ▄██████▄     ▄████████    ▄████████ ████████▄     ▄████████    ▄████████ 
███    ███ ███▀▀▀██▄ ███    ███    ███   ███    ███   ███    ███ ███   ▀███   ███    ███   ███    ███ 
███    ███ ███   ███ ███▌   ███    █▀    ███    ███   ███    ███ ███    ███   ███    █▀    ███    █▀  
███    ███ ███   ███ ███▌  ▄███         ▄███▄▄▄▄██▀   ███    ███ ███    ███  ▄███▄▄▄       ███        
███    ███ ███   ███ ███▌ ▀▀███ ████▄  ▀▀███▀▀▀▀▀   ▀███████████ ███    ███ ▀▀███▀▀▀     ▀███████████ 
███    ███ ███   ███ ███    ███    ███ ▀███████████   ███    ███ ███    ███   ███    █▄           ███ 
███    ███ ███   ███ ███    ███    ███   ███    ███   ███    ███ ███   ▄███   ███    ███    ▄█    ███ 
████████▀   ▀█   █▀  █▀     ████████▀    ███    ███   ███    █▀  ████████▀    ██████████  ▄████████▀  
`

// TitleColor is the green color used for the title banner
const TitleColor = lipgloss.Color("#17b118")

// LinkColor is the grey color used for the repository link
const LinkColor = lipgloss.Color("245")

// repoURL is the GitHub repository URL for UniGrades
const repoURL = "https://github.com/Radu-Cristian-Sarau/UniGrades"

// RenderTitle returns a formatted string containing the application title,
// banner, and a clickable link to the repository.
func RenderTitle() string {
	// Get terminal width to center the title
	width, _, _ := term.GetSize(uintptr(os.Stdout.Fd()))
	titleStyle := lipgloss.NewStyle().Foreground(TitleColor).Bold(true).Width(width).Align(lipgloss.Center)

	// Create OSC 8 hyperlink (terminal-compatible hyperlink format)
	// Format: \033]8;;URL\033\\LABEL\033]8;;\033\\
	link := "\033]8;;" + repoURL + "\033\\" + repoURL + "\033]8;;\033\\"
	linkStyle := lipgloss.NewStyle().Foreground(LinkColor).Width(width).Align(lipgloss.Center)

	return titleStyle.Render(title) + "\n" + linkStyle.Render(link)
}
