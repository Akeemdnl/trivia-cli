package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	input           textinput.Model
	state           GameState
	points          int
	streakPoints    int
	streak          bool
	currentQuestion int
	err             error
	question        list.Model
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (gs GameState) String() string {
	return stateName[gs]
}

func initialModel() model {
	var input = textinput.New()
	input.Focus()
	input.CharLimit = 30
	input.Width = 30
	input.Placeholder = "y/n ?"
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 50, 50)
	l.SetShowStatusBar(false)

	return model{
		input:           input,
		state:           newGame,
		points:          0,
		streakPoints:    0,
		streak:          false,
		currentQuestion: 0,
		question:        l,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.question.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		return handleInputKeys(msg, m)
	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.question, cmd = m.question.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case quitting:
		return ""
	case newGame:
		return newGameView(m)
	case playing:
		return playingView(m)
	case gameOver:
		return gameOverView(m)
	case gameError:
		return gameErorView(m)
	default:
		return errorView()
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops something went wrong: %v", err)
		os.Exit(1)
	}
}
