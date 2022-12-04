package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type (
	Entry struct {
		A Range
		B Range
	}

	Range struct {
		Start int
		End   int
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	entries := getEntries(lines)
	total := 0

	for _, r := range entries {
		if r.hasIntersection() {
			total++
		}
	}

	fmt.Printf("Entries: %v total: %v\n", len(entries), total)
}

func (e *Entry) hasIntersection() bool {
	return e.A.hasIntersection(e.B)
}

func (e *Entry) hasFullOverlap() bool {
	return e.A.isSubsetOrEqual(e.B) || e.B.isSubsetOrEqual(e.A)
}

func (r *Range) isSubsetOrEqual(a Range) bool {
	return r.Start <= a.Start && r.End >= a.End
}

func (r *Range) hasIntersection(a Range) bool {
	return r.isPointInRange(a.Start) || r.isPointInRange(a.End) || a.isPointInRange(r.Start) || a.isPointInRange(r.End)
}

func (r *Range) isPointInRange(x int) bool {
	return r.Start <= x && r.End >= x
}

func getEntries(lines []string) []*Entry {
	entries := make([]*Entry, 0)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		parts := strings.Split(line, ",")
		a := strings.Split(parts[0], "-")
		b := strings.Split(parts[1], "-")

		entries = append(entries, &Entry{
			A: Range{
				Start: getInt(a[0]),
				End:   getInt(a[1]),
			},
			B: Range{
				Start: getInt(b[0]),
				End:   getInt(b[1]),
			},
		})
	}

	return entries
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

func getInt(s string) int {
	x, err := strconv.Atoi(s)
	handleError(err)
	return x
}
