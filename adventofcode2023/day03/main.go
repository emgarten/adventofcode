package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
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
	SquareDigit
	SquareSymbol
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	m := getMatrix(lines)

	fmt.Printf("Lines: %v missing parts: %v\n", len(lines), sumMissingParts(m, 1))
	fmt.Printf("Lines: %v gear ratio: %v\n", len(lines), sumGearRatios(m))
}

// Part 2
func sumGearRatios(m Matrix) int {
	return walkMatrix(m, getGearRatio)
}

// Get the gear ratio of a symbol
func getGearRatio(m Matrix, i, j int, r rune) int {
	seen := make(map[string]bool)

	if r == '*' {
		nums := make([]int, 0)

		for _, c := range getAdj(m, i, j, j, 1) {
			if coordSquareType(m, c) == SquareDigit {
				row := c.I
				start := numStart(m, row, c.J)

				key := fmt.Sprintf("%d|%d", row, start)
				if !seen[key] {
					seen[key] = true
					nums = append(nums, getNum(m, row, start))
				}
			}
		}

		if len(nums) == 2 {
			return nums[0] * nums[1]
		}
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

// Part 1
func sumMissingParts(m Matrix, distance int) int {
	starts := allNumStarts(m)
	total := 0

	for _, s := range starts {
		if hasAdjSymbols(m, s, distance) {
			total += getNum(m, s.I, s.J)
		}

		fmt.Printf("%d %v\n", getNum(m, s.I, s.J), hasAdjSymbols(m, s, distance))
	}

	return total
}

// True if the target has adjacent symbols within the distance
func hasAdjSymbols(m Matrix, c Coords, distance int) bool {
	// Get all adjacent squares
	adj := getAdj(m, c.I, c.J, numEnd(m, c.I, c.J), distance)

	// Get the square type of each adjacent square
	for _, a := range adj {
		if coordSquareType(m, a) == SquareSymbol {
			return true
		}
	}

	return false
}

// Walk the matrix and return all number starting coordinates
func allNumStarts(m Matrix) []Coords {
	starts := make([]Coords, 0)
	last := -1

	for i, row := range m {
		last = -1

		for j, r := range row {
			if squareType(r) == SquareDigit {
				start := numStart(m, i, j)
				if start != last {
					last = start
					starts = append(starts, Coords{
						I: i,
						J: j,
					})
				}
			}
		}
	}

	return starts
}

// Return the coordinates of all adjacent squares, not include the target squares
func getAdj(m Matrix, i, start, end, distance int) []Coords {
	adj := make([]Coords, 0)

	// Find all adjacent squares, define offsets relative to the start digit
	for rowOffset := (0 - distance); rowOffset <= distance; rowOffset++ {
		mRow := i + rowOffset
		if mRow < 0 || mRow >= len(m) {
			continue
		}

		for colOffset := -1; colOffset <= (end - start + distance); colOffset++ {
			mCol := start + colOffset
			if mCol < 0 || mCol >= len(m[mRow]) {
				continue
			}

			// Skip the digits of the actual number we are looking at
			if rowOffset == 0 && (mCol >= start && mCol <= end) {
				continue
			}

			// Valid adjacent square
			adj = append(adj, Coords{
				I: mRow,
				J: mCol,
			})
		}
	}

	return adj
}

// Read the surrounding number from the matrix
func getNum(m Matrix, i, j int) int {
	start := numStart(m, i, j)
	end := numEnd(m, i, j)

	if start == -1 || end == -1 {
		return -1
	}

	num := 0
	for k := start; k <= end; k++ {
		x := runeToInt(m[i][k])
		num = num*10 + x
	}

	return num
}

// Move backwards to find the start of the digit
func numStart(m Matrix, i, j int) int {
	pos := -1

	for k := j; k >= 0; k-- {
		if squareType(m[i][k]) == SquareDigit {
			pos = k
		} else {
			break
		}
	}

	return pos
}

// Move forwards to find the end of the digit
func numEnd(m Matrix, i, j int) int {
	pos := -1

	for k := j; k < len(m[i]); k++ {
		if squareType(m[i][k]) == SquareDigit {
			pos = k
		} else {
			break
		}
	}

	return pos
}

func coordSquareType(m Matrix, c Coords) SquareType {
	return squareType(m[c.I][c.J])
}

func squareType(r rune) SquareType {
	if r == '.' {
		return SquareEmpty
	}

	if unicode.IsDigit(r) {
		return SquareDigit
	}

	return SquareSymbol
}

// Get all digits from a string
func digits(s string) []int {
	var digits []int
	for _, r := range s {
		if unicode.IsDigit(r) {
			digits = append(digits, runeToInt(r))
		}
	}

	return digits
}

// convert ascii char # to digit
func runeToInt(r rune) int {
	return int(r - '0')
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
