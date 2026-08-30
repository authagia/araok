package matcher

import (
	"fmt"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Service int

const (
	JOYSOUND Service = iota
	DAM
)

func (s Service) String() string {
	switch s {
	case JOYSOUND:
		return "Joysound"
	case DAM:
		return "DAM"
	default:
		return fmt.Sprintf("Service(%d)", s)
	}
}

type Song struct {
	Title   string
	Artist  string
	Service Service
}

type Songs []Song

func FilterByArtist(songs Songs, artist string) Songs {
	if artist == "" {
		return songs
	}

	query := normalizeArtist(artist)

	result := make(Songs, 0, len(songs))

	for _, song := range songs {
		candidate := normalizeArtist(song.Artist)

		if strings.Contains(candidate, query) {
			result = append(result, song)
		}
	}

	return result
}

func normalizeArtist(s string) string {
	s = norm.NFKC.String(s)
	s = hiraganaToKatakana(s)
	s = strings.ToLower(s)
	return s
}

func hiraganaToKatakana(s string) string {
	runes := []rune(s)

	for i, r := range runes {
		if r >= 'ぁ' && r <= 'ゖ' {
			runes[i] = r + ('ァ' - 'ぁ')
		}
	}

	return string(runes)
}
