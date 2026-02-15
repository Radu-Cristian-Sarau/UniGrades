package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var mongoClient *mongo.Client

// InitServer initializes the HTTP server with the MongoDB client
func InitServer(client *mongo.Client) {
	mongoClient = client
}

// handleCreateCourse handles POST /courses requests
func handleCreateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := AddCourse(mongoClient, course)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add course: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Course created successfully",
	})
}

// handleGetCourses handles GET /courses requests
func handleGetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	courses := GetAllCourses(mongoClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// StartServer starts the HTTP server on the given port
func StartServer(port string) {
	http.HandleFunc("/courses", handleCreateCourse)
	http.HandleFunc("/courses", handleGetCourses)

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
