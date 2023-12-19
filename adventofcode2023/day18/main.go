package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Direction int

const (
	Up Direction = iota + 1
	Down
	Left
	Right
)

var directions = []Direction{Up, Down, Left, Right}

type (
	point struct {
		x int
		y int
	}

	entry struct {
		d     rune
		len   int
		color string
	}

	line struct {
		p1 point
		p2 point
	}
)

var (
	lineRE = regexp.MustCompile(`^([A-Za-z])\s+(\d+)\s+\((#[0-9A-Fa-f]{6})\)$`)
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)

	// Parse
	entries := make([]entry, 0)
	for _, line := range lines {
		// Part 1
		l := parseLine(line)

		// Part 2
		l2 := parseHexLine(l)

		entries = append(entries, l2)
	}

	// Get points
	points := getPoints(entries)

	// Get lines
	pointLines := getLinesFromPoints(points)

	// Find area
	area := getArea(pointLines)

	fmt.Printf("Lines: %v area: %v\n", len(lines), area)
}

func (l line) len() int {
	return int(math.Sqrt(math.Pow(float64(l.p1.x-l.p2.x), 2) + math.Pow(float64(l.p1.y-l.p2.y), 2)))
}

// https://en.wikipedia.org/wiki/Shoelace_formula
func getArea(lines []line) int {
	sum := 0

	// Find the area of the inner polygon
	for i := 0; i < len(lines); i++ {
		sum += (lines[i].p1.x - lines[i].p2.x) * (lines[i].p1.y + lines[i].p2.y)
	}

	// Include the outer border
	for _, l := range lines {
		sum += l.len()
	}

	// Get the area
	return sum/2 + 1
}

func getLinesFromPoints(points []point) []line {
	lines := make([]line, 0)

	for i := 0; i < len(points); i++ {
		if i == 0 {
			// Wrap around to the end
			lines = append(lines, line{points[len(points)-1], points[i]})
		} else {
			lines = append(lines, line{points[i-1], points[i]})
		}
	}

	return lines
}

func getPoints(entries []entry) []point {
	points := make([]point, 0)

	// Start at 0,0
	lastPoint := point{0, 0}
	points = append(points, lastPoint)

	for _, entry := range entries {
		// Start at previous point
		p := point{lastPoint.x, lastPoint.y}

		// Make it inclusive to the square
		distance := entry.len

		// Move
		switch entry.d {
		case 'U':
			p.y -= distance
		case 'D':
			p.y += distance
		case 'L':
			p.x -= distance
		case 'R':
			p.x += distance
		}

		// Add point
		points = append(points, p)
		lastPoint = p
	}

	return points
}

func parseLine(line string) entry {
	matches := lineRE.FindStringSubmatch(line)
	return entry{
		d:     rune(matches[1][0]),
		len:   getInt(matches[2]),
		color: matches[3],
	}
}

func parseHexLine(e entry) entry {
	d := 'R'
	switch e.color[len(e.color)-1:] {
	case "0":
		d = 'R'
	case "1":
		d = 'D'
	case "2":
		d = 'L'
	case "3":
		d = 'U'
	}

	s := e.color[1 : len(e.color)-1]
	return entry{
		len:   getIntFromHex(s),
		color: "", // not used
		d:     d,
	}
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

func getIntFromHex(s string) int {
	x, err := strconv.ParseInt(strings.Trim(s, " "), 16, 64)
	if err != nil {
		fmt.Println("Error converting to int: ", s)
	}

	handleError(err)
	return int(x)
}

func getInt(s string) int {
	x, err := strconv.Atoi(strings.Trim(s, " "))
	if err != nil {
		fmt.Println("Error converting to int: ", s)
	}

	handleError(err)
	return x
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}
