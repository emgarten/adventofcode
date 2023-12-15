package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	total := 0

	// Only the first line
	input := lines[0]
	parts := strings.Split(input, ",")
	for _, p := range parts {
		x := hash(p)
		// fmt.Printf("%v: %v\n", p, x)
		total += x
	}

	fmt.Printf("rn: %d\n", hash("rn"))

	fmt.Printf("Lines: %v total: %v\n", len(lines), total)
}

// Hash
func hash(s string) int {
	x := 0
	for _, r := range s {
		x = calcHash(x, r)
	}
	return x
}

func calcHash(x int, r rune) int {
	// Add ascii value of rune to incoming val
	y := x + int(r)
	// current val * 17
	z := y * 17
	// return mod 256
	result := z % 256
	return result
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
