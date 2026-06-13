package utilits

import "testing"

func TestNormalizeTrack_BasicLowercaseAndTrim(t *testing.T) {
	result := NormalizeTrack("  ARTIST  ", "  SONG  ")

	if result.Artist != "artist" {
		t.Errorf("expected artist %q, got %q", "artist", result.Artist)
	}

	if result.Title != "song" {
		t.Errorf("expected title %q, got %q", "song", result.Title)
	}

	if result.NormalizedKey != "artist|song" {
		t.Errorf("expected key %q, got %q", "artist|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_ExtractsFeatFromTitle(t *testing.T) {
	result := NormalizeTrack("Artist", "Song feat. Guest")

	if result.Artist != "artist, guest" {
		t.Errorf("expected artist %q, got %q", "artist, guest", result.Artist)
	}

	if result.Title != "song" {
		t.Errorf("expected title %q, got %q", "song", result.Title)
	}

	if result.NormalizedKey != "artist, guest|song" {
		t.Errorf("expected key %q, got %q", "artist, guest|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_ExtractsFeatFromArtist(t *testing.T) {
	result := NormalizeTrack("Artist feat. Guest", "Song")

	if result.Artist != "artist, guest" {
		t.Errorf("expected artist %q, got %q", "artist, guest", result.Artist)
	}

	if result.Title != "song" {
		t.Errorf("expected title %q, got %q", "song", result.Title)
	}

	if result.NormalizedKey != "artist, guest|song" {
		t.Errorf("expected key %q, got %q", "artist, guest|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_SortsAndDeduplicatesArtists(t *testing.T) {
	result := NormalizeTrack("Guest & Artist and Guest", "Song")

	if result.Artist != "artist, guest" {
		t.Errorf("expected artist %q, got %q", "artist, guest", result.Artist)
	}

	if result.NormalizedKey != "artist, guest|song" {
		t.Errorf("expected key %q, got %q", "artist, guest|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_RemovesNonImportantBracketsFromTitle(t *testing.T) {
	result := NormalizeTrack("Artist", "Song (Official Video)")

	if result.Title != "song" {
		t.Errorf("expected title %q, got %q", "song", result.Title)
	}

	if result.NormalizedKey != "artist|song" {
		t.Errorf("expected key %q, got %q", "artist|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_KeepsRemixBracket(t *testing.T) {
	result := NormalizeTrack("Artist", "Song (Cool Remix)")

	if result.Title != "song (remix)" {
		t.Errorf("expected title %q, got %q", "song (remix)", result.Title)
	}

	if result.NormalizedKey != "artist|song (remix)" {
		t.Errorf("expected key %q, got %q", "artist|song (remix)", result.NormalizedKey)
	}
}

func TestNormalizeTrack_KeepsAcousticBracket(t *testing.T) {
	result := NormalizeTrack("Artist", "Song (Live Acoustic Version)")

	if result.Title != "song (acoustic)" {
		t.Errorf("expected title %q, got %q", "song (acoustic)", result.Title)
	}

	if result.NormalizedKey != "artist|song (acoustic)" {
		t.Errorf("expected key %q, got %q", "artist|song (acoustic)", result.NormalizedKey)
	}
}

func TestNormalizeTrack_RemovesNoiseWords(t *testing.T) {
	result := NormalizeTrack("Artist", "Song Official Video HD")

	if result.Title != "song" {
		t.Errorf("expected title %q, got %q", "song", result.Title)
	}

	if result.NormalizedKey != "artist|song" {
		t.Errorf("expected key %q, got %q", "artist|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_KeepSongNameBeforeRemix(t *testing.T) {
	result := NormalizeTrack("Artist", "Summer Remix")

	if result.Title != "summer (remix)" {
		t.Errorf("expected title %q, got %q", "summer (remix)", result.Title)
	}
}
func TestNormalizeTrack_NormalizesAcousticOutsideBrackets(t *testing.T) {
	result := NormalizeTrack("Artist", "Song Acoustic")

	if result.Title != "song (acoustic)" {
		t.Errorf("expected title %q, got %q", "song (acoustic)", result.Title)
	}

	if result.NormalizedKey != "artist|song (acoustic)" {
		t.Errorf("expected key %q, got %q", "artist|song (acoustic)", result.NormalizedKey)
	}
}

func TestNormalizeTrack_SortsArtistsAlphabetically(t *testing.T) {
	result := NormalizeTrack("Zebra & Alpha", "Song")

	if result.Artist != "alpha, zebra" {
		t.Errorf("expected artist %q, got %q", "alpha, zebra", result.Artist)
	}

	if result.NormalizedKey != "alpha, zebra|song" {
		t.Errorf("expected key %q, got %q", "alpha, zebra|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_SupportsXSeparator(t *testing.T) {
	result := NormalizeTrack("Artist x Guest", "Song")

	if result.Artist != "artist, guest" {
		t.Errorf("expected artist %q, got %q", "artist, guest", result.Artist)
	}

	if result.NormalizedKey != "artist, guest|song" {
		t.Errorf("expected key %q, got %q", "artist, guest|song", result.NormalizedKey)
	}
}

func TestNormalizeTrack_SupportsSemicolonSeparator(t *testing.T) {
	result := NormalizeTrack("Artist; Guest", "Song")

	if result.Artist != "artist, guest" {
		t.Errorf("expected artist %q, got %q", "artist, guest", result.Artist)
	}

	if result.NormalizedKey != "artist, guest|song" {
		t.Errorf("expected key %q, got %q", "artist, guest|song", result.NormalizedKey)
	}
}
func TestCapitalize_CapitalizesWords(t *testing.T) {
	result := Capitalize("artist name")

	if result != "Artist Name" {
		t.Errorf("expected %q, got %q", "Artist Name", result)
	}
}

func TestCapitalize_CapitalizesWordsAfterBrackets(t *testing.T) {
	result := Capitalize("song (acoustic version)")

	if result != "Song (Acoustic Version)" {
		t.Errorf("expected %q, got %q", "Song (Acoustic Version)", result)
	}
}

func TestCapitalize_LowercasesOtherLetters(t *testing.T) {
	result := Capitalize("aRTIST nAME")

	if result != "Artist Name" {
		t.Errorf("expected %q, got %q", "Artist Name", result)
	}
}
