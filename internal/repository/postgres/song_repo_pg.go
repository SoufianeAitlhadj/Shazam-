package postgres

import "database/sql"

func InsertSong(
	db *sql.DB,
	title, artist, album string,
	duration int,
	youtubeLink string,
) (int, error) {

	var songID int

	err := db.QueryRow(
		`INSERT INTO songs (title, artist, album, duration, youtube_link)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		title, artist, album, duration, youtubeLink,
	).Scan(&songID)

	if err != nil {
		return 0, err
	}

	return songID, nil
}

type Song struct {
	ID          int
	Title       string
	Artist      string
	Album       string
	Duration    int
	YoutubeLink string
}

func GetSongByID(db *sql.DB, id int) (*Song, error) {
	var song Song

	err := db.QueryRow(
		`SELECT id, title, artist, album, duration, youtube_link
		 FROM songs
		 WHERE id = $1`,
		id,
	).Scan(
		&song.ID,
		&song.Title,
		&song.Artist,
		&song.Album,
		&song.Duration,
		&song.YoutubeLink,
	)

	if err != nil {
		return nil, err
	}

	return &song, nil
}
