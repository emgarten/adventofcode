package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type (
	workflow struct {
		id    string
		flows []*flow
	}

	flow struct {
		exp        *exp // nil if always true
		ifTrueGoto string
	}

	exp struct {
		left  string
		op    string
		right string
	}

	entry map[string]int

	entrySets map[string][]int

	pathEntry struct {
		parent  *pathEntry
		isTrue  []*exp
		isFalse []*exp
	}

	// Inclusive
	entryRange struct {
		min int
		max int
	}

	// All ranges for each letter
	resultRanges map[string][]entryRange

	permCombo struct {
		x entryRange
		m entryRange
		a entryRange
		s entryRange
	}
)

func main() {
	path := os.Args[1]
	fmt.Println("reading: ", path)
	lines := getLines(path)

	// Parse
	entries := make([]entry, 0)
	flows := make(map[string]workflow)
	for _, line := range lines {
		if strings.HasPrefix(line, "{") {
			entries = append(entries, parseEntry(line))
		} else if strings.HasSuffix(line, "}") {
			w := parseWorkflow(line)
			flows[w.id] = w
		}
	}

	// Part 1
	total := 0
	for _, e := range entries {
		if accept(e, flows) {
			total += e.sum()
		}
	}

	fmt.Printf("Entries: %v Flows: %v total: %v\n", len(entries), len(flows), total)

	// Part 2
	paths := getPaths(flows, nil, "in")
	fmt.Printf("Paths: %v\n", len(paths))

	// Evaluate sets of all possible values for each letter
	// The resulting sets will be the ranges of values that are accepted
	// Arrays containing the actual numbers
	entrySetResults := getEntrySetResults(paths)

	fmt.Println("Entry set results: ", len(entrySetResults))

	// Remove duplication between sets
	// x -> set id
	// Mark the end
	// Add to the result to each set that had it
	pathResultRanges := getResultRanges(entrySetResults)

	fmt.Println("pathResultRanges: ", len(pathResultRanges))

	// Permutations
	totalPerms := getPermutationCount(pathResultRanges)

	fmt.Printf("Total: %v\n", totalPerms)
}

func (p *permCombo) count() uint64 {
	return uint64(p.x.count() * p.m.count() * p.a.count() * p.s.count())
}

// This takes a while
// To avoid using too much memory find all permutations for a single range
// from X at a time.
// There are probably way better ways to do this, but this works
func getPermutationCount(ranges []resultRanges) uint64 {
	fmt.Printf("Ranges: %v\n", len(ranges))

	totalPerms := uint64(0)
	rangesToId := mapRangeToSetIds(ranges, "x")
	for x, setIds := range rangesToId {
		fmt.Printf("Finding perms for X: %v\n", x)
		perms := make(map[permCombo]struct{})

		for _, setId := range setIds {
			r := ranges[setId]
			for _, m := range r["m"] {
				for _, a := range r["a"] {
					for _, s := range r["s"] {
						p := permCombo{
							x: x,
							m: m,
							a: a,
							s: s,
						}

						perms[p] = struct{}{}
					}
				}
			}
		}

		for k := range perms {
			totalPerms += k.count()
		}
	}

	return totalPerms
}

// Map of range -> set ids
func mapRangeToSetIds(ranges []resultRanges, letter string) map[entryRange][]int {
	// Range entryRange -> set ids
	m := make(map[entryRange][]int)

	for setId, r := range ranges {
		for _, rangeEntry := range r[letter] {
			m[rangeEntry] = append(m[rangeEntry], setId)
		}
	}

	return m
}

// Count unique propr values between a set of set ids
// Unique props across the union of set ids
func countForSetIds(ranges []resultRanges, setIds []int, letter string) int {
	// id -> count
	seen := make(map[entryRange]struct{})

	for _, setId := range setIds {
		for _, r := range ranges[setId][letter] {
			seen[r] = struct{}{}
		}
	}

	x := 0
	for r := range seen {
		x += r.count()
	}
	return x
}

func (r *entryRange) count() int {
	return r.max - r.min + 1
}

func getEntrySetResults(paths []*pathEntry) []entrySets {
	entrySetResults := make([]entrySets, 0)
	for _, p := range paths {
		es := make(entrySets)
		for _, c := range []string{"x", "m", "a", "s"} {
			es[c] = createSet(4000)
		}

		cur := p
		for cur != nil {

			for _, e := range cur.isTrue {
				e.evalSet(es, false)
			}

			for _, e := range cur.isFalse {
				e.evalSet(es, true)
			}

			cur = cur.parent
		}

		entrySetResults = append(entrySetResults, es)
	}
	return entrySetResults
}

