package utilits

import (
	"strings"
	"testing"
)

func TestBuildListenLinks_ReturnsAppleMusicAndYandexLinks(t *testing.T) {
	links := BuildListenLinks("Linkin Park", "Numb")

	if links.AppleMusic == "" {
		t.Errorf("expected Apple Music link")
	}

	if links.YandexMusic == "" {
		t.Errorf("expected Yandex Music link")
	}

	if !strings.Contains(links.AppleMusic, "music.apple.com") {
		t.Errorf("expected Apple Music domain, got %s", links.AppleMusic)
	}

	if !strings.Contains(links.YandexMusic, "music.yandex.ru") {
		t.Errorf("expected Yandex Music domain, got %s", links.YandexMusic)
	}
}

func TestBuildListenLinks_EncodesQuery(t *testing.T) {
	links := BuildListenLinks("Linkin Park", "Numb")

	if !strings.Contains(links.AppleMusic, "Linkin+Park+Numb") {
		t.Errorf("expected encoded query in Apple Music link, got %s", links.AppleMusic)
	}

	if !strings.Contains(links.YandexMusic, "Linkin+Park+Numb") {
		t.Errorf("expected encoded query in Yandex Music link, got %s", links.YandexMusic)
	}
}

func TestBuildListenLinks_EncodesSpecialCharacters(t *testing.T) {
	links := BuildListenLinks("AC/DC", "Back & Black")

	if !strings.Contains(links.AppleMusic, "AC%2FDC+Back+%26+Black") {
		t.Errorf("expected encoded special characters in Apple Music link, got %s", links.AppleMusic)
	}

	if !strings.Contains(links.YandexMusic, "AC%2FDC+Back+%26+Black") {
		t.Errorf("expected encoded special characters in Yandex Music link, got %s", links.YandexMusic)
	}
}
