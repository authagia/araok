package matcher

import (
	"math"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Service int

const (
	JOYSOUND Service = iota
	DAM
)

type Song struct {
	Title   string
	Artist  string
	Service Service
}

type Query struct {
	Title  string
	Artist string
}

type Result struct {
	Song Song

	// ArtistMatched はアーティスト候補として採用されたか。
	ArtistMatched bool

	// TitleScore は曲名の一致度。0.0〜1.0。
	TitleScore float64

	// TitleMatch は曲名の一致方法。
	TitleMatch MatchType
}

type MatchType int

const (
	MatchNone MatchType = iota
	MatchExact
	MatchPrefix
	MatchContains
	MatchEditDistance
)

type Matcher struct {
	titlePrefixes []string
}

func New() *Matcher {
	return &Matcher{
		titlePrefixes: []string{
			"[生音]",
			"[良音]",
			"[テクノ]",
			"[プロオケ]",
		},
	}
}

// Match は Song が Query の候補として妥当かを判定する。
//
// Artist:
//   - NFKC
//   - ひらがな → カタカナ
//   - 大文字小文字を無視
//   - Contains
//
// Title:
//   - NFKC
//   - 既知の prefix を除去
//   - Exact / Prefix / Contains / 編集距離
func (m *Matcher) Match(song Song, query Query) Result {
	artistMatched := matchArtist(query.Artist, song.Artist)

	titleScore, titleMatch := matchTitle(
		query.Title,
		song.Title,
		m.titlePrefixes,
	)

	return Result{
		Song:          song,
		ArtistMatched: artistMatched,
		TitleScore:    titleScore,
		TitleMatch:    titleMatch,
	}
}

// FilterArtists はアーティスト条件に一致する候補だけを返す。
func (m *Matcher) FilterArtists(songs []Song, artist string) []Song {
	if artist == "" {
		return songs
	}

	result := make([]Song, 0, len(songs))

	for _, song := range songs {
		if matchArtist(artist, song.Artist) {
			result = append(result, song)
		}
	}

	return result
}

func matchArtist(query, candidate string) bool {
	if query == "" {
		return true
	}

	q := normalizeArtist(query)
	c := normalizeArtist(candidate)

	return strings.Contains(c, q)
}

func normalizeArtist(s string) string {
	s = norm.NFKC.String(s)
	s = hiraganaToKatakana(s)
	return strings.ToLower(s)
}

func normalizeTitle(s string, prefixes []string) string {
	s = norm.NFKC.String(s)

	for _, prefix := range prefixes {
		prefix = norm.NFKC.String(prefix)

		if strings.HasPrefix(s, prefix) {
			s = strings.TrimSpace(strings.TrimPrefix(s, prefix))
			break
		}
	}

	return strings.TrimSpace(s)
}

func matchTitle(query, candidate string, prefixes []string) (float64, MatchType) {
	if query == "" {
		return 0, MatchNone
	}

	q := normalizeTitle(query, nil)
	c := normalizeTitle(candidate, prefixes)

	if q == "" || c == "" {
		return 0, MatchNone
	}

	// 完全一致
	if c == q {
		return 1.0, MatchExact
	}

	// 前方一致
	if strings.HasPrefix(c, q) {
		return 0.9, MatchPrefix
	}

	// 部分一致
	if strings.Contains(c, q) {
		return 0.75, MatchContains
	}

	// 編集距離
	score := editDistanceScore(q, c)

	return score, MatchEditDistance
}

func editDistanceScore(a, b string) float64 {
	distance := levenshtein([]rune(a), []rune(b))

	maxLen := max(len([]rune(a)), len([]rune(b)))
	if maxLen == 0 {
		return 1.0
	}

	return math.Max(0, 1.0-float64(distance)/float64(maxLen))
}

func levenshtein(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}

	if len(b) == 0 {
		return len(a)
	}

	// b のほうを短くしてメモリ使用量を抑える。
	if len(a) < len(b) {
		a, b = b, a
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)

	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i

		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			curr[j] = min(
				curr[j-1]+1,    // insertion
				prev[j]+1,      // deletion
				prev[j-1]+cost, // substitution
			)
		}

		prev, curr = curr, prev
	}

	return prev[len(b)]
}

func hiraganaToKatakana(s string) string {
	runes := []rune(s)

	for i, r := range runes {
		// ぁ〜ゖ → ァ〜ヶ
		if r >= 'ぁ' && r <= 'ゖ' {
			runes[i] = r + ('ァ' - 'ぁ')
		}
	}

	return string(runes)
}
