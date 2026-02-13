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

const TitleColor = lipgloss.Color("#17b118")

func RenderTitle() string {
	width, _, _ := term.GetSize(uintptr(os.Stdout.Fd()))
	titleStyle := lipgloss.NewStyle().Foreground(TitleColor).Bold(true).Width(width).Align(lipgloss.Center)
	return titleStyle.Render(title)
}
