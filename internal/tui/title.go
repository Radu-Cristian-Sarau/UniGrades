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

const repoURL = "https://github.com/Radu-Cristian-Sarau/UniGrades"

func RenderTitle() string {
	width, _, _ := term.GetSize(uintptr(os.Stdout.Fd()))
	titleStyle := lipgloss.NewStyle().Foreground(TitleColor).Bold(true).Width(width).Align(lipgloss.Center)

	// OSC 8 hyperlink: \033]8;;URL\033\\LABEL\033]8;;\033\\
	link := "\033]8;;" + repoURL + "\033\\" + repoURL + "\033]8;;\033\\"
	linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(width).Align(lipgloss.Center)

	return titleStyle.Render(title) + "\n" + linkStyle.Render(link)
}
