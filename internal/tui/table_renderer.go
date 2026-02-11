package tui

import (
	"UniGrades/internal/api"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongoDBClient() *mongo.Client {
	godotenv.Load(".env")
	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func RenderTable(uniColor lipgloss.Color) string {

	client := setupMongoDBClient()
	headers := api.GetTableHeaders(client)
	courses := api.GetAllCourses(client)

	rows := make([][]string, 0, len(courses))
	for _, course := range courses {
		row := make([]string, 0, len(headers))
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", course[h]))
		}
		rows = append(rows, row)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(uniColor)).
		StyleFunc(TableStyleFunc(uniColor)).
		Headers(headers...).
		Rows(rows...)

	// You can also add tables row-by-row
	// t.Row("English", "You look absolutely fabulous.", "How's it going?")

	return t.Render()
}
