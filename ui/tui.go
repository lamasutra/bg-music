package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type tickMsg time.Time

type tuiMusicProgressbar struct {
	title string
	value float64
	model progress.Model
}

type tuiVolumeProgressbar struct {
	value float64
	model progress.Model
}
type tuiLogsView struct {
	vp       viewport.Model
	renderer *glamour.TermRenderer
}

type tuiModel struct {
	log      string
	music    tuiMusicProgressbar
	volume   tuiVolumeProgressbar
	logsView tuiLogsView
}

func NewTui() *tuiModel {
	tm := &tuiModel{
		music: tuiMusicProgressbar{
			model: progress.New(),
		},
		volume: tuiVolumeProgressbar{
			model: progress.New(),
		},
	}
	p := tea.NewProgram(tm, tea.WithAltScreen(), tea.WithMouseAllMotion())
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	return tm
}

func (m *tuiModel) Debug(args ...any) {
	var newArgs []any
	newArgs = []any{time.Now().Format("15:04:05.000")}
	newArgs = append(newArgs, args...)
	m.log += fmt.Sprintln(newArgs...)
}

func (m *tuiModel) Error(args ...any) {
	newArgs := []any{"ERR:"}
	newArgs = append(newArgs, args...)
	m.Debug(newArgs...)
}

func (m *tuiModel) SetCurrentMusicProgress(progress float64) {
	m.music.value = progress
}

func (m *tuiModel) SetCurrentMusicTitle(title string) {
	m.music.title = title
}

func (m *tuiModel) SetCurrentVolume(volume float64) {
	m.volume.value = volume
}

func (m *tuiModel) Init() tea.Cmd {
	vp, _ := logsViewport(78)
	renderer, _ := logsRenderer(&vp, 78)
	m.logsView = tuiLogsView{
		vp:       vp,
		renderer: renderer,
	}

	return tickCmd()
}

func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		default:
			// fmt.Println(msg.String())
			var cmd tea.Cmd
			m.logsView.vp, _ = m.logsView.vp.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.music.model.Width = msg.Width - padding
		if m.music.model.Width > maxWidth {
			m.music.model.Width = maxWidth
		}
		m.volume.model.Width = msg.Width - padding
		if m.volume.model.Width > maxWidth {
			m.volume.model.Width = maxWidth
		}
		m.logsView.vp.Width = msg.Width - padding
		m.logsView.vp.Height = msg.Height - padding - 8
		m.logsView.renderer.Close()

		m.logsView.renderer, _ = logsRenderer(&m.logsView.vp, m.logsView.vp.Width)
		return m, nil

	case tickMsg:
		str, _ := m.logsView.renderer.Render(m.log)

		m.logsView.vp.SetContent(str)

		return m, tickCmd()
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *tuiModel) View() string {
	pad := strings.Repeat(" ", padding)

	s := ""

	s += pad + m.music.title + "\n"
	s += pad + m.music.model.ViewAs(m.music.value) + "\n\n"

	s += pad + "Volume:\n"
	s += pad + m.volume.model.ViewAs(m.volume.value) + "\n\n"

	s += m.logsView.vp.View() + m.helpView()

	return s
}

func (m *tuiModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func logsViewport(width int) (viewport.Model, error) {
	vp := viewport.New(width, 16)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return vp, nil
}

func logsRenderer(vp *viewport.Model, width int, str ...string) (*glamour.TermRenderer, error) {
	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return nil, err
	}

	return renderer, nil
}
