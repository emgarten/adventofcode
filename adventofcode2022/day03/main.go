package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type (
	Entry struct {
		Items []string
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	entries := getEntries2(lines)
	total := 0

	for _, r := range entries {
		match := r.Match()
		fmt.Printf("a: %v b: %v match: %v \n", r.Items[0], r.Items[1], match)
		total += match
	}

	fmt.Printf("Entries: %v score: %v\n", len(entries), total)
}

func getEntries(lines []string) []*Entry {
	entries := make([]*Entry, 0)

	for _, line := range lines {
		len := len(line)
		if len > 0 {
			entry := &Entry{
				Items: []string{line[:(len / 2)], line[(len / 2):]},
			}
			entries = append(entries, entry)
		}
	}

	return entries
}

func getEntries2(lines []string) []*Entry {
	entries := make([]*Entry, 0)

	for i := 0; i < len(lines); i++ {
		if (i+1)%3 == 0 {
			entries = append(entries, &Entry{
				Items: []string{
					lines[i-2],
					lines[i-1],
					lines[i],
				},
			})
		}
	}

	return entries
}

func (e *Entry) Match() int {
	found := make(map[int]uint8)
	allFlags := uint8(0)

	for i, s := range e.Items {
		flag := uint8(1 << i)
		allFlags |= flag

		for _, p := range getPriorities(s) {
			found[p] |= flag
		}
	}

	for k, v := range found {
		if v == allFlags {
			return k
		}
	}

	return -1
}

func getPriorities(s string) []int {
	p := make([]int, len(s))
	for i := range s {
		c := s[i]

		if c > 90 {
			// lower
			p[i] = int(c) - 96
		} else {
			// upper
			p[i] = int(c) - 38
		}
	}

	sort.Ints(p)

	return p
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
