package university

import "github.com/charmbracelet/lipgloss"

// University represents a university with its display name and brand color.
type University struct {
	Name  string
	Color lipgloss.Color
}

// All returns the list of available universities.
func All() []University {
	return []University{
		{Name: "TU/e", Color: lipgloss.Color("#c81919")},
		{Name: "TUD", Color: lipgloss.Color("#00a0da")},
		{Name: "TUM", Color: lipgloss.Color("#0066c1")},
	}
}

// Names returns just the university names.
func Names() []string {
	unis := All()
	names := make([]string, len(unis))
	for i, u := range unis {
		names[i] = u.Name
	}
	return names
}

// ColorMap returns a map from university name to its brand color.
func ColorMap() map[string]lipgloss.Color {
	colors := make(map[string]lipgloss.Color)
	for _, u := range All() {
		colors[u.Name] = u.Color
	}
	return colors
}
