package utilits

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type NormalizedTrack struct {
	Artist        string
	Title         string
	NormalizedKey string
}

func extractFeat(artist, title string) ([]string, string) {
	reFeat := regexp.MustCompile(`(?i)\b(feat|ft|featuring|with)\.?\s+([^\(\)\[\]]+)`)

	mainArtists := splitArtists(artist)
	var featArtists []string

	artistMatches := reFeat.FindAllStringSubmatch(artist, -1)
	for _, m := range artistMatches {
		if len(m) >= 3 {
			part := strings.TrimSpace(m[2])
			featArtists = append(featArtists, splitArtists(part)...)
		}
	}
	artist = reFeat.ReplaceAllString(artist, "")
	artist = removeArtistBrackets(artist)
	mainArtists = splitArtists(artist)

	titleMatches := reFeat.FindAllStringSubmatch(title, -1)
	for _, m := range titleMatches {
		if len(m) >= 3 {
			part := strings.TrimSpace(m[2])
			featArtists = append(featArtists, splitArtists(part)...)
		}
	}
	title = reFeat.ReplaceAllString(title, "")

	reEmptyBrackets := regexp.MustCompile(`\(\s*\)`)
	title = reEmptyBrackets.ReplaceAllString(title, "")

	reSpaces := regexp.MustCompile(`\s+`)
	title = reSpaces.ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)

	allArtists := append(mainArtists, featArtists...)
	allArtists = unique(allArtists)

	return allArtists, title
}

func normalizeTitle(title string) string {
	reSpaces := regexp.MustCompile(`\s+`)

	title = strings.ToLower(strings.TrimSpace(title))
	title = removeTitleNoise(title)

	title = normalizeAcousticOutsideBrackets(title)
	title = normalizeRemixOutsideBrackets(title)

	title = reSpaces.ReplaceAllString(title, " ")

	title = normalizeBrackets(title)

	title = cleanSpaces(title)

	return title
}

func normalizeBrackets(title string) string {
	re := regexp.MustCompile(`\((.*?)\)`)

	return re.ReplaceAllStringFunc(title, func(s string) string {
		content := strings.ToLower(strings.Trim(s, "()"))
		content = strings.TrimSpace(content)

		switch {
		case strings.Contains(content, "remix"):
			return " (remix)"
		case strings.Contains(content, "acoustic"):
			return " (acoustic)"
		default:
			return ""
		}
	})
}

func normalizeArtist(artists []string) string {
	var cleaned []string

	for _, a := range artists {
		a = strings.ToLower(strings.TrimSpace(a))
		if a != "" {
			cleaned = append(cleaned, a)
		}
	}

	cleaned = unique(cleaned)

	sort.Strings(cleaned)

	return strings.Join(cleaned, ", ")
}

func cleanSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

func splitArtists(s string) []string {
	s = strings.ToLower(strings.TrimSpace(s))

	s = strings.ReplaceAll(s, "&", ",")
	s = strings.ReplaceAll(s, ";", ",")

	reAnd := regexp.MustCompile(`\s+and\s+`)
	s = reAnd.ReplaceAllString(s, ",")

	reX := regexp.MustCompile(`\s+x\s+`)
	s = reX.ReplaceAllString(s, ",")

	reSpaces := regexp.MustCompile(`\s+`)
	s = reSpaces.ReplaceAllString(s, " ")

	parts := strings.Split(s, ",")

	var result []string

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func removeArtistBrackets(s string) string {
	re := regexp.MustCompile(`[\(\[].*?[\)\]]`)
	return cleanSpaces(re.ReplaceAllString(s, " "))
}

func removeTitleNoise(title string) string {
	reNoise := regexp.MustCompile(`(?i)\b(official video|official audio|lyric video|lyrics video|lyrics|audio|video|hd|hq)\b`)
	return cleanSpaces(reNoise.ReplaceAllString(title, " "))
}

func normalizeRemixOutsideBrackets(title string) string {
	re := regexp.MustCompile(`(?i)\b\w*\s*remix\b`)
	if re.MatchString(title) {
		title = re.ReplaceAllString(title, "")
		return cleanSpaces(title) + " (remix)"
	}
	return title
}

func normalizeAcousticOutsideBrackets(title string) string {
	re := regexp.MustCompile(`(?i)\b\w*\s*acoustic\b`)
	if re.MatchString(title) {
		title = re.ReplaceAllString(title, "")
		return cleanSpaces(title) + " (acoustic)"
	}
	return title
}

func buildKey(artist, title string) string {
	artist = strings.TrimSpace(artist)
	title = strings.TrimSpace(title)

	return artist + "|" + title
}

func NormalizeTrack(artist, title string) NormalizedTrack {
	artists, title := extractFeat(artist, title)

	artist = normalizeArtist(artists)
	title = normalizeTitle(title)

	key := buildKey(artist, title)

	return NormalizedTrack{
		Artist:        artist,
		Title:         title,
		NormalizedKey: key,
	}
}

func unique(input []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))

	for _, v := range input {
		if v == "" {
			continue
		}

		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func toUpper(r rune) rune {
	return []rune(strings.ToUpper(string(r)))[0]
}

func Capitalize(s string) string {
	runes := []rune(strings.ToLower(s))

	capNext := true

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case capNext && isLetter(r):
			runes[i] = toUpper(r)
			capNext = false

		case r == ' ':
			capNext = true

		case r == '(' || r == '[':
			capNext = true

		default:
			capNext = false
		}
	}

	return string(runes)
}
