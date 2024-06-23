package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	hotPink = lipgloss.Color("#FF06B7")
)

type GameState int

const (
	newGame = iota
	playing
	gameOver
)

var (
	textStyle   = lipgloss.NewStyle().Foreground(hotPink)
	centerStyle = lipgloss.NewStyle().Align(lipgloss.Center).Width(50)
	stateName   = map[GameState]string{
		newGame:  "new game",
		playing:  "playing",
		gameOver: "game over",
	}
)

type model struct {
	input     textinput.Model
	state     GameState
	points    int
	cursor    int
	err       error
	questions []string
	answers   []string
	selected  map[int]struct{}
}

type (
	errMsg error
)

func (gs GameState) String() string {
	return stateName[gs]
}

func initialModel() model {
	var input = textinput.New()
	input.Focus()
	input.CharLimit = 30
	input.Width = 30
	input.Placeholder = "y/n ?"

	return model{
		input:  input,
		state:  newGame,
		points: 0,
	}
}

func newGameView(m model) string {
	return centerStyle.Render(fmt.Sprintf(
		`
		%s
		%s
		`,
		textStyle.Render("Ready to start?"),
		m.input.View(),
	))
}

func playingView() string {
	return centerStyle.Render(
		fmt.Sprintf(
			`
			%s
			`,
			textStyle.Render("Playing.."),
		),
	)
}

func errorView() string {
	return centerStyle.Render(
		fmt.Sprintf(
			`
			%s
			`,
			textStyle.Render("Something went wrong"),
		),
	)
}

func gameOverView() string {
	return centerStyle.Render(
		fmt.Sprintf(
			`
			%s
			`,
			textStyle.Render("Press ctrl+c / esc to quit"),
		),
	)
}

func handleInputKeys(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc", "ctrl+c":
		return m, tea.Quit
	case "enter":
		inputVal := m.input.Value()
		if inputVal == "y" || inputVal == "Y" {
			m.state = playing
			m.input.Reset()
		} else if inputVal == "n" || inputVal == "N" {
			m.state = gameOver
			m.input.Reset()
		} else {
			m.state = newGame
			m.input.Reset()
			m.input.Placeholder = "Please enter y or n"
		}

		return m, nil
	case "ctrl+n":
		m.state = newGame
		m.input.Reset()
		return m, nil
	default:
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleInputKeys(msg, m)
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case newGame:
		return newGameView(m)
	case playing:
		return playingView()
	case gameOver:
		return gameOverView()
	default:
		return errorView()
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
