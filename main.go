package main

import (
	"araok/dam"
	"araok/joysound"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func searchJoysound(keyword string) (*joysound.ListInfo, error) {
	u, err := url.Parse(joysound.JoysoundURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("keyword", keyword)
	query.Set("match", "1")
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
	return joysound.DecodeListInfo(results.Text())
}

func searchDam(keyword string) (*dam.SearchResponse, error) {
	payload := map[string]any{
		"modelTypeCode": "1",
		"serialNo":      "BA000001",
		"keyword":       keyword,
		"compId":        "1",
		"authKey":       "2/Qb9R@8s*", //https://www.clubdam.com/assets/dkcommon/js/karaokesearch.js
		"contentsCode":  nil,
		"serviceCode":   nil,
		"sort":          "2",
		"dispCount":     "100",
		"pageNo":        "1",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		// panic(err)
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		dam.DamAPI,
		bytes.NewReader(body),
	)
	if err != nil {
		// panic(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// panic(err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		// panic(fmt.Sprintf("unexpected HTTP status: %s", resp.Status))
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var result dam.SearchResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	keyword := "妄想感傷代償連盟"

	listInfo, err := searchJoysound(keyword)
	if err != nil {
		panic(err)
	}
	for _, song := range listInfo.Items {
		fmt.Printf(
			"%s - %s [%s]\n",
			song.Title,
			song.Date,
			song.Href,
		)
	}

	// result, err := searchDam(keyword)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("API Status:", result.Result.StatusCode)
	// fmt.Println("Message:", result.Result.Message)
	// fmt.Println("Keyword:", result.Data.Keyword)
	// fmt.Println("Total:", result.Data.TotalCount)

	// for _, song := range result.List {
	// 	fmt.Printf(
	// 		"%s - %s [%s]\n",
	// 		song.Title,
	// 		song.Artist,
	// 		song.RequestNo,
	// 	)
	// }
}
