package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("✅ Connected successfully to PostgreSQL!")

	return db
}

func RunMigrations(db *sql.DB) {
	createSongsTable(db)
	createFingerprintsTable(db)
	createIndexes(db)
	log.Println("✅ Database migrations completed.")
}

func createSongsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		album VARCHAR(255) NOT NULL,
		duration INT NOT NULL,
		youtube_link TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating songs table:", err)
	}
}

func createFingerprintsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS fingerprints (
		song_id INT REFERENCES songs(id) ON DELETE CASCADE,
		hash BIGINT NOT NULL,
		time_offset INT NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating fingerprints table:", err)
	}
}

func createIndexes(db *sql.DB) {
	query := `
	CREATE INDEX IF NOT EXISTS idx_fingerprints_hash
	ON fingerprints(hash)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating hash index:", err)
	}
}
