package utils

import "fmt"

func ExtractTweet(text string) (string, int) {
	if len(text) <= 280 {
		return text, 0
	}

	var lastPunctuation, lastWhitespace int
	for pos, rune := range text {
		if pos > 280 {
			break
		}
		if rune == '.' || rune == '!' || rune == '?' {
			lastPunctuation = pos + 1
		} else if rune == ' ' {
			lastWhitespace = pos + 1
		}
	}

	if lastPunctuation > 0 {
		return text[:lastPunctuation], lastPunctuation
	} else if lastWhitespace > 0 {
		return text[:lastWhitespace], lastWhitespace
	} else {
		return text[:280], 280
	}
}

func GenerateTweetURL(username string, id int64) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%d", username, id)
}
