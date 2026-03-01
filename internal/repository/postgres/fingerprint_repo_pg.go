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

package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"shazam/internal/fingerprint"
)

func FindMatches(db *sql.DB, queryFP []fingerprint.Fingerprint) (int, error) {

	if len(queryFP) == 0 {
		return 0, nil
	}

	// Collect unique hashes
	hashSet := make(map[int64]struct{})
	for _, fp := range queryFP {
		hashSet[fp.Hash] = struct{}{}
	}

	var hashes []string
	for h := range hashSet {
		hashes = append(hashes, fmt.Sprintf("%d", h))
	}

	query := fmt.Sprintf(`
		SELECT song_id, COUNT(*) as match_count
		FROM fingerprints
		WHERE hash IN (%s)
		GROUP BY song_id
		ORDER BY match_count DESC
		LIMIT 1
	`, strings.Join(hashes, ","))

	var songID int
	var count int

	err := db.QueryRow(query).Scan(&songID, &count)
	if err != nil {
		return 0, nil
	}

	return songID, nil
}