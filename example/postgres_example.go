package main

import (
	"fmt"
	"log"

	"github.com/yoockh/dbyoc/config"
	"github.com/yoockh/dbyoc/db/sql"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create a new PostgreSQL client
	dbClient, err := sql.NewPostgresClient(cfg.Database)
	if err != nil {
		log.Fatalf("Error creating PostgreSQL client: %v", err)
	}
	defer dbClient.Close()

	// Example query
	rows, err := dbClient.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	// Process query results
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("User: %d, Name: %s\n", id, name)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error with rows: %v", err)
	}
}
