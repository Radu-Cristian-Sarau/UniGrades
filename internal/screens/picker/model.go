package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"UniGrades/internal/university"
)

type Model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func InitialModel() Model {
	return Model{
		choices:  university.Names(),
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

var uniColors = university.ColorMap()

func (m Model) View() string {
	s := "Select university: \n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		style := lipgloss.NewStyle().Foreground(uniColors[choice])
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, style.Render(choice))
	}

	s += "\nPress Q or Ctrl + C to quit.\n"

	return s
}
