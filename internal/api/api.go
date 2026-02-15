package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetAllCourses(client *mongo.Client) []bson.M {
	coll := client.Database("CourseInfo").Collection("TUe")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	/* jsonData, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData) */
	return results
}

func GetTableHeaders(client *mongo.Client) []string {
	coll := client.Database("CourseInfo").Collection("TUe")
	var result bson.D
	opts := options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 0}})
	err := coll.FindOne(context.TODO(), bson.D{}, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("Cannot retrieve table headers: No courses were found in the database.")
		return nil
	}
	if err != nil {
		panic(err)
	}
	headers := make([]string, 0, len(result))
	for _, elem := range result {
		headers = append(headers, elem.Key)
	}
	return headers
}

func getCourseDataByName(client *mongo.Client, name string) {
	coll := client.Database("CourseInfo").Collection("TUe")
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"Name", name}}).
		Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No course was found with the name %s\n", name)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}

// Course represents the structure of a course document
type Course struct {
	Name  string  `bson:"Name"`
	Year  int     `bson:"Year"`
	Grade float64 `bson:"Grade"`
	ECTS  int     `bson:"ECTS"`
}

// AddCourse inserts a new course into the MongoDB database
func AddCourse(client *mongo.Client, course Course) (string, error) {
	coll := client.Database("CourseInfo").Collection("TUe")

	result, err := coll.InsertOne(context.TODO(), course)
	if err != nil {
		return "", fmt.Errorf("failed to insert course: %w", err)
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

// DeleteCourse removes a course from the MongoDB database by name
func DeleteCourse(client *mongo.Client, courseName string) error {
	coll := client.Database("CourseInfo").Collection("TUe")

	result, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "Name", Value: courseName}})
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("course '%s' not found", courseName)
	}

	return nil
}

// UpdateCourse updates a course field in the MongoDB database
func UpdateCourse(client *mongo.Client, courseName, field, value string) error {
	coll := client.Database("CourseInfo").Collection("TUe")

	// Convert value to appropriate type based on field
	var updateValue interface{}
	var err error

	switch field {
	case "Name":
		updateValue = value
	case "Grade":
		updateValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("Grade must be a number")
		}
	case "Year":
		year, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Year must be an integer")
		}
		updateValue = year
	case "ECTS":
		ects, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("ECTS must be an integer")
		}
		updateValue = ects
	default:
		return fmt.Errorf("invalid field: %s. Valid fields are: Name, Year, Grade, ECTS", field)
	}

	result, err := coll.UpdateOne(
		context.TODO(),
		bson.D{{Key: "Name", Value: courseName}},
		bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: updateValue}}}},
	)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("course '%s' not found", courseName)
	}

	return nil
}

func Run() {
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
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// getCourseDataByName(client, "DZC10_Game_Design_I")
	// fmt.Println(GetAllCourses(client))
	GetTableHeaders(client)
}
