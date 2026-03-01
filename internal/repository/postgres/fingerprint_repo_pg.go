package postgres

import (
	"database/sql"

	"shazam/internal/fingerprint"
)

func InsertFingerprints(db *sql.DB, songID int, fps []fingerprint.Fingerprint) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO fingerprints (song_id, hash, time_offset)
		 VALUES ($1, $2, $3)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, fp := range fps {
		_, err := stmt.Exec(songID, fp.Hash, fp.TimeOffset)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func FindMatches(db *sql.DB, queryFP []fingerprint.Fingerprint) (int, error) {

	scores := make(map[int]int)

	for _, fp := range queryFP {

		rows, err := db.Query(
			`SELECT song_id FROM fingerprints WHERE hash = $1`,
			fp.Hash,
		)
		if err != nil {
			return 0, err
		}

		for rows.Next() {
			var songID int
			if err := rows.Scan(&songID); err != nil {
				rows.Close()
				return 0, err
			}
			scores[songID]++
		}
		rows.Close()
	}

	// find best match
	bestSongID := 0
	maxScore := 0

	for songID, score := range scores {
		if score > maxScore {
			maxScore = score
			bestSongID = songID
		}
	}

	return bestSongID, nil
}
