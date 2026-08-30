package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"araok/matcher"
)

type Model struct {
	title string
	all   matcher.Songs

	input textinput.Model
}

func New(title string, songs matcher.Songs) Model {
	input := textinput.New()
	input.Prompt = "Artist: "
	input.Focus()

	return Model{
		title: title,
		all:   songs,
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	// filtered := m.all.FilterByArtist(m.input.Value())
	filtered := matcher.FilterByArtist(m.all, m.input.Value())

	var b strings.Builder

	fmt.Fprintf(&b, "─ %s ─\n\n", m.title)

	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	fmt.Fprintf(
		&b,
		"%d / %d results\n\n",
		len(filtered),
		len(m.all),
	)

	for _, song := range filtered {
		fmt.Fprintf(
			&b,
			"  %-30s %-30s %-10s\n",
			song.Title,
			song.Artist,
			song.Service,
		)
	}

	b.WriteString("\nCtrl+C: quit\n")

	return b.String()
}
