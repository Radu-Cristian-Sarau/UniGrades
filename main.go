// Package main is the entry point for the UniGrades application.
// It initializes the MongoDB connection and launches the terminal UI.
package main

import (
	// Standard library imports for utility functions
	"fmt" // Formatted I/O
	"log" // Logging support
	"os"  // Operating system operations

	// Internal packages for application functionality
	"UniGrades/internal/api"            // MongoDB database operations
	"UniGrades/internal/screens/picker" // University picker screen
	"UniGrades/internal/tui"            // Terminal UI rendering

	// Third-party packages
	tea "github.com/charmbracelet/bubbletea"       // TUI framework
	"github.com/joho/godotenv"                     // Environment variable loading
	"go.mongodb.org/mongo-driver/v2/mongo"         // MongoDB client
	"go.mongodb.org/mongo-driver/v2/mongo/options" // MongoDB options
)

// main initializes the application, sets up the MongoDB connection,
// and launches the terminal UI with the university picker screen.
func main() {
	// Display the application title/banner
	fmt.Println(tui.RenderTitle())

	// Load environment variables from .env file
	godotenv.Load(".env")

	// Retrieve MongoDB connection URI from environment
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: www.mongodb.com/docs/drivers/go/current/" +
			"usage-examples/#environment-variable")
	}

	// Establish MongoDB connection
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Fetch table headers and course data from database
	headers := api.GetTableHeaders(client)
	courses := api.GetAllCourses(client)

	// Initialize and run the Bubble Tea program with the picker screen
	p := tea.NewProgram(picker.InitialModel(headers, courses, client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
