package main

import (
	"fmt"
	"os"
	"strings"
)

type SquareType int

type (
	// Matrix is a 2d array of runes
	// i is the row, j is the column
	Matrix [][]rune

	// Coordinates for a square
	Coords struct {
		I int
		J int
	}
)

const (
	SquareEmpty SquareType = iota + 1
	SquareRoundRock
	SquareCubeRock
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	m := getMatrix(lines)

	for walkMatrix(m, shiftUp) > 0 {
		// Shifting
	}

	// Calculate the load
	load := walkMatrix(m, calcLoad)

	fmt.Printf("Lines: %v load: %v\n", len(lines), load)
}

func calcLoad(m Matrix, i, j int, r rune) int {
	if r == 'O' {
		// Row 0 is the heaviest, count them as the length of the matrix
		// Remove 1 for each row closer to the bottom
		// Last row is 1
		return len(m) - i
	}

	return 0
}

func shiftUp(m Matrix, i, j int, r rune) int {
	// Top row cannot shift
	if i == 0 {
		return 0
	}

	// Roll round rocks upwards if the space is empty
	if coordSquareType(m, Coords{i, j}) == SquareRoundRock && coordSquareType(m, Coords{i - 1, j}) == SquareEmpty {
		m[i-1][j] = 'O'
		m[i][j] = '.'

		// Indicate a change
		return 1
	}

	return 0
}

// Walk and sum up the matrix
func walkMatrix(m Matrix, f func(m Matrix, i, j int, r rune) int) int {
	total := 0

	for i, row := range m {
		for j, r := range row {
			total += f(m, i, j, r)
		}
	}

	return total
}

func coordSquareType(m Matrix, c Coords) SquareType {
	return squareType(m[c.I][c.J])
}

func squareType(r rune) SquareType {
	switch r {
	case 'O':
		return SquareRoundRock
	case '#':
		return SquareCubeRock
	default:
		return SquareEmpty
	}
}

// create a 2d matrix of the input
func getMatrix(lines []string) [][]rune {
	matrix := make([][]rune, len(lines))

	for i, line := range lines {
		matrix[i] = make([]rune, len(line))
		for j, r := range line {
			matrix[i][j] = r
		}
	}

	return matrix
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
