package main

import (
	"fmt"
	"os"

	"UniGrades/internal/api"
	"UniGrades/internal/screens/picker"
	"UniGrades/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	api.Run()

	tableStr := tui.RenderTable(tui.DefaultColor)

	p := tea.NewProgram(picker.InitialModel(tableStr))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
