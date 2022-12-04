package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type (
	elf struct {
		calories []int
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)

	// create an elf for each set of calories
	elves := getElves(lines)
	elfCount := len(elves)
	calories := make([]int, elfCount)

	for i, e := range elves {
		calories[i] = e.totalCalories()
	}

	// Sort inline
	sort.Ints(calories[:])

	highest := calories[elfCount-1]

	fmt.Printf("Total elves: %v highest calories: %v\n", len(elves), highest)

	top3 := calories[elfCount-1] + calories[elfCount-2] + calories[elfCount-3]

	fmt.Printf("Sum of top 3: %v\n", top3)
}

func (e *elf) totalCalories() int {
	x := 0
	for _, i := range e.calories {
		x += i
	}

	return x
}

func getElves(lines []string) []*elf {
	var curElf *elf
	elves := make([]*elf, 0)

	for _, line := range lines {
		// fmt.Println("line: ", line)

		if line != "" {
			if curElf == nil {
				curElf = &elf{
					calories: make([]int, 0),
				}
			}

			i, _ := strconv.Atoi(line)
			curElf.calories = append(curElf.calories, i)
		} else if curElf != nil {
			elves = append(elves, curElf)
			curElf = nil
		}
	}

	if curElf != nil {
		elves = append(elves, curElf)
		curElf = nil
	}

	return elves
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
