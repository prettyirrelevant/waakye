package utils

import (
	"regexp"
	"strings"
)

// CleanTrackTitle cleans up a track title by removing any text enclosed in square brackets or parentheses,
// and adding the word "Remix" to the end of the title if it contains the word "remix" (case-insensitive).
func CleanTrackTitle(title string) string {
	re := regexp.MustCompile(`[\(\[].*?[\)\]]`)

	cleanedTitle := strings.TrimSpace(re.ReplaceAllString(title, ""))
	if strings.Contains(strings.ToLower(title), "remix") {
		cleanedTitle += " Remix"
	}

	return cleanedTitle
}
