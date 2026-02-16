// Package api provides MongoDB database operations for managing course information.
// It handles CRUD operations (Create, Read, Update, Delete) for courses stored in MongoDB,
// including operations like retrieving all courses, getting table headers, adding courses,
// updating course details, and deleting courses.
package api

import (
	// Standard library imports for core functionality
	"context"       // Used for context management in MongoDB operations
	"encoding/json" // JSON marshaling/unmarshaling utilities
	"fmt"           // Formatted I/O package
	"log"           // Logging utilities
	"os"            // Operating system functionality
	"strconv"       // String conversion utilities

	// Third-party packages
	"github.com/joho/godotenv"                     // Loads environment variables from .env files
	"go.mongodb.org/mongo-driver/v2/bson"          // BSON encoding/decoding for MongoDB
	"go.mongodb.org/mongo-driver/v2/mongo"         // MongoDB driver
	"go.mongodb.org/mongo-driver/v2/mongo/options" // MongoDB connection options
)

// GetAllCourses retrieves all course documents from the MongoDB database.
// It queries the "TUe" collection in the "CourseInfo" database and returns
// all documents as a slice of BSON maps.
//
// Parameters:
//
//	client: MongoDB client connection
//
// Returns:
//
//	A slice of bson.M (BSON maps) containing all course documents.
//	Panics if the database query fails.
func GetAllCourses(client *mongo.Client) []bson.M {
	// Access the "TUe" collection from the "CourseInfo" database
	coll := client.Database("CourseInfo").Collection("TUe")

	// Execute a query to find all documents (empty filter returns all documents)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())

	// Decode all cursor results into a slice of BSON maps
	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// Optional debug code (commented out) that could print JSON representation of results
	/* jsonData, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData) */

	return results
}

// GetTableHeaders retrieves the field names of the first course document in the database.
// This is useful for determining the structure/schema of course documents and
// extracting column headers for display in tables. The _id field is excluded.
//
// Parameters:
//
//	client: MongoDB client connection
//
// Returns:
//
//	A slice of strings containing the field names from the first document.
//	Returns nil if no courses are found in the database.
func GetTableHeaders(client *mongo.Client) []string {
	// Access the "TUe" collection from the "CourseInfo" database
	coll := client.Database("CourseInfo").Collection("TUe")

	var result bson.D
	// Configure projection to exclude the MongoDB _id field
	opts := options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 0}})

	// Retrieve the first document from the collection
	err := coll.FindOne(context.TODO(), bson.D{}, opts).Decode(&result)

	// Handle case where no documents exist in the collection
	if err == mongo.ErrNoDocuments {
		fmt.Println("Cannot retrieve table headers: No courses were found in the database.")
		return nil
	}
	if err != nil {
		panic(err)
	}

	// Extract field names from the BSON document
	headers := make([]string, 0, len(result))
	for _, elem := range result {
		headers = append(headers, elem.Key)
	}

	return headers
}

// getCourseDataByName queries the database for a course by its name and prints
// the course data as formatted JSON. This is a private helper function (lowercase name).
//
// Parameters:
//
//	client: MongoDB client connection
//	name: The name of the course to search for
//
// Note: This function prints to stdout rather than returning a value. It's typically
// used for debugging or manual data inspection.
func getCourseDataByName(client *mongo.Client, name string) {
	// Access the "TUe" collection from the "CourseInfo" database
	coll := client.Database("CourseInfo").Collection("TUe")

	var result bson.M
	// Query for a single document matching the given course name
	err := coll.FindOne(context.TODO(), bson.D{{"Name", name}}).
		Decode(&result)

	// Handle case where no document with the given name exists
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No course was found with the name %s\n", name)
		return
	}
	if err != nil {
		panic(err)
	}

	// Marshal the result to formatted JSON for pretty printing
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}

	// Print the formatted JSON to stdout
	fmt.Printf("%s\n", jsonData)
}

// Course represents a university course with its core information.
// This struct maps to course documents in the MongoDB database with BSON tags
// defining how fields are serialized/deserialized.
type Course struct {
	// Name is the identifier for the course (e.g., "DZC10_Game_Design_I")
	Name string `bson:"Name"`

	// Year is the academic year in which the course was taken
	Year int `bson:"Year"`

	// Grade is the numerical grade/mark received for the course
	Grade float64 `bson:"Grade"`

	// ECTS is the number of European Credit Transfer System points earned
	ECTS int `bson:"ECTS"`
}

