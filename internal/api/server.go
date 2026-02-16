// Package api provides HTTP server functionality for the UniGrades application,
// offering RESTful endpoints for managing courses via a web API.
package api

import (
	// Standard library imports
	"encoding/json" // JSON encoding/decoding
	"fmt"           // Formatted I/O
	"net/http"      // HTTP server and handlers

	// MongoDB driver
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// mongoClient is a module-level variable holding the MongoDB client connection.
// It is initialized via InitServer() and used by HTTP handlers.
var mongoClient *mongo.Client

// InitServer initializes the HTTP server with the provided MongoDB client.
// This must be called before starting the server to ensure database operations work.
func InitServer(client *mongo.Client) {
	mongoClient = client
}

// handleCreateCourse handles HTTP POST requests to /courses for creating new courses.
// Expects a JSON body with course information and returns the created course's ID.
func handleCreateCourse(w http.ResponseWriter, r *http.Request) {
	// Check that the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body into a Course struct
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add the course to the database
	id, err := AddCourse(mongoClient, course)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add course: %v", err), http.StatusInternalServerError)
		return
	}

	// Return a JSON response with the new course ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "Course created successfully",
	})
}

// handleGetCourses handles HTTP GET requests to /courses for retrieving all courses.
func handleGetCourses(w http.ResponseWriter, r *http.Request) {
	// Check that the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all courses from the database
	courses := GetAllCourses(mongoClient)

	// Return the courses as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// StartServer initializes and starts the HTTP server on the specified port.
// Registers HTTP handlers for course operations and handles errors.
func StartServer(port string) {
	// Register HTTP handlers for course endpoints
	http.HandleFunc("/courses", handleCreateCourse)
	http.HandleFunc("/courses", handleGetCourses)

	// Start the HTTP server
	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
