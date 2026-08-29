package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	damURL         = "https://www.clubdam.com/karaokesearch/"
	resultSelector = "div.result-wrap > ul > li"
	userAgent      = "curl/8.21.0"
)

func main() {
	keyword := "KING"

	// URLを組み立てる
	u, err := url.Parse(damURL)
	if err != nil {
		panic(err)
	}

	query := u.Query()
	query.Set("keyword", keyword)
	query.Set("type", "song")
	u.RawQuery = query.Encode()

	fmt.Println("Request URL:", u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", userAgent)
	fmt.Println(req.Header)

	resp, err := http.DefaultClient.Do(req)
	// HTTPリクエスト
	// resp, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("URL:", u.String())
	fmt.Println("Status:", resp.Status)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println("Content-Length:", resp.Header.Get("Content-Length"))

	fmt.Println("Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	// HTMLを解析
	w := os.Stdout
	teeReader := io.TeeReader(resp.Body, w)
	defer io.ReadAll(teeReader)
	doc, err := goquery.NewDocumentFromReader(teeReader)
	// doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// 検索結果を取得
	results := doc.Find(resultSelector)

	fmt.Println("Result count:", results.Length())

	// 各検索結果を表示
	results.Each(func(i int, s *goquery.Selection) {
		fmt.Printf("\n--- Result %d ---\n", i+1)
		fmt.Println(s.Text())
	})
}
