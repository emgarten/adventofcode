package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var words = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	total := 0

	for _, r := range lines {
		// total += calibrationValue(digits(r))
		total += calibrationValue(digitsWithWords(r))
	}

	fmt.Printf("Entries: %v total: %v\n", len(lines), total)
}

// Concat first and last digits
func calibrationValue(d []int) int {
	// fmt.Printf("%v\n", d)
	return d[0]*10 + d[len(d)-1]
}

// Get all digits including words
func digitsWithWords(s string) []int {
	var digits []int

	// For each rune in the string check if it is a digit or the start of a word
	for i, r := range s {
		// Check for digits first
		if unicode.IsDigit(r) {
			digits = append(digits, runeToInt(r))
			continue
		}

		// Check if we are at a word
		for word, digit := range words {
			if strings.HasPrefix(s[i:], word) {
				digits = append(digits, digit)
				break
			}
		}
	}

	return digits
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
