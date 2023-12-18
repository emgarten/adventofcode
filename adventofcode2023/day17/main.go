package main

import (
	"fmt"
	"math"
	"os"
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

type (
	// Matrix is a 2d array of runes
	// i is the row, j is the column
	Matrix [][]int

	PathEntry struct {
		Prev     *PathEntry
		I        int
		J        int
		Sum      int
		Distance int
	}

	CacheContext struct {
		Matrix     *Matrix
		MaxI       int
		MaxJ       int
		LowestPath map[string]*PathEntry
	}

	Coords struct {
		I int
		J int
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	m := Matrix(getMatrix(lines))

	lowest := lowestCostPath(&m)

	fmt.Printf("Lines: %v load: %v\n", len(lines), lowest)
}

func printMatrix(m Matrix) {
	fmt.Println("-------------------")
	for _, row := range m {
		for _, r := range row {
			fmt.Printf("%v", r)
		}
		fmt.Println()
	}
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

func lowestCostPath(m *Matrix) int {
	cacheContext := &CacheContext{
		Matrix:     m,
		MaxI:       len(*m),
		MaxJ:       len((*m)[0]),
		LowestPath: make(map[string]*PathEntry),
	}
	path := &PathEntry{
		Prev:     nil,
		I:        0,
		J:        0,
		Sum:      0,
		Distance: 999999999999,
	}
	target := Coords{cacheContext.MaxI - 1, cacheContext.MaxJ - 1}
	// target := Coords{0, 8}
	starting := []*PathEntry{path}
	// v := walkPreferLowest(cacheContext, starting, target)
	// printPath(cacheContext, v)
	// fmt.Println("Lowest path: ", v.Sum)
	// return v.Sum

	walk(cacheContext, starting, target)

	// Find the lowest path to the target
	targetKey := fmt.Sprintf("%v|%v|", target.I, target.J)
	lowest := -1
	for k, v := range cacheContext.LowestPath {
		if strings.HasPrefix(k, targetKey) {
			// printPath(cacheContext, v)
			// fmt.Println("Lowest path: ", v.Sum, " ", getPathShort(cacheContext, v))
			// return v.Sum
			if lowest == -1 || v.Sum < lowest {
				lowest = v.Sum
			}
		}
	}

	return lowest
}

func walkPreferLowest(cacheContext *CacheContext, startingPaths []*PathEntry, target Coords) *PathEntry {
	paths := startingPaths
	var winner *PathEntry

	// Loop until we run out of paths to check
	for len(paths) > 0 {
		// Paths to expand on
		toWalk := make([]*PathEntry, 0)
		// lowestInPaths := -1
		lowestDistanceAway := -1

		// Check incoming paths
		for _, path := range paths {
			// Check if we hit the target
			if path.I == target.I && path.J == target.J {
				// We hit the target
				if winner == nil || path.Sum < winner.Sum {
					winner = path
					fmt.Println("New winner: ", path.Sum)
				}
			} else {
				// Continue walking this path only if it is lower than the current winner
				if winner == nil || path.Sum < winner.Sum {
					toWalk = append(toWalk, path)

					// if lowestInPaths == -1 || path.Sum < lowestInPaths {
					// 	lowestInPaths = path.Sum
					// }

					if lowestDistanceAway == -1 || path.Distance < lowestDistanceAway {
						lowestDistanceAway = path.Distance
					}
				}
			}
		}

		// Find all possible paths
		paths = make([]*PathEntry, 0)

		// Expand paths for the next round
		for _, path := range toWalk {
			if path.Distance == lowestDistanceAway {
				for _, p := range getNextPaths(cacheContext, path, target) {
					paths = append(paths, p)
				}
			} else {
				// Ignore this path for now and put it back for later
				paths = append(paths, path)
			}
		}
	}

	return winner
}

func walk(cacheContext *CacheContext, startingPaths []*PathEntry, target Coords) {
	paths := make(map[string]*PathEntry)
	best := -1
	var lowestDistancePath *PathEntry

	for _, path := range startingPaths {
		paths[getPathShort(cacheContext, path)] = path
	}

	cycles := 0

	// Loop until we run out of paths to check
	for len(paths) > 0 {
		cycles++
		if cycles%10_000 == 0 {
			fmt.Printf("Cycles: %v Paths: %v\n", cycles, len(paths))
		}

		// Find the shortest distance
		if lowestDistancePath == nil {
			// To purge
			purge := make([]string, 0)

			for _, path := range paths {
				if best > 0 && path.Sum >= best {
					purge = append(purge, getPathShort(cacheContext, path))
					continue
				}

				if lowestDistancePath == nil || path.Distance < lowestDistancePath.Distance {
					lowestDistancePath = path
				}
			}

			if len(purge) > 0 {
				// fmt.Println("Purging: ", len(purge))
				for _, s := range purge {
					delete(paths, s)
				}
			}
		}

		if lowestDistancePath == nil {
			// No more paths
			continue
		}

		path := lowestDistancePath
		pathKey := getPathShort(cacheContext, path)
		delete(paths, pathKey)
		lowestDistancePath = nil

		if best > 0 && path.Sum >= best {
			continue
		}

		// Check if we have a lower path already
		if path.Prev != nil {
			hasNewOptions := false
			for _, d := range getAllowed(cacheContext, path) {
				// Key the coords + the possible next move
				key := fmt.Sprintf("%d|%d|%v", path.I, path.J, dirStr(d))
				if lowestPath, ok := cacheContext.LowestPath[key]; ok {
					if lowestPath.Sum < path.Sum {
						// Ignore this path
						continue
					}

					hasNewOptions = true

					if lowestPath.Sum > path.Sum {
						// New lowest
						cacheContext.LowestPath[key] = path
					}
				} else {
					// New path
					cacheContext.LowestPath[key] = path
					hasNewOptions = true
				}
			}

			// Skip this path, everything is cached
			if !hasNewOptions {
				continue
			}
		}

		// Check if we hit the target
		if path.I == target.I && path.J == target.J {
			if best == -1 || path.Sum < best {
				best = path.Sum
				fmt.Printf("New best: %v Len: %v\n", path.Sum, getPathLength(path))
			}

			// We hit the target
			continue
		}

		// Double check that a parent path wasn't cached already
		// parentsValid := true
		// curPath := path.Prev
		// for curPath != nil && curPath.Prev != nil {
		// 	key := fmt.Sprintf("%d|%d|%v", curPath.Prev.I, curPath.Prev.J, dirStr(getDirection(curPath.Prev.I, curPath.Prev.J, curPath.I, curPath.J)))
		// 	if lowestPath, ok := cacheContext.LowestPath[key]; ok {
		// 		if curPath.Prev.Sum > lowestPath.Sum {
		// 			// Ignore this path, it is now too high
		// 			parentsValid = false
		// 			break
		// 		}
		// 	}
		// }

		// if !parentsValid {
		// 	continue
		// }

		// Expand all paths for the next round
		for _, p := range getNextPaths(cacheContext, path, target) {
			// Skip anything that can be ruled out already as too high
			if best == -1 || p.Sum < best {
				// Add all possible paths to the full list
				paths[getPathShort(cacheContext, p)] = p

				// Shortcut to set the next path
				if p.Distance < path.Distance && (lowestDistancePath == nil || p.Distance < lowestDistancePath.Distance) {
					lowestDistancePath = p
				}
			}
		}
	}
}

func walkOld(cacheContext *CacheContext, startingPaths []*PathEntry, target Coords) {
	paths := startingPaths
	best := -1

	// Loop until we run out of paths to check
	for len(paths) > 0 {
		// Paths to expand on
		toWalk := make([]*PathEntry, 0)
		lowestDistance := -1

		// Check incoming paths
		for _, path := range paths {
			// test := getPathShort(cacheContext, path)
			// // if test == "4115453231" {
			// if test == ">4>1v1>5>4>5^3>2>3>1" {
			// 	fmt.Println("Found it!")
			// }

			if lowestDistance == -1 || path.Distance < lowestDistance {
				lowestDistance = path.Distance
			}

			// Check if we have a lower path already
			if path.Prev != nil {
				key := fmt.Sprintf("%d|%d|%v", path.I, path.J, getBlockedString(cacheContext, path))
				if lowestPath, ok := cacheContext.LowestPath[key]; ok {
					if lowestPath.Sum < path.Sum {
						// fmt.Printf("Ignoring path: %v Lowest: %v Cur: %v\n", key, getPathShort(cacheContext, lowestPath), getPathShort(cacheContext, path))

						// test := getPathShort(cacheContext, path)
						// // if test == "4115453231" {
						// if test == ">4>1v1>5>4>5>3^2" {
						// 	fmt.Println("Found it!")
						// }

						// Ignore this path
						continue
					}

					if lowestPath.Sum > path.Sum {
						// New lowest
						cacheContext.LowestPath[key] = path
					}
				} else {
					// New path
					cacheContext.LowestPath[key] = path
				}
			}

			// Check if we hit the target
			if path.I == target.I && path.J == target.J {
				if best == -1 || path.Sum < best {
					best = path.Sum
					fmt.Println("New best: ", path.Sum)
				}

				// We hit the target
				continue
			}

			if best == -1 || path.Sum < best {
				// Continue walking this path
				toWalk = append(toWalk, path)
			}
		}

		// Find all possible paths
		paths = make([]*PathEntry, 0)

		fmt.Println("Lowest distance: ", lowestDistance)

		// Expand all paths for the next round
		for _, path := range toWalk {
			// test := getPathShort(cacheContext, path)
			// if test == "4115453231" {
			// if test == ">4>1v1>5>4>5" {
			// 	fmt.Println("Found it!")
			// }

			if path.Distance == lowestDistance {
				for _, p := range getNextPaths(cacheContext, path, target) {
					paths = append(paths, p)
				}
			} else {
				// Ignore this path for now and put it back for later
				paths = append(paths, path)
			}
		}
	}
}

// Find the distance between two points
func getDistance(path *PathEntry, target Coords) int {
	return int(math.Abs(float64(target.I)-float64(path.I)) + math.Abs(float64(target.J)-float64(path.J)))
}

// Find the distance between two points
func getPathLength(path *PathEntry) int {
	x := 0
	for path != nil {
		x++
		path = path.Prev
	}
	return x
}

// Find all squares that could be visited next
func getNextPaths(cacheContext *CacheContext, path *PathEntry, target Coords) []*PathEntry {
	toWalk := make([]*PathEntry, 0)

	// Find everything around
	for _, c := range allPossible(path.I, path.J) {
		newPath := &PathEntry{
			Prev: path,
			I:    c.I,
			J:    c.J,
		}

		if isPossible(cacheContext, newPath) {
			newPath.Sum = sumPath(cacheContext, newPath)
			newPath.Distance = getDistance(newPath, target)
			toWalk = append(toWalk, newPath)
		}
	}

	return toWalk
}

func allPossible(i, j int) []Coords {
	return []Coords{
		{i - 1, j},
		{i + 1, j},
		{i, j - 1},
		{i, j + 1},
	}
}

func isPossible(c *CacheContext, path *PathEntry) bool {
	i := path.I
	j := path.J

	// Check bounds
	if i >= 0 && i < c.MaxI && j >= 0 && j < c.MaxJ {
		// Ensure we don't go back or cross the same path
		if path.Prev == nil || !isVisited(path.Prev, i, j) {
			// Ensure we don't go in the same direction for more than 3 squares
			if x, _ := getLastRepeatedDirection(path); x <= 3 {
				return true
			}
		}
	}

	// Invalid
	return false
}

func getAllowed(c *CacheContext, path *PathEntry) []Direction {
	r := make([]Direction, 0)

	// Up
	up := &PathEntry{
		Prev: path,
		I:    path.I - 1,
		J:    path.J,
	}
	if isPossible(c, up) {
		r = append(r, Up)
	}

	// Down
	down := &PathEntry{
		Prev: path,
		I:    path.I + 1,
		J:    path.J,
	}
	if isPossible(c, down) {
		r = append(r, Down)
	}

	// Left
	left := &PathEntry{
		Prev: path,
		I:    path.I,
		J:    path.J - 1,
	}
	if isPossible(c, left) {
		r = append(r, Left)
	}

	// Right
	right := &PathEntry{
		Prev: path,
		I:    path.I,
		J:    path.J + 1,
	}
	if isPossible(c, right) {
		r = append(r, Right)
	}

	return r
}

func getBlockedString(c *CacheContext, path *PathEntry) string {
	s := ""

	// Up
	up := &PathEntry{
		Prev: path,
		I:    path.I - 1,
		J:    path.J,
	}
	if !isPossible(c, up) {
		s += dirStr(Up)
	}

	// Down
	down := &PathEntry{
		Prev: path,
		I:    path.I + 1,
		J:    path.J,
	}
	if !isPossible(c, down) {
		s += dirStr(Down)
	}

	// Left
	left := &PathEntry{
		Prev: path,
		I:    path.I,
		J:    path.J - 1,
	}
	if !isPossible(c, left) {
		s += dirStr(Left)
	}

	// Right
	right := &PathEntry{
		Prev: path,
		I:    path.I,
		J:    path.J + 1,
	}
	if !isPossible(c, right) {
		s += dirStr(Right)
	}

	return s
}

// Shows either the last direction, or the last repeated directions
func getLastRepeatedDirectionString(path *PathEntry) string {
	s := ""

	x, d := getLastRepeatedDirection(path)
	for i := 0; i < x; i++ {
		s += dirStr(d)
	}

	return s
}

// Get count of existing moves in the same direction
func getLastRepeatedDirection(path *PathEntry) (int, Direction) {
	var lastDirection *Direction

	count := 0
	for path != nil && path.Prev != nil {
		d := getDirection(path.Prev.I, path.Prev.J, path.I, path.J)

		if lastDirection == nil {
			lastDirection = &d
			count++
		} else if *lastDirection == d {
			count++
		} else {
			// Change direction
			break
		}
		path = path.Prev
	}

	return count, *lastDirection
}

func getDirection(prevI, prevJ, i, j int) Direction {
	if prevI == i {
		if prevJ > j {
			return Left
		} else {
			return Right
		}
	} else {
		if prevI > i {
			return Up
		} else {
			return Down
		}
	}
}

func isVisited(path *PathEntry, i, j int) bool {
	for path != nil {
		if path.I == i && path.J == j {
			return true
		}
		path = path.Prev
	}

	return false
}

func isBackwards(path *PathEntry, i, j int) bool {
	for path != nil && path.Prev != nil {
		if path.I == i && path.J == j {
			return true
		}
	}

	return false
}

func sumPath(c *CacheContext, path *PathEntry) int {
	x := 0

	// for path != nil && path.Prev != nil {
	// 	x += (*c.Matrix)[path.I][path.J]
	// 	path = path.Prev
	// }

	if path != nil && path.Prev != nil {
		x = path.Prev.Sum + (*c.Matrix)[path.I][path.J]
	}

	return x
}

func printPath(c *CacheContext, path *PathEntry) {
	fmt.Println("-------------------")
	for path != nil && path.Prev != nil {
		fmt.Printf("(%v, %v) %v Val: %v\n", path.I, path.J, dirStr(getDirection(path.Prev.I, path.Prev.J, path.I, path.J)), (*c.Matrix)[path.I][path.J])
		path = path.Prev
	}
}

func getPathShort(c *CacheContext, path *PathEntry) string {
	s := ""
	for path != nil && path.Prev != nil {
		s += fmt.Sprintf("%v%v", (*c.Matrix)[path.I][path.J], dirStr(getDirection(path.Prev.I, path.Prev.J, path.I, path.J)))
		path = path.Prev
	}

	r := ""
	for _, c := range s {
		r = string(c) + r
	}
	return r
}

func dirStr(d Direction) string {
	switch d {
	case Up:
		return "^"
	case Down:
		return "v"
	case Left:
		return "<"
	case Right:
		return ">"
	default:
		return "x"
	}
}

// create a 2d matrix of the input
func getMatrix(lines []string) [][]int {
	matrix := make([][]int, len(lines))

	for i, line := range lines {
		matrix[i] = make([]int, len(line))
		for j, r := range line {
			matrix[i][j] = runeToInt(r)
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

func getInt(s string) int {
	x, err := strconv.Atoi(strings.Trim(s, " "))
	if err != nil {
		fmt.Println("Error converting to int: ", s)
	}

	handleError(err)
	return x
}

// convert ascii char # to digit
func runeToInt(r rune) int {
	return int(r - '0')
}
