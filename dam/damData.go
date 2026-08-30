package dam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const DamAPI = "https://www.clubdam.com/dkwebsys/search-api/SearchMusicByKeywordApi"

type SearchResponse struct {
	Result Result `json:"result"`
	Data   Data   `json:"data"`
	List   []Song `json:"list"`
}

type Result struct {
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
}

type Data struct {
	Keyword    string `json:"keyword"`
	PageNo     int    `json:"pageNo"`
	PageCount  int    `json:"pageCount"`
	HasNext    string `json:"hasNext"`
	HasPreview string `json:"hasPreview"`
	TotalCount int    `json:"totalCount"`
	OverFlag   string `json:"overFlag"`
}

type Song struct {
	RequestNo         string `json:"requestNo"`
	Title             string `json:"title"`
	TitleYomi         string `json:"titleYomi"`
	ArtistCode        int    `json:"artistCode"`
	Artist            string `json:"artist"`
	ArtistYomi        string `json:"artistYomi"`
	ReleaseDate       string `json:"releaseDate"`
	NewReleaseFlag    string `json:"newReleaseFlag"`
	FutureReleaseFlag string `json:"futureReleaseFlag"`

	HighlightLyrics string `json:"highlightLyrics"`

	HonninFlag     string `json:"honninFlag"`
	AnimeFlag      string `json:"animeFlag"`
	LiveFlag       string `json:"liveFlag"`
	KidsFlag       string `json:"kidsFlag"`
	MamaotoFlag    string `json:"mamaotoFlag"`
	NamaotoFlag    string `json:"namaotoFlag"`
	DuetFlag       string `json:"duetFlag"`
	GuideVocalFlag string `json:"guideVocalFlag"`
	ProokeFlag     string `json:"prookeFlag"`
	DuetDxFlag     string `json:"duetDxFlag"`

	DamTomoMovieFlag           string `json:"damTomoMovieFlag"`
	DamTomoRecordingFlag       string `json:"damTomoRecordingFlag"`
	DamTomoPublicVocalFlag     string `json:"damTomoPublicVocalFlag"`
	DamTomoPublicMovieFlag     string `json:"damTomoPublicMovieFlag"`
	DamTomoPublicRecordingFlag string `json:"damTomoPublicRecordingFlag"`

	ScoreFlag    string `json:"scoreFlag"`
	MyListFlag   string `json:"myListFlag"`
	Shift        string `json:"shift"`
	PlaybackTime int    `json:"playbackTime"`
}

func fetch(keyword string, page int) (*SearchResponse, error) {
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
		"pageNo":        strconv.Itoa(page),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		// panic(err)
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		DamAPI,
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

	// fmt.Println("HTTP Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		// panic(fmt.Sprintf("unexpected HTTP status: %s", resp.Status))
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var result SearchResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Prefetch(keyword string) (*SearchResponse, error) {
	return fetch(keyword, 1)
}

func FetchAll(prefetchResult *SearchResponse) ([]Song, error) {
	keyword := prefetchResult.Data.Keyword
	songs := prefetchResult.List
	for p := 2; p <= prefetchResult.Data.PageCount; p++ {
		resp, err := fetch(keyword, p+1)
		if err != nil {
			// return nil, err
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		songs = append(songs, resp.List...)
	}
	return songs, nil
}

func Search(keyword string) ([]Song, error) {
	prefetchResult, err := Prefetch(keyword)
	if err != nil {
		return nil, err
	}
	return FetchAll(prefetchResult)
}

// TODO: make Prefetch and FetchAll private
