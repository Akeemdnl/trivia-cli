package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	hotPink   = lipgloss.Color("#FF06B7")
	triviaUrl = "https://the-trivia-api.com/v2/questions/"
	newGame   = iota
	playing
	gameOver
	gameError
	quitting
)

var (
	textStyle   = lipgloss.NewStyle().Foreground(hotPink)
	centerStyle = lipgloss.NewStyle().Align(lipgloss.Center).Width(80)
	stateName   = map[GameState]string{
		newGame:   "new game",
		playing:   "playing",
		gameOver:  "game over",
		gameError: "game error",
		quitting:  "quitting",
	}
)

var (
	questions        []string
	correctAnswers   []string
	incorrectAnswers map[int][]string
)

type GameState int
type item string

type Question struct {
	Text string `json:"text"`
}

type Response struct {
	CorrectAnswer    string   `json:"correctAnswer"`
	IncorrectAnswers []string `json:"incorrectAnswers"`
	Question         Question
}

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }

func getQuestions() ([]Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(triviaUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	var resp []Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func handleInputKeys(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.state = quitting
		return m, tea.Quit
	case "enter":
		inputVal := m.input.Value()
		if inputVal == "y" || inputVal == "Y" {
			startGame(&m)
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
		var inputCmd tea.Cmd
		var questionCmd tea.Cmd

		m.input, inputCmd = m.input.Update(msg)
		m.question, questionCmd = m.question.Update(msg)
		cmds := []tea.Cmd{inputCmd, questionCmd}
		return m, tea.Batch(cmds...)
	}
}

func startGame(m *model) {
	m.input.Reset()

	res, err := getQuestions()
	if err != nil {
		m.state = gameError
		m.err = err
		return
	}

	incorrectAnswers = make(map[int][]string)
	m.state = playing
	for i, item := range res {
		questions = append(questions, item.Question.Text)
		correctAnswers = append(correctAnswers, item.CorrectAnswer)
		incorrectAnswers[i] = item.IncorrectAnswers
	}
	m.setQuestion(questions[0], correctAnswers[0], incorrectAnswers[0])
}

func (m *model) setQuestion(question string, correctAnswer string, incorrectAnswers []string) {
	var totalAnswers []string
	m.question.Title = question
	items := m.question.Items()

	if items == nil {
		items = []list.Item{}
	}

	fmt.Println(reflect.TypeOf(items))
	totalAnswers = incorrectAnswers
	totalAnswers = append(totalAnswers, correctAnswer)

	for _, v := range totalAnswers {
		items = append(items, item(v))
	}

	m.question.SetItems(items)
}
