package utilits

import "net/url"

type ListenLinks struct {
	AppleMusic  string
	YandexMusic string
}

func BuildListenLinks(artist, title string) ListenLinks {
	query := url.QueryEscape(artist + " " + title)

	return ListenLinks{
		AppleMusic:  "https://music.apple.com/us/search?term=" + query,
		YandexMusic: "https://music.yandex.ru/search?text=" + query,
	}
}
