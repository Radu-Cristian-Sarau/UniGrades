package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
