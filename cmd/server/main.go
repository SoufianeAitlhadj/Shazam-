package main

import (
	"fmt"
	"log"
	"net/http"

	"shazam/internal/repository/postgres"
	"shazam/internal/web"
)

func main() {
	//configuration
	connStr := "host=localhost port=5432 user=postgres dbname=songs_db sslmode=disable"
	port := ":8080"

	//database Initialization
	db := postgres.NewDB(connStr)
	defer db.Close()

	postgres.RunMigrations(db)

	//HTTP Routing
	http.HandleFunc("/", web.UploadHandler(db))

	//start Server
	fmt.Println("Shazam server running on port", port)
	fmt.Println("Open http://localhost" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}
