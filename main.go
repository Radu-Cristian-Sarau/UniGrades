package main

import (
	"fmt"
	"log"
	"os"

	"UniGrades/internal/api"
	"UniGrades/internal/screens/picker"
	"UniGrades/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	fmt.Println(tui.RenderTitle())

	// Set up MongoDB client once
	godotenv.Load(".env")
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: www.mongodb.com/docs/drivers/go/current/" +
			"usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Fetch data once, pass to picker
	headers := api.GetTableHeaders(client)
	courses := api.GetAllCourses(client)

	p := tea.NewProgram(picker.InitialModel(headers, courses))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
