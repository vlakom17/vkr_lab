package utilits

import "testing"

func TestValidateTrackInput_ReturnsNilForValidInput(t *testing.T) {
	err := ValidateTrackInput("Linkin Park", "Numb")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateTrackInput_ReturnsErrorWhenArtistIsEmpty(t *testing.T) {
	err := ValidateTrackInput("   ", "Numb")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	err := ValidateTrackInput("Linkin Park", "   ")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenArtistIsTooLong(t *testing.T) {
	longArtist := makeLongString(111)

	err := ValidateTrackInput(longArtist, "Numb")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenTitleIsTooLong(t *testing.T) {
	longTitle := makeLongString(111)

	err := ValidateTrackInput("Linkin Park", longTitle)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenArtistHasNoLetters(t *testing.T) {
	err := ValidateTrackInput("12345", "Numb")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenTitleHasNoLetters(t *testing.T) {
	err := ValidateTrackInput("Linkin Park", "12345")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenInputContainsURL(t *testing.T) {
	err := ValidateTrackInput("Linkin Park", "https://example.com")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenInputContainsHTML(t *testing.T) {
	err := ValidateTrackInput("Linkin Park", "<script>alert(1)</script>")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateTrackInput_ReturnsErrorWhenInputLooksSuspicious(t *testing.T) {
	err := ValidateTrackInput("Linkin Park", "!!!@@@###")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func makeLongString(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		result += "a"
	}

	return result
}
