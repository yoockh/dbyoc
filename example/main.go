package main

import (
	"fmt"
	"log"

	"github.com/yoockh/dbyoc/config"
	"github.com/yoockh/dbyoc/db"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database client
	client, err := db.NewClient(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing database client: %v", err)
	}
	defer client.Close()

	// Example query
	result, err := client.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	fmt.Println("Query result:", result)
}
