package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yoockh/dbyoc/config"
	"github.com/yoockh/dbyoc/db/nosql"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create a new MongoDB client
	mongoClient, err := nosql.NewMongoClient(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}

	// Example operation: Insert a document
	document := map[string]interface{}{
		"name":  "John Doe",
		"email": "john.doe@example.com",
		"age":   30,
	}

	err = mongoClient.Insert("users", document)
	if err != nil {
		log.Fatalf("Error inserting document: %v", err)
	}

	fmt.Println("Document inserted successfully")

	// Example operation: Find a document
	var result map[string]interface{}
	err = mongoClient.Find("users", map[string]interface{}{"name": "John Doe"}, &result)
	if err != nil {
		log.Fatalf("Error finding document: %v", err)
	}

	fmt.Printf("Found document: %+v\n", result)

	// Close the MongoDB client
	defer mongoClient.Close()

	// Wait for a while before exiting
	time.Sleep(2 * time.Second)
}
