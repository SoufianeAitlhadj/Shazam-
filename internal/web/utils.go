package web

import "strings"

func ConvertToEmbed(link string) string {
	if strings.Contains(link, "watch?v=") {
		videoID := strings.Split(link, "watch?v=")[1]
		return "https://www.youtube.com/embed/" + videoID
	}
	return link
}
