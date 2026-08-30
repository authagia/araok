package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"araok/matcher"
)

type Model struct {
	title string
	all   matcher.Songs

	input textinput.Model

	marqueePos int
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
	return tea.Batch(
		textinput.Blink,
		marqueeTick(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case marqueeTickMsg:
		m.marqueePos++
		return m, marqueeTick()
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Bold(true)

	artistStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))

	joysoundStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Background(lipgloss.Color("#3A2030")).
			Bold(true).
			Padding(0, 1)

	damStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Background(lipgloss.Color("#18343D")).
			Bold(true).
			Padding(0, 1)
)

const (
	titleWidth   = 34
	artistWidth  = 28
	serviceWidth = 12
)

func (m Model) View() string {

	filtered := matcher.FilterByArtist(m.all, m.input.Value())

	var b strings.Builder

	fmt.Fprintf(&b, "─ %s ─\n\n", m.title)
	b.WriteString(m.input.View())

	b.WriteString("\n\n")

	var countStyle lipgloss.Style

	if len(filtered) == 0 {
		countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)
	} else {
		countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true)
	}
	fmt.Fprintf(
		&b,
		"%s / %d results\n\n",
		countStyle.Render(strconv.Itoa(len(filtered))),
		len(m.all),
	)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Bold(true)

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#44475A"))
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		strings.Repeat(" ", len("* ")),
		lipgloss.NewStyle().
			Width(titleWidth).
			Render(headerStyle.Render("TITLE")),

		lipgloss.NewStyle().
			Width(artistWidth).
			Render(headerStyle.Render("ARTIST")),

		lipgloss.NewStyle().
			Width(serviceWidth).
			Render(headerStyle.Render("SOURCE")),
	)

	b.WriteString("  ")
	b.WriteString(header)
	b.WriteString("\n")

	// b.WriteString("  ")
	b.WriteString(separatorStyle.Render(
		strings.Repeat("─", 2+titleWidth+artistWidth+serviceWidth),
	))
	b.WriteString("\n")

	for _, song := range filtered {
		b.WriteString("  ")
		b.WriteString(m.renderSong(song))
		b.WriteString("\n")
	}

	b.WriteString("\nCtrl+C: quit\n")

	return b.String()
}

func renderService(service matcher.Service) string {
	switch service {
	case matcher.JOYSOUND:
		return joysoundStyle.Render("JOY")
	case matcher.DAM:
		return damStyle.Render("DAM")
	default:
		return ""
	}
}

func (m Model) renderSong(song matcher.Song) string {
	title := marquee(
		song.Title,
		titleWidth,
		m.marqueePos,
	)

	artist := marquee(
		song.Artist,
		artistWidth,
		m.marqueePos,
	)

	title = lipgloss.NewStyle().
		Width(titleWidth).
		Render(titleStyle.Render(title))

	artist = lipgloss.NewStyle().
		Width(artistWidth).
		Render(artistStyle.Render(artist))

	service := lipgloss.NewStyle().
		Width(serviceWidth).
		Render(renderService(song.Service))

	return "* " + lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		artist,
		service,
	)
}

func marquee(s string, width, offset int) string {
	if lipgloss.Width(s) <= width {
		return s
	}

	// 末尾まで行ったら少し間を空けて先頭に戻る。
	text := s + "   "

	runes := []rune(text)

	if len(runes) == 0 {
		return ""
	}

	offset %= len(runes)

	var result []rune
	resultWidth := 0

	for i := 0; resultWidth < width; i++ {
		r := runes[(offset+i)%len(runes)]
		w := lipgloss.Width(string(r))

		if resultWidth+w > width {
			break
		}

		result = append(result, r)
		resultWidth += w
	}

	if resultWidth < width {
		result = append(
			result,
			[]rune(strings.Repeat(" ", width-resultWidth))...,
		)
	}

	return string(result)
}

const marqueeInterval = 150 * time.Millisecond

type marqueeTickMsg struct{}

func marqueeTick() tea.Cmd {
	return tea.Tick(marqueeInterval, func(time.Time) tea.Msg {
		return marqueeTickMsg{}
	})
}
