package joysound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const JoysoundURL = "https://www.joysound.com/web/search/song"

type ListInfo struct {
	Items []Item `json:"items"`
}

func (a *ListInfo) Append(b *ListInfo) {
	a.Items = append(a.Items, b.Items...)
}

type Item struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Href  string `json:"href"`
}

type Pagination struct {
	Page         int    `json:"page"`
	TotalPages   int    `json:"totalPages"`
	ItemsPerPage int    `json:"itemsPerPage"`
	HashAnchor   string `json:"hashAnchors"`
}

func DecodeListInfo(s string) (*ListInfo, error) {
	const marker = `\"listInfo\":{`

	start := strings.Index(s, marker)
	if start == -1 {
		return nil, fmt.Errorf("listInfo not found")
	}

	// marker 内の `{` の位置
	open := start + strings.Index(marker, "{")

	// JSON文字列内のエスケープも考慮して、対応する `}` を探す
	depth := 0
	inString := false
	// escaped := false
	close := -1

	for i := open; i < len(s); i++ {
		c := s[i]

		if inString {
			// if escaped {
			// 	escaped = false
			// 	continue
			// }
			// if c == '\\' {
			// 	escaped = true
			// 	continue
			// }
			if c == '"' {
				inString = false
			}
			continue
		}

		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
		}
		if c == '}' && depth == 0 {
			close = i
			break
		}
	}

	if close == -1 {
		return nil, fmt.Errorf("matching } not found")
	}

	raw := s[open : close+1]

	// バックスラッシュでエスケープされたJSONを逆エスケープ。
	// 例: {\"foo\":\"bar\"} -> {"foo":"bar"}
	unescaped, err := strconv.Unquote(`"` + raw + `"`)
	if err != nil {
		return nil, fmt.Errorf("unescape listInfo: %w", err)
	}

	var result ListInfo
	if err := json.NewDecoder(strings.NewReader(unescaped)).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode listInfo: %w", err)
	}

	return &result, nil
}

func DecodePagination(s string) (*Pagination, error) {
	const marker = `\"pagination\":{`

	start := strings.Index(s, marker)
	if start == -1 {
		return nil, fmt.Errorf("pagination not found")
	}

	// marker 内の `{` の位置
	open := start + strings.Index(marker, "{")

	// JSON文字列内のエスケープも考慮して、対応する `}` を探す
	depth := 0
	inString := false
	// escaped := false
	close := -1

	for i := open; i < len(s); i++ {
		c := s[i]

		if inString {
			// if escaped {
			// 	escaped = false
			// 	continue
			// }
			// if c == '\\' {
			// 	escaped = true
			// 	continue
			// }
			if c == '"' {
				inString = false
			}
			continue
		}

		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
		}
		if c == '}' && depth == 0 {
			close = i
			break
		}
	}

	if close == -1 {
		return nil, fmt.Errorf("matching } not found")
	}

	raw := s[open : close+1]

	// バックスラッシュでエスケープされたJSONを逆エスケープ。
	// 例: {\"foo\":\"bar\"} -> {"foo":"bar"}
	unescaped, err := strconv.Unquote(`"` + raw + `"`)
	if err != nil {
		return nil, fmt.Errorf("unescape pagination: %w", err)
	}

	var result Pagination
	if err := json.NewDecoder(strings.NewReader(unescaped)).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode Pagination: %w", err)
	}

	return &result, nil
}

func fetch(keyword string, page int) (*goquery.Selection, error) {
	u, err := url.Parse(JoysoundURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("keyword", keyword)
	query.Set("match", "1")
	query.Set("page", strconv.Itoa(page))
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := doc.Find(
		`script`,
	).FilterFunction(func(i int, s *goquery.Selection) bool {
		code := strings.TrimSpace(s.Text())
		return strings.Contains(code, `listInfo`)
	})

	l := results.Length()
	if l != 1 {
		return nil, fmt.Errorf("Expected exactly 1 script tag matches, but %d matched.", l)
	}
	return results, nil
	// pg, err := DecodePagination(results.Text())
	// fmt.Println(pg)
	// return DecodeListInfo(results.Text())
}

func Prefetch(keyword string) (*Pagination, error) {
	resp, err := fetch(keyword, 1)
	if err != nil {
		return nil, err
	}
	return DecodePagination(resp.Text())

}

func FetchAll(keyword string, pg *Pagination) (*ListInfo, error) {
	result := &ListInfo{}
	for p := 1; p <= pg.TotalPages; p++ {
		resp, err := fetch(keyword, p+1)
		if err != nil {
			// return nil, err
			fmt.Fprintln(os.Stderr, err)

			continue
		}
		li, err := DecodeListInfo(resp.Text())
		if err != nil {
			// return nil, err
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		result.Append(li)
	}
	return result, nil
}

func Search(keyword string) (*ListInfo, error) {
	pagination, err := Prefetch(keyword)
	if err != nil {
		return nil, err
	}
	return FetchAll(keyword, pagination)
}

// TODO: make Prefetch and FetchAll private
