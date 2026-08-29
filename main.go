package main

import (
	"araok/dam"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func searchDam(keyword string) (dam.SearchResponse, error) {
	payload := map[string]any{
		"modelTypeCode": "1",
		"serialNo":      "BA000001",
		"keyword":       keyword,
		"compId":        "1",
		"authKey":       "2/Qb9R@8s*",
		"contentsCode":  nil,
		"serviceCode":   nil,
		"sort":          "2",
		"dispCount":     "100",
		"pageNo":        "1",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		// panic(err)
		return dam.SearchResponse{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		dam.DamAPI,
		bytes.NewReader(body),
	)
	if err != nil {
		// panic(err)
		return dam.SearchResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// panic(err)
		return dam.SearchResponse{}, err
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		// panic(fmt.Sprintf("unexpected HTTP status: %s", resp.Status))
		return dam.SearchResponse{}, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var result dam.SearchResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func main() {
	keyword := "妄想感傷代償連盟"
	result, err := searchDam(keyword)
	if err != nil {
		panic(err)
	}

	fmt.Println("API Status:", result.Result.StatusCode)
	fmt.Println("Message:", result.Result.Message)
	fmt.Println("Keyword:", result.Data.Keyword)
	fmt.Println("Total:", result.Data.TotalCount)

	for _, song := range result.List {
		fmt.Printf(
			"%s - %s [%s]\n",
			song.Title,
			song.Artist,
			song.RequestNo,
		)
	}
}