// Returns an array that is set id -> ranges
func getResultRanges(entrySetResults []entrySets) []resultRanges {
	// id is the order of the entrySetResults above
	// id -> ranges
	pathResultRanges := make([]resultRanges, len(entrySetResults))

	// init resultRanges
	for i := 0; i < len(entrySetResults); i++ {
		pathResultRanges[i] = make(resultRanges)

		for _, c := range []string{"x", "m", "a", "s"} {
			pathResultRanges[i][c] = make([]entryRange, 0)
		}
	}

	for _, c := range []string{"x", "m", "a", "s"} {
		numToSet := make(map[int][]int)

		for esId, es := range entrySetResults {
			for _, v := range es[c] {
				numToSet[v] = append(numToSet[v], esId)
			}
		}

		var curRange *entryRange
		curSets := make([]int, 0)
		for i := 1; i <= 4000; i++ {
			endRange := false

			sets, ok := numToSet[i]
			if ok {
				if curRange == nil {
					curRange = &entryRange{
						min: i,
					}
					curSets = sets
					continue
				}

				if slices.Equal(curSets, sets) {
					continue
				}

				endRange = true
			} else {
				if curRange != nil {
					endRange = true
				}
			}

			if endRange {
				curRange.max = i - 1

				for _, setId := range curSets {
					pathResultRanges[setId][c] = append(pathResultRanges[setId][c], *curRange)
				}

				if len(sets) > 0 {
					// Some sets are still valid
					curRange = &entryRange{
						min: i,
					}
					curSets = sets
				} else {
					// No one needs this part of the range
					curRange = nil
				}
			}
		}

		if curRange != nil {
			curRange.max = 4000
			for _, setId := range curSets {
				pathResultRanges[setId][c] = append(pathResultRanges[setId][c], *curRange)
			}
		}
	}

	return pathResultRanges
}

func (p *pathEntry) String() string {
	steps := make([]string, 0)
	steps = append(steps, "A")
	cur := p
	for cur != nil {
		condParts := make([]string, 0)
		if len(cur.isTrue) > 0 {
			condParts = append(condParts, fmt.Sprintf("(%v)", getExpArrayStrings(cur.isTrue)))
		}

		if len(cur.isFalse) > 0 {
			condParts = append(condParts, fmt.Sprintf("(!%v)", getExpArrayStrings(cur.isFalse)))
		}

		steps = append(steps, strings.Join(condParts, " && "))
		cur = cur.parent
	}

	return strings.Join(steps, " -> ")
}

func getExpArrayStrings(a []*exp) []string {
	r := make([]string, 0)
	for _, e := range a {
		r = append(r, e.String())
	}
	return r
}

func (e *exp) String() string {
	return fmt.Sprintf("%v%v%v", e.left, e.op, e.right)
}

func (r *entryRange) String() string {
	return fmt.Sprintf("[%d,%d]", r.min, r.max)
}

func createSet(max int) []int {
	set := make([]int, max)
	for i := 0; i < max; i++ {
		set[i] = i + 1
	}
	return set
}

func getPaths(m map[string]workflow, parent *pathEntry, target string) []*pathEntry {
	paths := make([]*pathEntry, 0)

	if target == "A" {
		// Accepted
		paths = append(paths, parent)
	} else if target == "R" {
		// Rejected, no need to add this
		return paths
	}

	// Lookup target and continue
	expSoFar := make([]*exp, 0)

	w := m[target]
	for _, f := range w.flows {
		p := &pathEntry{
			parent:  parent,
			isFalse: expSoFar, // Does this need a clone?
		}

		if f.exp != nil {
			p.isTrue = append(p.isTrue, f.exp)

			// For the next exp it will need to be false
			expSoFar = append(expSoFar, f.exp)
		}

		for _, child := range getPaths(m, p, f.ifTrueGoto) {
			paths = append(paths, child)
		}
	}

	return paths
}

func (e *entry) sum() int {
	x := 0
	for _, v := range *e {
		x += v
	}
	return x
}

func accept(input entry, flows map[string]workflow) bool {
	// Start at first workflow
	w := flows["in"]

	// Run workflow until we get to an accept or reject
	for {
		r := w.eval(input)
		switch r {
		case "A":
			return true
		case "R":
			return false
		default:
			w = flows[r]
		}
	}
}

func (w *workflow) eval(e entry) string {
	for _, f := range w.flows {
		if f.exp == nil || f.exp.eval(e) {
			return f.ifTrueGoto
		}
	}
	panic("No flow found")
}

func (e *exp) evalSet(sets entrySets, negate bool) {
	s := sets[e.left]
	y := getInt(e.right)

	var f func(int) bool

	switch e.op {
	case "<":
		f = func(x int) bool { return x < y }
	case ">":
		f = func(x int) bool { return x > y }
	}

	filtered := make([]int, 0)
	for _, v := range s {
		b := f(v)

		if (negate && !b) || (!negate && b) {
			filtered = append(filtered, v)
		}
	}
	sets[e.left] = filtered
}

func (e *exp) eval(input entry) bool {
	x := input[e.left]
	y := getInt(e.right)

	switch e.op {
	case "<":
		return x < y
	case ">":
		return x > y
	}
	panic("Unknown operator: " + e.op)
}

// px{a<2006:qkq,m>2090:A,rfg}
func parseWorkflow(line string) workflow {
	parts := strings.Split(line, "{")
	w := workflow{
		id:    parts[0],
		flows: make([]*flow, 0),
	}

	flowsSection := strings.Trim(parts[1], "}")
	for _, f := range strings.Split(flowsSection, ",") {
		e := flow{}

		fp := strings.Split(f, ":")
		a := fp[0]

		// No exp
		if len(fp) == 1 {
			e.ifTrueGoto = a
		} else {
			expression := exp{}
			expression.left = a[0:1]
			expression.op = a[1:2]
			expression.right = a[2:]

			e.exp = &expression
			e.ifTrueGoto = fp[1]
		}

		w.flows = append(w.flows, &e)
	}

	return w
}

// {x=787,m=2655,a=1222,s=2876}
func parseEntry(line string) entry {
	line = strings.Trim(line, "{}")
	e := entry{}
	for _, prop := range strings.Split(line, ",") {
		p := strings.Split(prop, "=")
		e[p[0]] = getInt(p[1])
	}
	return e
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
