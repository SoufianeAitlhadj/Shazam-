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
	<title>Shazam Clone</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background: linear-gradient(to right, #141e30, #243b55);
			color: white;
			text-align: center;
			padding-top: 50px;
		}

		h1 {
			font-size: 40px;
			margin-bottom: 30px;
		}

		.upload-box {
			background: rgba(255,255,255,0.1);
			padding: 30px;
			border-radius: 10px;
			display: inline-block;
			box-shadow: 0px 0px 20px rgba(0,0,0,0.3);
		}

		input[type="file"] {
			margin-bottom: 20px;
			color: white;
		}

		button {
			background-color: #1db954;
			border: none;
			padding: 12px 25px;
			color: white;
			font-size: 16px;
			border-radius: 5px;
			cursor: pointer;
		}

		button:hover {
			background-color: #17a74a;
		}

		.result-card {
			margin-top: 40px;
			background: white;
			color: black;
			padding: 20px;
			border-radius: 10px;
			width: 60%%;
			margin-left: auto;
			margin-right: auto;
			box-shadow: 0px 0px 20px rgba(0,0,0,0.4);
		}

		iframe {
			margin-top: 20px;
			border-radius: 10px;
		}

		.loading {
			display: none;
			margin-top: 20px;
			font-size: 18px;
		}
	</style>

	<script>
		function showLoading() {
			document.getElementById("loading").style.display = "block";
		}
	</script>

</head>
<body>

	<h1>🎧 Shazam MVP</h1>

	<div class="upload-box">
		<form method="POST" enctype="multipart/form-data" onsubmit="showLoading()">
			<input type="file" name="audio" accept=".wav" required>
			<br>
			<button type="submit">Upload & Identify</button>
		</form>
		<div id="loading" class="loading">
			⏳ Processing audio... please wait
		</div>
	</div>

	%s

</body>
</html>
`, result)

	w.Write([]byte(html))
}

func BuildResultHTML(title, artist, album, link, embed string) string {
	return fmt.Sprintf(`
	<div class="result-card">
		<h2>🎵 Match Found!</h2>
		<p><strong>Title:</strong> %s</p>
		<p><strong>Artist:</strong> %s</p>
		<p><strong>Album:</strong> %s</p>
		<p><a href="%s" target="_blank">▶ Open on YouTube</a></p>
		<iframe width="560" height="315"
			src="%s"
			frameborder="0"
			allowfullscreen>
		</iframe>
	</div>
	`, title, artist, album, link, embed)
}
