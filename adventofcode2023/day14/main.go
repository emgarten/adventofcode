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
)

const (
	SquareEmpty     = '.'
	SquareRoundRock = 'O'
	SquareCubeRock  = '#'
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	m := getMatrix(lines)

	// input -> output
	cache := make(map[string]string)

	// Part 2
	// Spin cycles
	max := 1_000_000_000
	// max := 3
	for i := 0; i < max; i++ {
		if i%10_000 == 0 {
			fmt.Printf("Run %d\n", i)
		}

		// Get matrix key to check the cache
		input := matrixString(m)

		// Once we hit a cache hit, we no longer have to run the shifting calculations
		_, hitCache := cache[input]
		if hitCache {
			// Everything from here to the end will be in the cache

			// Calculate the cycle length
			// Even looking up the results in the cache is too slow
			target := input
			searchInput := input
			cycleLength := 0
			for j := i; j < max; j++ {
				c, ok := cache[searchInput]
				if !ok {
					panic("Cache miss!")
				}

				cycleLength++

				if c == target {
					fmt.Println("Cycle length: ", cycleLength)
					break
				}
				searchInput = c
			}

			// Now jump to the start of the last cycle within the max
			// Use our current position, then divide what is left by the cycle length
			// Take the remainder and subtract from the max to find where to start
			// Then we will continue running from there against the cache
			lastCycleStart := max - ((max - i) % cycleLength)

			// Start at the final cycle, input is the same since we are in a repeating cycle
			for j := lastCycleStart; j < max; j++ {
				if j%1_000_000 == 0 {
					fmt.Printf("Cache run %d\n", j)
				}

				c, ok := cache[input]
				if !ok {
					panic("Cache miss!")
				}

				// Use the cached value as the next input
				input = c
			}

			// Convert to matrix for the final load calculation
			m = matrixFromString(input)
			break
		}

		// Shift round rocks
		fullShift(m, shiftNorth)
		fullShift(m, shiftWest)
		fullShift(m, shiftSouth)
		fullShift(m, shiftEast)

		// After cycling
		output := matrixString(m)

		// cache result
		cache[input] = output
	}

	// Part 1
	// for walkMatrix(m, shiftNorth) > 0 {
	// 	// Shift round rocks
	// }

	// Calculate the load
	load := walkMatrix(m, calcLoad)

	fmt.Printf("Lines: %v load: %v\n", len(lines), load)
}

func fullShift(m Matrix, f func(m Matrix, i, j int, r rune) int) {
	for walkMatrix(m, f) > 0 {
		// Shift round rocks until they stop
	}
}

func printMatrix(m Matrix) {
	fmt.Println("-------------------")
	for _, row := range m {
		for _, r := range row {
			fmt.Printf("%c", r)
		}
		fmt.Println()
	}
}

func matrixFromString(s string) Matrix {
	return getMatrix(strings.Split(s, "\n"))
}

// Get the original input back
func matrixString(m Matrix) string {
	lines := make([]string, 0)
	for _, row := range m {
		s := ""
		for _, r := range row {
			s += string(r)
		}
		lines = append(lines, s)
	}

	return strings.Join(lines, "\n")
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

// Base function for shifting
func shiftOffset(m Matrix, i, j int, r rune, iOffset, jOffset int) int {
	if m[i][j] != SquareRoundRock {
		return 0
	}

	targetI := i + iOffset
	targetJ := j + jOffset
	if targetI < 0 || targetJ < 0 || targetI >= len(m) || targetJ >= len(m[0]) {
		return 0
	}

	if m[targetI][targetJ] != SquareEmpty {
		return 0
	}

	// Roll round rocks if the space is empty
	m[targetI][targetJ] = SquareRoundRock
	m[i][j] = SquareEmpty

	// Indicate a change
	return 1
}

func shiftNorth(m Matrix, i, j int, r rune) int {
	return shiftOffset(m, i, j, r, -1, 0)
}

func shiftSouth(m Matrix, i, j int, r rune) int {
	return shiftOffset(m, i, j, r, 1, 0)
}

func shiftWest(m Matrix, i, j int, r rune) int {
	return shiftOffset(m, i, j, r, 0, -1)
}

func shiftEast(m Matrix, i, j int, r rune) int {
	return shiftOffset(m, i, j, r, 0, 1)
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

// Walk and sum up the matrix
// func walkMatrixRev(m Matrix, f func(m Matrix, i, j int, r rune) int) int {
// 	total := 0

// 	for i := len(m) - 1; i >= 0; i-- {
// 		for j := len(m[i]) - 1; j >= 0; j-- {
// 			total += f(m, i, j, m[i][j])
// 		}
// 	}

// 	return total
// }

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
