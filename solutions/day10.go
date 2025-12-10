package solutions

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day10", &day10{})
}

type day10 struct{}

type machine10 struct {
	target      int
	buttonMasks []int
	joltages    []int
}

type queue struct {
	fifo [][2]int
	lifo [][2]int
}

func (q *queue) push(x [2]int) {
	q.lifo = append(q.lifo, x)
}
func (q *queue) pop() ([2]int, bool) {
	if len(q.fifo) > 0 {
		n := len(q.fifo)
		x := q.fifo[n-1]
		q.fifo = q.fifo[:n-1]
		return x, true
	}
	if len(q.lifo) == 0 {
		return [2]int{}, false
	}
	q.fifo = q.lifo
	slices.Reverse(q.fifo)
	q.lifo = nil
	return q.pop()
}

func (*day10) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.Trim(string(input), "\n"), "\n")
	machines := make([]machine10, 0)
	for _, line := range lines {
		elements := strings.Split(line, " ")
		if elements[0][0] != '[' || elements[0][len(elements[0])-1] != ']' {
			return fmt.Errorf("first element not wrapped in [] %q", elements[0])
		}
		targetMask := 0
		for i, c := range elements[0][1 : len(elements[0])-1] {
			if c == '#' {
				targetMask |= (1 << i)
			}
		}
		buttonMasks := make([]int, 0)
		for _, button := range elements[1 : len(elements)-1] {
			if button[0] != '(' || button[len(button)-1] != ')' {
				return fmt.Errorf("button not wrapped in () %q", button)
			}
			buttonNs := strings.Split(button[1:len(button)-1], ",")
			buttonMask := 0
			for _, buttonN := range buttonNs {
				n, err := strconv.Atoi(buttonN)
				if err != nil {
					return err
				}
				buttonMask |= (1 << n)
			}
			buttonMasks = append(buttonMasks, int(buttonMask))
		}
		lastElement := elements[len(elements)-1]
		if lastElement[0] != '{' || lastElement[len(lastElement)-1] != '}' {
			return fmt.Errorf("last element not wrapped in {} %q", lastElement)
		}
		joltages := make([]int, 0)
		joltageStrs := strings.Split(lastElement[1:len(lastElement)-1], ",")
		for _, joltageStr := range joltageStrs {
			n, err := strconv.Atoi(joltageStr)
			if err != nil {
				return err
			}
			joltages = append(joltages, n)
		}
		machines = append(machines, machine10{target: int(targetMask), buttonMasks: buttonMasks, joltages: joltages})
	}
	slog.Debug("input", "machines", machines)
	totalPresses := 0
machineLoop:
	for mi, machine := range machines {
		q := queue{}
		visited := make(map[int]bool)
		q.push([2]int{0, 0})
		visited[0] = true
		slog.Debug("starting machine", "index", mi)
		for {
			x, ok := q.pop()
			if !ok {
				return fmt.Errorf("pop failed?")
			}
			current := x[0]
			d := x[1]
			for _, buttonMask := range machine.buttonMasks {
				next := current ^ buttonMask
				if next == machine.target {
					totalPresses += d + 1
					slog.Debug("completing machine", "nPresses", d+1)
					continue machineLoop
				}
				if visited[next] {
					continue
				}
				visited[next] = true
				q.push([2]int{next, d + 1})
			}
		}
	}
	context.Solution("A", totalPresses)
	return nil
}
