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
	fmt.Println(tui.RenderTable())
	p := tea.NewProgram(picker.InitialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
