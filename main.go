package main

import (
	"araok/dam"
	"araok/joysound"
	"fmt"

	"araok/matcher"
)

func main() {
	songKeyword := "シャルル"
	artistKeyword := "バルーン"

	listInfo, err := joysound.Search(songKeyword)
	if err != nil {
		panic(err)
	}

	songs, err := dam.Search(songKeyword)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hit %d on joysound and %d on dam\n", len(listInfo.Items), len(songs))
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

	if len(candidates) < 20 {
		for _, cand := range candidates {
			fmt.Printf(
				"%-20s %s %s\n",
				cand.Title,
				cand.Artist,
				cand.Service,
			)
		}
		return
	}

	results := matcher.FilterByArtist(
		candidates,
		artistKeyword,
	)

	for _, result := range results {
		fmt.Printf(
			"%-20s %s %s\n",
			result.Title,
			result.Artist,
			result.Service,
		)
	}
}
