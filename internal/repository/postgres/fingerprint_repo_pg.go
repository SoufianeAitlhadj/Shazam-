package postgres

import (
	"database/sql"
	"sort"

	"shazam/internal/fingerprint"

	"github.com/lib/pq"
)

// INSERT FINGERPRINTS

func InsertFingerprints(db *sql.DB, songID int, fps []fingerprint.Fingerprint) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO fingerprints (song_id, hash, time_offset)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		tx.Rollback()
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

type MatchResult struct {
	SongID int
	Score  int
}

func FindMatches(db *sql.DB, queryFP []fingerprint.Fingerprint) ([]MatchResult, error) {

	if len(queryFP) == 0 {
		return nil, nil
	}

	// Build map: hash -> list of query time offsets
	queryMap := make(map[int64][]int)
	hashSet := make(map[int64]struct{})

	for _, fp := range queryFP {
		queryMap[fp.Hash] = append(queryMap[fp.Hash], fp.TimeOffset)
		hashSet[fp.Hash] = struct{}{}
	}

	var hashes []int64
	for h := range hashSet {
		hashes = append(hashes, h)
	}

	query := `
		SELECT song_id, hash, time_offset
		FROM fingerprints
		WHERE hash = ANY($1)
	`

	rows, err := db.Query(query, pq.Array(hashes))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Structure:
	// songID -> delta -> count
	votes := make(map[int]map[int]int)

	for rows.Next() {
		var songID int
		var hash int64
		var dbOffset int

		err := rows.Scan(&songID, &hash, &dbOffset)
		if err != nil {
			return nil, err
		}

		queryOffsets := queryMap[hash]

		for _, qOffset := range queryOffsets {
			delta := dbOffset - qOffset

			if votes[songID] == nil {
				votes[songID] = make(map[int]int)
			}
			votes[songID][delta]++
		}
	}

	// Now compute best score per song
	var results []MatchResult

	for songID, deltaMap := range votes {
		maxVotes := 0
		for _, count := range deltaMap {
			if count > maxVotes {
				maxVotes = count
			}
		}
		results = append(results, MatchResult{
			SongID: songID,
			Score:  maxVotes,
		})
	}

	// Sort descending by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Return top 3
	if len(results) > 3 {
		results = results[:3]
	}

	return results, nil
}
