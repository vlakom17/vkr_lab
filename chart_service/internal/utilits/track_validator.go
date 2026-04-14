package utilits

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	maxLength  = 110
	minLetters = 1
)

func ValidateTrackInput(artist, title string) error {
	artist = strings.TrimSpace(artist)
	title = strings.TrimSpace(title)

	if artist == "" || title == "" {
		return fmt.Errorf("artist and title are required")
	}

	if len([]rune(artist)) > maxLength {
		return fmt.Errorf("artist is too long")
	}
	if len([]rune(title)) > maxLength {
		return fmt.Errorf("title is too long")
	}

	if countLetters(artist) < minLetters {
		return fmt.Errorf("artist must contain at least %d letters", minLetters)
	}
	if countLetters(title) < minLetters {
		return fmt.Errorf("title must contain at least %d letters", minLetters)
	}

	if looksLikeURL(artist) || looksLikeURL(title) {
		return fmt.Errorf("links are not allowed")
	}

	if looksLikeHTML(artist) || looksLikeHTML(title) {
		return fmt.Errorf("html content is not allowed")
	}

	if isSuspicious(artist) {
		return fmt.Errorf("artist looks suspicious")
	}
	if isSuspicious(title) {
		return fmt.Errorf("title looks suspicious")
	}

	return nil
}

func countLetters(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			count++
		}
	}
	return count
}

func looksLikeURL(s string) bool {
	re := regexp.MustCompile(`(?i)(https?://|www\.)`)
	return re.MatchString(s)
}

func looksLikeHTML(s string) bool {
	re := regexp.MustCompile(`(?i)<[^>]+>`)
	return re.MatchString(s)
}

func isSuspicious(s string) bool {
	total := 0
	special := 0

	for _, r := range s {
		if unicode.IsSpace(r) {
			continue
		}
		total++

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			special++
		}
	}

	if total == 0 {
		return true
	}

	return special > 2*total/3
}
