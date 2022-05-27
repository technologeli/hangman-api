package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type game struct {
	ID              string   `json:"id"`
	Lives           int      `json:"lives"`
	Answer          string   `json:"answer"`
	LowercaseAnswer string   `json:"lowercase_answer"`
	LetterGuesses   string   `json:"letter_guesses"`
	PhraseGuesses   []string `json:"phrase_guesses"`
	Current         string   `json:"current"`
}

var games = []game{}
var alphabet = "abcdefghijklmnopqrstuvwxyz"

func getID() string {
	return fmt.Sprint(1 + len(games))
}

func toUnderscores(answer string) string {
	u := ""
	for _, c := range answer {
		cLower := strings.ToLower(string(c))
		if strings.Contains(alphabet, cLower) {
			u += "_"
		} else {
			u += string(c)
		}
	}
	return u
}

func getCurrent(answer string, guesses string) string {
	curr := []rune(toUnderscores(answer))
	for i, l := range answer {
		lLower := strings.ToLower(string(l))
		if strings.Contains(guesses, lLower) {
			curr[i] = l
		}
	}
	return string(curr)
}

func createGame(answer string) *game {
	g := game{
		ID:              getID(),
		Lives:           5,
		LetterGuesses:   "",
		PhraseGuesses:   []string{},
		Answer:          answer,
		LowercaseAnswer: strings.ToLower(answer),
		Current:         toUnderscores(answer),
	}
	games = append(games, g)
	return &g
}

func getGameByID(id string) (*game, error) {
	for i, g := range games {
		if g.ID == id {
			return &games[i], nil
		}
	}

	return nil, errors.New("Game with ID of " + id + " not found")
}

func guessLetter(id string, guess string) (*game, error) {
	g, err := getGameByID(id)

	if err != nil {
		return nil, err
	}

	// if already guessed letter
	if strings.Contains(g.LetterGuesses, guess) {
		return g, errors.New("Already guessed " + guess)
	}

	g.LetterGuesses += guess

	// if answer contains guess, update guess
	if strings.Contains(g.LowercaseAnswer, guess) {
		g.Current = getCurrent(g.Answer, g.LetterGuesses)
	} else {
		// otherwise decrement lives
		g.Lives--
	}

	return g, nil
}

func guessPhrase(id string, guess string) (*game, error) {
	g, err := getGameByID(id)

	if err != nil {
		return nil, err
	}

	// if already guessed phrase
	for _, phrase := range g.PhraseGuesses {
		if phrase == guess {
			return g, errors.New("Already guessed" + guess)
		}
	}

	g.PhraseGuesses = append(g.PhraseGuesses, guess)

	if g.LowercaseAnswer == guess {
		// if answer == guess
		g.Current = g.Answer
	} else {
		// otherwise decrement lives
		g.Lives--
	}

	return g, nil
}

func makeGuess(id string, guess string) (*game, error) {
	var g *game
	var err error

	guess = strings.ToLower(guess)

	// if guess is a letter
	if len(guess) == 1 {
		g, err = guessLetter(id, guess)
	} else {
		// otherwise phrase guess
		g, err = guessPhrase(id, guess)
	}
	return g, err
}

func getGameStatus(id string) (string, error) {
	g, err := getGameByID(id)

	if err != nil {
		return "", err
	}

	if g.Lives <= 0 {
		return "LOSS", nil
	}

	if g.Answer == g.Current {
		return "WIN", nil
	}

	return "UNFINISHED", nil
}

func playText() {
	var err error
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Input your phrase: ")
	ans, readErr := in.ReadString('\n')
	ans = strings.Replace(ans, "\n", "", -1)
	if readErr != nil {
		fmt.Println(readErr)
	}

	g := createGame(ans)
	for {
		fmt.Println()
		fmt.Println(g.Current)
		fmt.Println("Lives: ", g.Lives)
		fmt.Println("Guessed letters: ", g.LetterGuesses)
		fmt.Println("Guessed phrases: ", g.PhraseGuesses)
		fmt.Print("Make a guess: ")

		line, readErr := in.ReadString('\n')
		if readErr != nil {
			fmt.Println(readErr)
		}
		line = strings.Replace(line, "\n", "", -1)

		g, err = makeGuess(g.ID, line)
		if err != nil {
			fmt.Println(err)
		}

		var status string
		status, err = getGameStatus(g.ID)
		switch status {
		case "WIN":
			fmt.Println()
			fmt.Println("You win! The correct phrase was:")
			fmt.Println(g.Answer)
			return
		case "LOSS":
			fmt.Println()
			fmt.Println("You lost. The correct phrase was:")
			fmt.Println(g.Answer)
			return
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	playText()
}
