package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	triviaUrl = "https://the-trivia-api.com/v2/questions/"
)

type Question struct {
	Text string `json:"text"`
}

type Response struct {
	CorrectAnswer    string   `json:"correctAnswer"`
	IncorrectAnswers []string `json:"incorrectAnswers"`
	Question         Question
}

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
