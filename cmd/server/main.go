package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"shazam/internal/audio"
	"shazam/internal/fingerprint"
	"shazam/internal/repository/postgres"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres dbname=songs_db sslmode=disable"

	db := postgres.NewDB(connStr)
	defer db.Close()

	postgres.RunMigrations(db)

	//I did this just to clear old data before indexing (safe for MVP)
	clearDatabase(db)

	songsDir := "./songs/wav_songs"

	files, err := os.ReadDir(songsDir)
	if err != nil {
		log.Fatal(err)
	}

	songMetadata := map[string]struct {
		Title       string
		Artist      string
		Album       string
		YouTubeLink string
	}{
		"Daft Punk - Instant Crush (Official Video) ft. Julian Casablancas.wav": {
			Title:       "Instant Crush",
			Artist:      "Daft Punk",
			Album:       "Random Access Memories",
			YouTubeLink: "https://www.youtube.com/watch?v=a5uQMwRMHcs",
		},
		"Daft Punk - Lose Yourself to Dance (Official Version).wav": {
			Title:       "Lose Yourself to Dance",
			Artist:      "Daft Punk",
			Album:       "Random Access Memories",
			YouTubeLink: "https://www.youtube.com/watch?v=NF-kLy44Hls",
		},
		"Daft Punk - One More Time (Official Video).wav": {
			Title:       "One More Time",
			Artist:      "Daft Punk",
			Album:       "Discovery",
			YouTubeLink: "https://www.youtube.com/watch?v=FGBhQbmPwH8",
		},
		"Kendrick Lamar - Count Me Out (Official Audio).wav": {
			Title:       "Count Me Out",
			Artist:      "Kendrick Lamar",
			Album:       "Mr. Morale & The Big Steppers",
			YouTubeLink: "https://www.youtube.com/watch?v=6nTcdw7bVdc",
		},
		"Madd - Motto.wav": {
			Title:       "Motto",
			Artist:      "Madd",
			Album:       "Sensus",
			YouTubeLink: "https://www.youtube.com/watch?v=vVLV9sJopBU",
		},
		"Pink Floyd - Money (Official Music Video).wav": {
			Title:       "Money",
			Artist:      "Pink Floyd",
			Album:       "The Dark Side of the Moon",
			YouTubeLink: "https://www.youtube.com/watch?v=-0kcet4aPpQ",
		},
		"Tame Impala - Let It Happen (Official Video).wav": {
			Title:       "Let It Happen",
			Artist:      "Tame Impala",
			Album:       "Currents",
			YouTubeLink: "https://www.youtube.com/watch?v=pFptt7Cargc",
		},
		"Tame Impala - New Person, Same Old Mistakes (Audio).wav": {
			Title:       "New Person, Same Old Mistakes",
			Artist:      "Tame Impala",
			Album:       "Currents",
			YouTubeLink: "https://www.youtube.com/watch?v=_9bw_VtMUGA",
		},
	}

	for _, file := range files {

		if filepath.Ext(file.Name()) != ".wav" {
			continue
		}

		meta, exists := songMetadata[file.Name()]
		if !exists {
			fmt.Println("⚠ Skipping (no metadata):", file.Name())
			continue
		}

		fullPath := filepath.Join(songsDir, file.Name())

		fmt.Println("Processing:", meta.Title)

		samples, sr, err := audio.ReadWav(fullPath)
		if err != nil {
			log.Fatal(err)
		}

		fps := fingerprint.Extract(samples, sr)

		duration := len(samples) / sr

		songID, err := postgres.InsertSong(
			db,
			meta.Title,
			meta.Artist,
			meta.Album,
			duration,
			meta.YouTubeLink,
		)
		if err != nil {
			log.Fatal(err)
		}

		err = postgres.InsertFingerprints(db, songID, fps)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted:", meta.Title)
	}

	fmt.Println("✅ Database indexing complete")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM songs").Scan(&count)
	fmt.Println("Songs currently in DB:", count)
}

func clearDatabase(db *sql.DB) {
	_, err := db.Exec(`TRUNCATE fingerprints, songs RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal("Failed to clear database:", err)
	}
	fmt.Println("🧹 Database cleared")
}
