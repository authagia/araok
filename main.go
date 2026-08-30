package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"araok/dam"
	"araok/joysound"
	"araok/matcher"
	"araok/tui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <title>\n", os.Args[0])
		os.Exit(1)
	}

	title := os.Args[1]

	// ここは既存のJOYSOUND/DAM検索処理に置き換える。
	songs, err := searchSongs(title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	model := tui.New(title, songs)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func searchSongs(title string) (matcher.Songs, error) {
	listInfo, err := joysound.Search(title)
	if err != nil {
		return nil, err
	}

	songs, err := dam.Search(title)
	if err != nil {
		return nil, err
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

	return candidates, nil
}
