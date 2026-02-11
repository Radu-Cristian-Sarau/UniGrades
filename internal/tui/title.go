package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

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

func RenderTitle() string {
	width, _, _ := term.GetSize(uintptr(os.Stdout.Fd()))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#17b118")).Bold(true).Width(width).Align(lipgloss.Center)
	return titleStyle.Render(title)
}
