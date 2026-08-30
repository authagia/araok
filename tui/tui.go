package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
	var (
		titleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#ccc")).
				Width(64).
				MaxWidth(64).
				Padding(0, 1).
				Bold(true)

		artistStyle = lipgloss.NewStyle().
				Width(32).
				Padding(0, 4).
				Foreground(lipgloss.Color("252"))

		joysoundStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF79C6")).
				Background(lipgloss.Color("#4A1834")).
				Padding(0, 1).
				Bold(true)

		damStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8BE9FD")).
				Background(lipgloss.Color("#163B46")).
				Padding(0, 1).
				Bold(true)

		serviceStyle = lipgloss.NewStyle().
				Width(10).
				Align(lipgloss.Left)
	)

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
		service := song.Service.String()

		switch song.Service {
		case matcher.JOYSOUND:
			service = joysoundStyle.Render(service)
		case matcher.DAM:
			service = damStyle.Render(service)
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			titleStyle.Render(song.Title),

			artistStyle.Render(song.Artist),

			serviceStyle.Render(service),
		)

		b.WriteString("  ")
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\nCtrl+C: quit\n")

	return b.String()
}
