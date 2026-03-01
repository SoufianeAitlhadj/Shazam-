package web

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"shazam/internal/audio"
	"shazam/internal/fingerprint"
	"shazam/internal/repository/postgres"
)

func UploadHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// GET → render upload form
		if r.Method == http.MethodGet {
			RenderUploadPage(w, "")
			return
		}

		start := time.Now()

		file, _, err := r.FormFile("audio")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		tmpFile, err := os.CreateTemp("", "upload-*.wav")
		if err != nil {
			http.Error(w, "Failed to create temp file", 500)
			return
		}
		defer os.Remove(tmpFile.Name())

		io.Copy(tmpFile, file)
		tmpFile.Close()

		samples, sr, err := audio.ReadWav(tmpFile.Name())
		if err != nil {
			http.Error(w, "Invalid WAV file", 500)
			return
		}

		queryFP := fingerprint.Extract(samples, sr)

		matches, err := postgres.FindMatches(db, queryFP)
		if err != nil || len(matches) == 0 {
			RenderUploadPage(w, "<h2>No match found</h2>")
			return
		}

		var resultHTML string

		for i, match := range matches {

			song, err := postgres.GetSongByID(db, match.SongID)
			if err != nil {
				continue
			}

			embedLink := ConvertToEmbed(song.YoutubeLink)

			resultHTML += BuildRankedResultHTML(
				i+1,
				song.Title,
				song.Artist,
				song.Album,
				song.YoutubeLink,
				embedLink,
				match.Score,
			)
		}

		elapsed := time.Since(start)

		resultHTML += fmt.Sprintf(
			`<div style="margin-top:20px;"><strong>Processing time:</strong> %d ms</div>`,
			elapsed.Milliseconds(),
		)

		RenderUploadPage(w, resultHTML)
	}
}
