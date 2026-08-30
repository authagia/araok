package joysound

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const JoysoundURL = "https://www.joysound.com/web/search/song"

type ListInfo struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Href  string `json:"href"`
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
