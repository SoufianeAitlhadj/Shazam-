package main

import (
	"fmt"
	"log"
	"net/http"

	"shazam/internal/repository/postgres"
	"shazam/internal/web"
)

func main() {
	// Database connection string
	connStr := "host=localhost port=5432 user=postgres dbname=songs_db sslmode=disable"

	// Connect to DB
	db := postgres.NewDB(connStr)
	defer db.Close()

	// Ensure tables + indexes exist
	postgres.RunMigrations(db)

	// Register routes
	http.HandleFunc("/", web.UploadHandler(db))

	fmt.Println("🚀 Shazam server running at http://localhost:8080")

	// Start HTTP server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