// AddCourse inserts a new course document into the MongoDB database.
// It creates a new entry in the "TUe" collection under the "CourseInfo" database.
//
// Parameters:
//
//	client: MongoDB client connection
//	course: A Course struct containing the course information to be added
//
// Returns:
//
//	A string containing the MongoDB ObjectID of the newly inserted document (in hex format).
//	An error if insertion fails, wrapped with context about the failure.
func AddCourse(client *mongo.Client, course Course) (string, error) {
	// Access the "TUe" collection from the "CourseInfo" database
	coll := client.Database("CourseInfo").Collection("TUe")

	// Insert the course document into the collection
	result, err := coll.InsertOne(context.TODO(), course)
	if err != nil {
		return "", fmt.Errorf("failed to insert course: %w", err)
	}

	// Return the automatically generated MongoDB ObjectID as a hexadecimal string
	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

// DeleteCourse removes a course from the MongoDB database by matching on its name.
// If no course with the specified name exists, an error is returned.
//
// Parameters:
//
//	client: MongoDB client connection
//	courseName: The exact name of the course to delete
//
// Returns:
//
//	An error if deletion fails or if the course is not found. Returns nil on success.
func DeleteCourse(client *mongo.Client, courseName string) error {
	// Access the "TUe" collection from the "CourseInfo" database
	coll := client.Database("CourseInfo").Collection("TUe")

	// Execute a delete operation targeting the course with the matching name
	result, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "Name", Value: courseName}})
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	// Check if the course was actually found and deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("course '%s' not found", courseName)
	}

	return nil
}

// UpdateCourse modifies a single field of a course document in the MongoDB database.
// It performs type validation and conversion based on the field being updated.
// Supported fields are: Name (string), Grade (float), Year (int), and ECTS (int).
//
// Parameters:
//
//	client: MongoDB client connection
//	courseName: The exact name of the course to update
//	field: The field name to update (Name, Grade, Year, or ECTS)
//	value: The new value as a string (will be converted to appropriate type)
//
// Returns:
//
//	An error if the update fails, the field is invalid, the value cannot be converted,
//	or the course is not found. Returns nil on success.
func UpdateCourse(client *mongo.Client, courseName, field, value string) error {
	// Access the "TUe" collection from the "CourseInfo" database
	coll := client.Database("CourseInfo").Collection("TUe")

	// Convert the string value to the appropriate type based on the field
	var updateValue interface{}
	var err error

	switch field {
	case "Name":
		// Name field is stored as a string, use value as-is
		updateValue = value
	case "Grade":
		// Grade field must be converted to float64
		updateValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("Grade must be a number")
		}
	case "Year":
		// Year field must be converted to integer
		year, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Year must be an integer")
		}
		updateValue = year
	case "ECTS":
		// ECTS field must be converted to integer
		ects, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("ECTS must be an integer")
		}
		updateValue = ects
	default:
		// Field name is not recognized
		return fmt.Errorf("invalid field: %s. Valid fields are: Name, Year, Grade, ECTS", field)
	}

	// Execute the update operation on the document matching the course name
	result, err := coll.UpdateOne(
		context.TODO(),
		bson.D{{Key: "Name", Value: courseName}}, // Filter: match by course name
		bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: updateValue}}}}, // Update: set the field to new value
	)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	// Check if the course was actually found (MatchedCount > 0 means a document was matched)
	if result.MatchedCount == 0 {
		return fmt.Errorf("course '%s' not found", courseName)
	}

	return nil
}

// Run initializes and manages the MongoDB database connection.
// It loads the MongoDB URI from environment variables and establishes a connection.
// This function serves as a setup routine and should be called before performing
// any database operations. The MongoDB client connection is properly deferred for cleanup.
//
// Environment Requirements:
//
//	MONGODB_URI: Environment variable containing the MongoDB connection string.
//	             If not set, the function will terminate with a fatal error.
//
// Note: This function currently contains commented-out debug operations that
// could be used to test database connectivity (getCourseDataByName, GetAllCourses).
func Run() {
	// Load environment variables from the .env file in the current directory
	godotenv.Load(".env")

	// Retrieve the MongoDB connection URI from environment variables
	uri := os.Getenv("MONGODB_URI")

	// Documentation reference for MongoDB Go driver
	docs := "www.mongodb.com/docs/drivers/go/current/"

	// Validate that the MONGODB_URI environment variable is set
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	// Establish a connection to the MongoDB server using the provided URI
	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ensure the MongoDB connection is properly closed when the function exits
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Optional debug operations (currently commented out)
	// Uncomment to test database connectivity and retrieve sample data:
	// getCourseDataByName(client, "DZC10_Game_Design_I")
	// fmt.Println(GetAllCourses(client))

	// Retrieve and set up table headers from the database
	GetTableHeaders(client)
}
