package main

import (
	"araok/dam"
	"araok/joysound"
	"fmt"
	"sort"

	"araok/matcher"
)

func main() {
	songKeyword := "NEVER"
	artistKeyword := "こめだわら"

	listInfo, err := joysound.Search(songKeyword)
	if err != nil {
		panic(err)
	}

	songs, err := dam.Search(songKeyword)
	if err != nil {
		panic(err)
	}

	candidates := make([]matcher.Song, 0, len(listInfo.Items)+len(songs))

	// JOYSOUND → matcher.Song
	for _, song := range listInfo.Items {
		candidates = append(candidates, matcher.Song{
			Title:   song.Title,
			Artist:  song.Date,
			Service: matcher.JOYSOUND,
		})
	}

	// DAM → matcher.Song
	for _, song := range songs {
		candidates = append(candidates, matcher.Song{
			Title:   song.Title,
			Artist:  song.Artist,
			Service: matcher.DAM,
		})
	}

	m := matcher.New()

	// まずアーティストで大幅に絞る。
	candidates = m.FilterArtists(candidates, artistKeyword)

	// その後、曲名を評価。
	results := make([]matcher.Result, 0, len(candidates))

	for _, song := range candidates {
		result := m.Match(song, matcher.Query{
			Title:  songKeyword,
			Artist: artistKeyword,
		})

		results = append(results, result)
	}

	// 曲名一致度の高い順。
	sort.Slice(results, func(i, j int) bool {
		return results[i].TitleScore > results[j].TitleScore
	})

	for _, result := range results {
		fmt.Printf(
			"%.3f %-20s %s\n",
			result.TitleScore,
			result.Song.Title,
			result.Song.Artist,
		)
	}
}
