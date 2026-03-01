package web

import (
	"database/sql"
	"io"
	"net/http"
	"os"

	"shazam/internal/audio"
	"shazam/internal/fingerprint"
	"shazam/internal/repository/postgres"
)

func UploadHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			RenderUploadPage(w, "")
			return
		}

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

		matchedID, err := postgres.FindMatches(db, queryFP)
		if err != nil || matchedID == 0 {
			RenderUploadPage(w, "<h2>No match found</h2>")
			return
		}

		song, err := postgres.GetSongByID(db, matchedID)
		if err != nil {
			http.Error(w, "Failed to fetch song", 500)
			return
		}

		embedLink := ConvertToEmbed(song.YoutubeLink)

		resultHTML := BuildResultHTML(song.Title, song.Artist, song.Album, song.YoutubeLink, embedLink)

		RenderUploadPage(w, resultHTML)
	}
}
