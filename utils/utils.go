package utils

import "fmt"

func ExtractTweet(text string) (tweet string, newPos int) {
	if len(text) <= 280 {
		return text, 0
	}

	for pos, rune := range text {
		if rune == '.' {
			if pos <= 280 {
				newPos = pos + 1
			} else if pos > 280 {
				break
			}
		}
	}

	if newPos == 0 {
		return text[:280], 281
	} else {
		return text[:newPos], newPos + 1
	}
}

func GenerateTweetURL(username string, id int64) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%d", username, id)
}