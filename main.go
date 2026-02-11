package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"UniGrades/internal/screens/picker"
)

func main() {
	p := tea.NewProgram(picker.InitialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
