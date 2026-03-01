package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres dbname=songs_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer db.Close()

	// i will check if i have a conn to the db
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Connected successfully to PostgreSQL!")

	createSongsTable(db)
	createFingerprintsTable(db)
}

func createSongsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,	
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		album VARCHAR(255) NOT NULL,
		duration INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating songs table:", err)
	}
}

func createFingerprintsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS fingerprints (
		song_id INT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
		hash BIGINT NOT NULL,
		time_offset INT NOT NULL,
		Primary Key (hash, time_offset, song_id)
		)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating fingerprints table:", err)
	}
}
