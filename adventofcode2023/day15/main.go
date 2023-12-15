package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	Lens struct {
		Val         string // string to hash
		FocalLength int    // 0 for remove
	}

	LensPos struct {
		Lens Lens
		Pos  int // Position will have gaps, what matters is that this is incrementing. Higher = farther back
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	part1total := 0
	boxes := make(map[int][]Lens, 256)

	// Only the first line
	input := lines[0]
	parts := strings.Split(input, ",")
	for _, p := range parts {
		// Part 1
		x := hash(p)
		// fmt.Printf("%v: %v\n", p, x)
		part1total += x

		// Part 2
		// Sort lenses into boxes
		lens := parseLens(p)
		h := hash(lens.Val)
		fmt.Printf("Lens: %v Hash: %v\n", lens.Val, h)
		boxes[h] = append(boxes[h], lens)
	}

	part2total := 0

	// Apply lense commands in each box
	for k, v := range boxes {
		fmt.Printf("Box %v: %v\n", k, v)

		// Process lenses in box
		lookup := make(map[string]LensPos)

		// Apply lens commands
		for i, lens := range v {
			if lens.FocalLength == 0 {
				// Remove
				delete(lookup, lens.Val)
			} else {
				existing, ok := lookup[lens.Val]
				if ok {
					// Replace existing
					lookup[lens.Val] = LensPos{
						Lens: lens,
						Pos:  existing.Pos, // Keep position
					}
				} else {
					// Add
					lookup[lens.Val] = LensPos{
						Lens: lens,
						Pos:  i, // Use position from array, this will always be the highest so far
					}
				}
			}
		}

		// Find slots
		slots := make([]int, 0)
		for i := 0; i < len(v); i++ {
			for _, lp := range lookup {
				if lp.Pos == i {
					slots = append(slots, lp.Lens.FocalLength)
				}
			}
		}

		for slot, focalLength := range slots {
			val := (k + 1) * (slot + 1) * focalLength
			fmt.Printf("Box: %v Slot %v: Focal: %v Val: %v\n", k, slot, focalLength, val)
			part2total += val
		}
	}

	fmt.Printf("Lines: %v total: %v part2: %v\n", len(lines), part1total, part2total)
}

func parseLens(s string) Lens {
	parts := strings.Split(s, "=")
	if len(parts) == 2 {
		return Lens{
			Val:         parts[0],
			FocalLength: getInt(parts[1]),
		}
	}

	return Lens{
		Val:         strings.TrimSuffix(s, "-"),
		FocalLength: 0,
	}
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
