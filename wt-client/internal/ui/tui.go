package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/internal/types"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type tickMsg time.Time

type tuiLogsView struct {
	vp       viewport.Model
	renderer *glamour.TermRenderer
}

type tuiInputView struct {
	vp       viewport.Model
	renderer *glamour.TermRenderer
	data     types.WtInput
}

type tuiModel struct {
	log       []string
	logFile   *os.File
	logsView  tuiLogsView
	inputView tuiInputView
}

func NewTui() *tuiModel {
	tm := &tuiModel{}
	var err error
	tm.logFile, err = os.OpenFile("tui.log", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(tm, tea.WithAltScreen())
	// tea.WithMouseAllMotion()
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	return tm
}

func (m *tuiModel) Init() tea.Cmd {
	if term.IsTerminal(0) {
		Debug("in a term")
	} else {
		Debug("not in a term")
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		width = 78
		height = 16
	} else {
		width -= 3
		height -= 10
	}

	Debug("term size: ", width, height)

	vp, _ := newViewport(width, 6)
	renderer, _ := newRenderer(&vp, width)
	m.inputView = tuiInputView{
		vp:       vp,
		renderer: renderer,
	}

	vp, _ = newViewport(width, height-6)
	renderer, _ = newRenderer(&vp, width)
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
		m.inputView.vp.Width = msg.Width - padding
		m.inputView.vp.Height = 6
		m.inputView.renderer.Close()

		m.logsView.vp.Width = msg.Width - padding
		m.logsView.vp.Height = msg.Height - padding - 6
		m.logsView.renderer.Close()

		m.logsView.renderer, _ = newRenderer(&m.logsView.vp, m.logsView.vp.Width)
		m.inputView.renderer, _ = newRenderer(&m.inputView.vp, m.logsView.vp.Width)
		return m, nil

	case tickMsg:
		// str, _ := m.inputView.renderer.Render()
		// @todo optimize
		m.inputView.vp.SetContent(renderInputData(&m.inputView.data))

		str, _ := m.logsView.renderer.Render(strings.Join(m.log, "\r\n\r\n"))
		m.logsView.vp.SetContent(str)

		return m, tickCmd()
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *tuiModel) View() string {
	// pad := strings.Repeat(" ", padding)

	s := ""

	s += m.inputView.vp.View() + "\n"
	s += m.logsView.vp.View() + m.helpView()

	return s
}

func (m *tuiModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func newViewport(width int, height int) (viewport.Model, error) {
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return vp, nil
}

func newRenderer(vp *viewport.Model, width int, str ...string) (*glamour.TermRenderer, error) {
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

func (m *tuiModel) Debug(args ...any) {
	length := len(args) + 1
	buf := make([]string, length)
	buf[0] = time.Now().Format("15:04:05.000")
	for i, val := range args {
		buf[i+1] = fmt.Sprint(val)
	}
	str := strings.Join(buf, " ")
	m.log = append(m.log, str)
	m.logFile.WriteString(str + "\n")
	if len(m.log) > 100 {
		newLog := m.log[1:]
		m.log = newLog
	}
}

func (m *tuiModel) Error(args ...any) {
	newArgs := []any{"ERR:"}
	newArgs = append(newArgs, args...)
	m.Debug(newArgs...)
}

func (m *tuiModel) Input(in *types.WtInput) {
	m.inputView.data = *in
}

func renderInputData(in *types.WtInput) string {
	return fmt.Sprintf(
		"  Game running: %s Map loaded: %s Mode: %s Mission started: %s Mission ended: %s\r\n"+
			"  Player: type %s vehicle: %s landed: %s dead: %s\r\n"+
			"  Enemies: last kill %d, air close: %s battle %s, ground close: %s, battle %s, nearest air %s, nearest ground: %s\r\n"+
			"  Map: not implemented yet",
		btyn(in.GameRunning),
		btyn(in.MapLoaded),
		in.GameMode,
		btyn(in.MissionStarted),
		btyn(in.MissionEnded),
		in.PlayerType,
		in.PlayerVehicle,
		btyn(in.PlayerLanded),
		btyn(in.PlayerDead),
		in.LastPlayerMadeKillTime,
		btyn(in.EnemyAirNear),
		btyn(in.EnemyAirClose),
		btyn(in.EnemyGroundNear),
		btyn(in.EnemyGroundClose),
		fts(in.NearestEnemyAir),
		fts(in.NearestEnemyGround),
	)
}

func btyn(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func fts(dist float64) string {
	if dist == -1.0 {
		return "none"
	} else {
		return fmt.Sprintf("%f km", dist/1000)
	}
}
