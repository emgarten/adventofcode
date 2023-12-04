package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	Round map[string]int

	Game struct {
		Id     int
		Rounds []*Round
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)

	m := make(map[int]Game)
	for _, line := range lines {
		game := parseGame(line)
		m[game.Id] = game
	}

	total := 0
	for _, game := range m {
		if isPossible(game) {
			total += game.Id
		}
	}

	totalPower := 0
	for _, game := range m {
		totalPower += gamePower(game)
	}

	fmt.Printf("Lines: %v total: %v power: %v\n", len(lines), total, totalPower)
}

func gamePower(g Game) int {
	red := maxByColor(g, "red")
	green := maxByColor(g, "green")
	blue := maxByColor(g, "blue")

	return red * green * blue
}

func maxByColor(g Game, c string) int {
	max := 0
	for _, round := range g.Rounds {
		if (*round)[c] > max {
			max = (*round)[c]
		}
	}

	return max
}

// only 12 red cubes, 13 green cubes, and 14 blue cubes
func isPossible(g Game) bool {
	for _, round := range g.Rounds {
		if (*round)["red"] > 12 || (*round)["green"] > 13 || (*round)["blue"] > 14 {
			return false
		}
	}

	return true
}

func parseGame(line string) Game {
	game := Game{}

	// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
	parts := strings.Split(line, ":")
	game.Id = getInt(strings.TrimPrefix(parts[0], "Game "))
	for _, roundEntry := range strings.Split(parts[1], ";") {
		r := make(Round)

		for _, colorEntry := range strings.Split(roundEntry, ",") {
			s := strings.Trim(colorEntry, " ")
			colorParts := strings.Split(s, " ")
			r[colorParts[1]] = getInt(colorParts[0])
		}

		game.Rounds = append(game.Rounds, &r)
	}

	return game
}

// File -> lines
func getLines(path string) []string {
	fileBytes, err := os.ReadFile(path)
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

func getInt(s string) int {
	x, err := strconv.Atoi(s)
	handleError(err)
	return x
}
