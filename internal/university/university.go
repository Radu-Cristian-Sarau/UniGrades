// Package university provides university information and styling for the UniGrades application.
package university

import "github.com/charmbracelet/lipgloss"

// University represents a university entity with its name and brand color.
type University struct {
	// Name is the display name of the university (e.g., "TU/e")
	Name string
	// Color is the brand/theme color for the university
	Color lipgloss.Color
}

// All returns a slice of all available universities with their configurations.
func All() []University {
	return []University{
		{Name: "TU/e", Color: lipgloss.Color("#c81919")}, // Eindhoven - Red
		{Name: "TUD", Color: lipgloss.Color("#00a0da")},  // Delft - Blue
		{Name: "TUM", Color: lipgloss.Color("#0066c1")},  // Munich - Dark Blue
	}
}

// Names returns just the names of all available universities.
func Names() []string {
	unis := All()
	names := make([]string, len(unis))
	for i, u := range unis {
		names[i] = u.Name
	}
	return names
}

// ColorMap returns a map from university names to their brand colors.
// This allows quick color lookups by university name.
func ColorMap() map[string]lipgloss.Color {
	colors := make(map[string]lipgloss.Color)
	for _, u := range All() {
		colors[u.Name] = u.Color
	}
	return colors
}
