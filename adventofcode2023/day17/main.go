package main

import (
	"container/heap"
	"fmt"
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

var directions = []Direction{Up, Down, Left, Right}

type (
	// Matrix is a 2d array of runes
	// i is the row, j is the column
	Matrix [][]int

	PathStateEntry struct {
		HeatLoss int
		State    State
	}

	CacheContext struct {
		Matrix     *Matrix
		MaxI       int
		MaxJ       int
		LowestPath map[State]int
	}

	Coords struct {
		I int
		J int
	}

	State struct {
		I     int
		J     int
		Dir   Direction
		Moves int
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)
	m := Matrix(getMatrix(lines))

	lowest := lowestCostPath(&m)

	fmt.Printf("Lines: %v lowest: %v\n", len(lines), lowest)
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
		LowestPath: make(map[State]int),
	}
	target := Coords{cacheContext.MaxI - 1, cacheContext.MaxJ - 1}
	// target := Coords{0, 8}
	start := Coords{0, 0}

	return getShortestPath(cacheContext, start, target)
}

// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func getShortestPath(cacheContext *CacheContext, start, target Coords) int {
	pq := make(PQ[PathStateEntry], 0)
	heap.Init(&pq)

	pq.HeapPush(PathStateEntry{HeatLoss: 0, State: State{I: 0, J: 0, Dir: Right, Moves: 0}}, 0)
	pq.HeapPush(PathStateEntry{HeatLoss: 0, State: State{I: 0, J: 0, Dir: Down, Moves: 0}}, 0)

	cycle := 0

	for pq.Len() > 0 {
		cycle++

		// Take the lowest cost path
		v, _ := pq.HeapPop()

		if cycle%1_000_000 == 0 {
			fmt.Printf("Cycle: %v Cur: %v\n", cycle, v.HeatLoss)
		}

		lowest, ok := cacheContext.LowestPath[v.State]
		if ok {
			if lowest < v.HeatLoss {
				// We already have a lower cost path here
				continue
			}
		}

		// Check if we reached the target, if so, do not walk further
		if v.State.I == target.I && v.State.J == target.J {
			fmt.Printf("Hit target: %v\n", v.HeatLoss)
			return v.HeatLoss
		}

		// Find neighbors
		// for _, n := range allPossible(v.State.I, v.State.J) {
		for _, d := range directions {
			n := getNeighborOrNil(cacheContext, &v.State, d)
			if n == nil {
				continue
			}

			// Calculate heat loss
			heatLoss := v.HeatLoss + (*cacheContext.Matrix)[n.I][n.J]

			lowest, ok := cacheContext.LowestPath[*n]
			if ok {
				if lowest < heatLoss {
					// We already have a lower cost path here
					continue
				}

				if lowest == heatLoss {
					// No need to re-do this
					continue
				}
			} else {
				// Add to cache
				cacheContext.LowestPath[*n] = heatLoss
			}

			// Add to priority queue to continue walking later
			p := PathStateEntry{
				HeatLoss: heatLoss,
				State:    *n,
			}

			pq.HeapPush(p, heatLoss)
		}
	}

	return -1
}

func revDirection(d Direction) Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	panic("Invalid direction")
}

func getNeighborOrNil(c *CacheContext, cur *State, d Direction) *State {
	// Do not go backwards
	if d == revDirection(cur.Dir) {
		return nil
	}

	// Move limit hit
	if d == cur.Dir && cur.Moves == 3 {
		return nil
	}

	r := &State{
		I:     cur.I,
		J:     cur.J,
		Dir:   d,
		Moves: 1,
	}

	if d == cur.Dir {
		r.Moves = cur.Moves + 1
	}

	switch d {
	case Up:
		r.I--
	case Down:
		r.I++
	case Left:
		r.J--
	case Right:
		r.J++
	}

	// Check bounds
	if r.I < 0 || r.I >= c.MaxI || r.J < 0 || r.J >= c.MaxJ {
		return nil
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

// Priority Queue based on https://pkg.go.dev/container/heap#pkg-examples
type pqi[T any] struct {
	v T
	p int
}

type PQ[T any] []pqi[T]

// Interface to use
func (q *PQ[T]) HeapPush(v T, p int) { heap.Push(q, pqi[T]{v, p}) }
func (q *PQ[T]) HeapPop() (T, int)   { x := heap.Pop(q).(pqi[T]); return x.v, x.p }

// Internal for Heap
func (q PQ[_]) Len() int           { return len(q) }
func (q PQ[_]) Less(i, j int) bool { return q[i].p < q[j].p }
func (q PQ[_]) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *PQ[T]) Push(x any)        { *q = append(*q, x.(pqi[T])) }
func (q *PQ[_]) Pop() (x any)      { x, *q = (*q)[len(*q)-1], (*q)[:len(*q)-1]; return x }
