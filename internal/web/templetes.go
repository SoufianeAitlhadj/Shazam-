package web

import (
	"fmt"
	"net/http"
)

func RenderUploadPage(w http.ResponseWriter, result string) {
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Shazam MVP</title>
	</head>
	<body style="font-family: Arial; text-align:center;">
		<h1>🎧 Shazam MVP</h1>
		<form method="POST" enctype="multipart/form-data">
			<input type="file" name="audio" accept=".wav" required>
			<br><br>
			<button type="submit">Upload & Identify</button>
		</form>
		<br><br>
		%s
	</body>
	</html>
	`, result)

	w.Write([]byte(html))
}

func BuildResultHTML(title, artist, album, link, embed string) string {
	return fmt.Sprintf(`
		<h2>🎵 Match Found!</h2>
		<p><strong>Title:</strong> %s</p>
		<p><strong>Artist:</strong> %s</p>
		<p><strong>Album:</strong> %s</p>
		<p><a href="%s" target="_blank">Open on YouTube</a></p>
		<iframe width="560" height="315"
			src="%s"
			frameborder="0"
			allowfullscreen>
		</iframe>
	`, title, artist, album, link, embed)
}
