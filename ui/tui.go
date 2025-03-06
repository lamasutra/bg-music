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

const content = `
# Today’s Menu

## Appetizers

| Name        | Price | Notes                           |
| ---         | ---   | ---                             |
| Tsukemono   | $2    | Just an appetizer               |
| Tomato Soup | $4    | Made with San Marzano tomatoes  |
| Okonomiyaki | $4    | Takes a few minutes to make     |
| Curry       | $3    | We can add squash if you’d like |

## Seasonal Dishes

| Name                 | Price | Notes              |
| ---                  | ---   | ---                |
| Steamed bitter melon | $2    | Not so bitter      |
| Takoyaki             | $3    | Fun to eat         |
| Winter squash        | $3    | Today it's pumpkin |

## Desserts

| Name         | Price | Notes                 |
| ---          | ---   | ---                   |
| Dorayaki     | $4    | Looks good on rabbits |
| Banana Split | $5    | A classic             |
| Cream Puff   | $3    | Pretty creamy!        |

All our dishes are made in-house by Karen, our chef. Most of our ingredients
are from our garden or the fish market down the street.

Some famous people that have eaten here lately:

* [x] René Redzepi
* [x] David Chang
* [ ] Jiro Ono (maybe some day)

Bon appétit!
`

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
	vp       *viewport.Model
	renderer *glamour.TermRenderer
}

type tuiModel struct {
	log      []string
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
	p := tea.NewProgram(tm, tea.WithAltScreen())
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
	m.log = append(m.log, fmt.Sprintln(newArgs...))
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
	renderer, _ := logsRenderer(vp, 78)
	m.logsView = tuiLogsView{
		vp:       vp,
		renderer: renderer,
	}
	str, _ := renderer.Render(content)
	vp.SetContent(str)
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return time.Time(t)
	})
}

func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		// case "up", "k":
		//     if m.cursor > 0 {
		//         m.cursor--
		//     }

		// The "down" and "j" keys move the cursor down
		// case "down", "j":
		//     if m.cursor < len(m.choices)-1 {
		//         m.cursor++
		//     }

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		// case "enter", " ":
		//     _, ok := m.selected[m.cursor]
		//     if ok {
		//         delete(m.selected, m.cursor)
		//     } else {
		//         m.selected[m.cursor] = struct{}{}
		//     }
		// }
		default:
			var cmd tea.Cmd
			m.logsView.vp.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.music.model.Width = msg.Width - padding*2 - 4
		if m.music.model.Width > maxWidth {
			m.music.model.Width = maxWidth
		}
		m.volume.model.Width = msg.Width - padding*2 - 4
		if m.volume.model.Width > maxWidth {
			m.volume.model.Width = maxWidth
		}
		m.logsView.vp.Width = msg.Width - padding*2 - 4
		m.logsView.renderer.Close()
		m.logsView.renderer, _ = logsRenderer(m.logsView.vp, m.logsView.vp.Width)
		return m, nil

	case tickMsg:
		// m.percent += 0.25
		// if m.percent > 1.0 {
		// 	m.percent = 1.0
		// 	return m, tea.Quit
		// }
		return m, tickCmd()
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *tuiModel) View() string {
	// The header

	pad := strings.Repeat(" ", padding)

	s := ""

	// s := pad + m.player. "What should we buy at the market?\n\n"

	// Iterate over our choices
	// for i, choice := range m.choices {

	//     // Is the cursor pointing at this choice?
	//     cursor := " " // no cursor
	//     if m.cursor == i {
	//         cursor = ">" // cursor!
	//     }

	//     // Is this choice selected?
	//     checked := " " // not selected
	//     if _, ok := m.selected[i]; ok {
	//         checked = "x" // selected!
	//     }

	//     // Render the row
	//     s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	// }

	// The footer
	// s += "\nPress q to quit.\n"

	s += pad + m.music.title + "\n"
	s += pad + m.music.model.ViewAs(m.music.value) + "\n\n"

	s += pad + "Volume:\n"
	s += pad + m.volume.model.ViewAs(m.volume.value) + "\n\n"

	// pad + helpStyle("Press any key to quit")

	var wholeLog string
	for _, l := range m.log {
		wholeLog += fmt.Sprint(l)
	}
	// str, _ := m.logsView.renderer.Render(wholeLog)
	// s += str

	// m.logsView.vp.SetContent(str)
	s += m.logsView.vp.View() + m.helpView()

	// Send the UI for rendering
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

func logsViewport(width int) (*viewport.Model, error) {
	vp := viewport.New(width, 16)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return &vp, nil
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
