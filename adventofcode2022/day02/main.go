package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type RPSMove int
type RPSOutcome int

const (
	Rock RPSMove = iota + 1
	Paper
	Scissors

	Loss RPSOutcome = 0
	Draw RPSOutcome = 3
	Win  RPSOutcome = 6
)

type (
	Round struct {
		Opponent       RPSMove
		Ours           RPSMove
		DesiredOutcome RPSOutcome
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)

	rounds := getRounds(lines)

	score := 0

	for _, r := range rounds {
		// Set the move for part 2
		r.SetOutMove()
		score += r.Score()
	}

	fmt.Printf("Rounds: %v score: %v\n", len(rounds), score)
}

func getMove(s string) RPSMove {
	switch strings.ToUpper(s) {
	case "A", "X":
		return Rock
	case "B", "Y":
		return Paper
	case "C", "Z":
		return Scissors
	}

	handleError(errors.New("Unknown move: " + s))
	return 0
}

func getOutcome(s string) RPSOutcome {
	switch strings.ToUpper(s) {
	case "X":
		return Loss
	case "Y":
		return Draw
	case "Z":
		return Win
	}

	handleError(errors.New("Unknown outcome: " + s))
	return 0
}

func (r *Round) Score() int {
	return int(r.Ours) + int(r.Outcome())
}

func (r *Round) SetOutMove() {
	for i := 1; i < 4; i++ {
		r.Ours = RPSMove(i)
		if r.Outcome() == r.DesiredOutcome {
			break
		}
	}
}

func (r *Round) Outcome() RPSOutcome {
	if r.Opponent == r.Ours {
		return Draw
	}

	switch r.Ours {
	case Rock:
		if r.Opponent == Scissors {
			return Win
		}
	case Paper:
		if r.Opponent == Rock {
			return Win
		}
	case Scissors:
		if r.Opponent == Paper {
			return Win
		}
	}

	return Loss
}

func getRounds(lines []string) []*Round {
	rounds := make([]*Round, 0)

	for _, line := range lines {
		if len(line) == 3 {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				round := &Round{
					Opponent:       getMove(parts[0]),
					Ours:           getMove(parts[1]),    // part 1
					DesiredOutcome: getOutcome(parts[1]), // part 2
				}
				rounds = append(rounds, round)
			}
		}
	}

	return rounds
}

func getLines(path string) []string {
	fileBytes, err := ioutil.ReadFile(path)
	handleError(err)

	file := string(fileBytes)
	file = strings.Trim(file, "\n ")
	lines := strings.Split(file, "\n")
	return lines
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %x \n", err)
		os.Exit(1)
	}
}
