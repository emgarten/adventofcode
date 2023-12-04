package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type (
	Card struct {
		Id    int
		Win   map[int]bool
		Other map[int]bool
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)

	m := make(map[int]Card)
	for _, line := range lines {
		c := parse(line)
		m[c.Id] = c
	}

	// Part 1
	total := 0
	for _, c := range m {
		x := scoreCard(c)
		total += x
		// fmt.Printf("Card %v: %v\n", c.Id, x)
	}

	// Part 2
	// Start with one copy of each card
	orig := make([]int, 0)
	for _, c := range m {
		orig = append(orig, c.Id)
	}

	// Play each card and cards won
	plays := playCount(m, orig)

	fmt.Printf("Lines: %v total: %v plays: %v\n", len(lines), total, plays)
}

// Recursively play the cards and track how many cards were played
func playCount(cards map[int]Card, play []int) int {
	plays := 0

	for _, id := range play {
		// id may not be valid, these should be skipped
		c, ok := cards[id]
		if ok {
			// Count this current card
			plays++
			// fmt.Printf("Card %v\n", id)

			// Wins add the cards after this current one
			next := make([]int, 0)
			for i := 0; i < len(intersect(c.Win, c.Other)); i++ {
				// First win is 0, so add 1 to get the next card
				next = append(next, id+i+1)
			}

			plays += playCount(cards, next)
		}
	}

	return plays
}

// Double the score for each matching number
// Starting with 1
// 2^0 = 1
// 2^1 = 2
func scoreCard(c Card) int {
	return int(math.Pow(2, float64(len(intersect(c.Win, c.Other)))-1))
}

func intersect(a, b map[int]bool) map[int]bool {
	m := make(map[int]bool)
	for k := range a {
		if b[k] {
			m[k] = true
		}
	}
	return m
}

func parse(line string) Card {
	card := Card{}

	// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
	parts := strings.Split(line, ":")
	card.Id = getInt(strings.TrimPrefix(parts[0], "Card "))
	for _, entry := range strings.Split(parts[1], "|") {
		m := make(map[int]bool)

		for _, s := range strings.Split(entry, " ") {
			if len(s) > 0 {
				m[getInt(s)] = true
			}
		}

		if len(card.Win) == 0 {
			card.Win = m
		} else {
			card.Other = m
		}
	}

	return card
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
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}

func getInt(s string) int {
	x, err := strconv.Atoi(strings.Trim(s, " "))
	if err != nil {
		fmt.Println("Error converting to int: ", s)
	}

	handleError(err)
	return x
}
