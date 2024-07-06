package main

import "fmt"

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

func playingView(m model) string {
	return m.question.View()

}

func gameOverView(m model) string {
	return centerStyle.Render(
		fmt.Sprintf(
			`
			%s
			%s
			`,
			textStyle.Render(fmt.Sprintf("Thanks for playing! You got %d points!", m.points)),
			textStyle.Render("Press ctrl+c / esc to quit"),
		),
	)
}

func gameErorView(m model) string {
	return centerStyle.Render(
		fmt.Sprintf(
			`
			%s%s
			`,
			textStyle.Render("Something went wrong: "),
			textStyle.Render(m.err.Error()),
		),
	)
}
