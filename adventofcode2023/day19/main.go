package main

import (
	"fmt"
	"math"
	"os"
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

	total := 0
	for _, e := range entries {
		if accept(e, flows) {
			total += e.sum()
		}
	}

	fmt.Printf("Entries: %v Flows: %v total: %v\n", len(entries), len(flows), total)
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
