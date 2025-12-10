package solutions

import (
	"container/heap"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
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
	target        int
	buttonMasks   []int
	buttonIndices [][]int
	joltages      []int
}

type queue[T any] struct {
	fifo []T
	lifo []T
}

func (q *queue[T]) push(x T) {
	q.lifo = append(q.lifo, x)
}
func (q *queue[T]) pop() (T, bool) {
	if len(q.fifo) > 0 {
		n := len(q.fifo)
		x := q.fifo[n-1]
		q.fifo = q.fifo[:n-1]
		return x, true
	}
	if len(q.lifo) == 0 {
		var zero T
		return zero, false
	}
	q.fifo = q.lifo
	slices.Reverse(q.fifo)
	q.lifo = nil
	return q.pop()
}

func encode(xs []int) string {
	strs := make([]string, 0)
	for _, x := range xs {
		strs = append(strs, strconv.FormatInt(int64(x), 10))
	}
	return strings.Join(strs, ",")
}

type day10b struct {
	heap           d10bHeap
	visited        map[string]int
	buttonIndices  [][]int
	targetDistance int
	widestButton   int
}

type d10bState struct {
	joltagesLeft []int
	pressesMade  int
	lowerBound   int
}

type d10bHeap []d10bState

func (h d10bHeap) Len() int {
	return len(h)
}

func (h d10bHeap) Less(i, j int) bool {
	if h[i].lowerBound != h[j].lowerBound {
		return h[i].lowerBound < h[j].lowerBound
	}
	if h[i].pressesMade != h[j].pressesMade {
		return h[i].pressesMade > h[j].pressesMade
	}
	return sum(h[i].joltagesLeft) < sum(h[j].joltagesLeft)
}

func sum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

func (h d10bHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *d10bHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(d10bState))
}

func (h *d10bHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (self *day10b) extractSolution() (int, bool) {

	if self.targetDistance > 0 {
		return self.targetDistance, true
	}
	return 0, false
}
func (self *day10b) init(originalTarget []int, buttonIndices [][]int) {
	widestButton := 0
	for _, button := range buttonIndices {
		widestButton = max(widestButton, len(button))
	}
	self.widestButton = widestButton
	self.buttonIndices = buttonIndices
	self.visited = make(map[string]int)
	self.visited[encode(originalTarget)] = 0
	self.heap = []d10bState{d10bState{
		joltagesLeft: slices.Clone(originalTarget),
		pressesMade:  0,
		lowerBound:   self.computeLowerBound(originalTarget),
	}}
	slog.Debug("init", "widestButton", widestButton)
}

func (self *day10b) computeLowerBound(joltagesLeft []int) int {
	buttonIndices := slices.Clone(self.buttonIndices)
	estimates := make([]int, 16)
	for ei := range estimates {
		jl := slices.Clone(joltagesLeft)
		for {
			startingMax := slices.Max(jl)
			rand.Shuffle(len(buttonIndices), func(i, j int) {
				buttonIndices[i], buttonIndices[j] = buttonIndices[j], buttonIndices[i]
			})
			for _, button := range buttonIndices {
				n := slices.Max(jl)
				for _, bi := range button {
					n = min(n, jl[bi])
				}
				for _, bi := range button {
					jl[bi] -= n
				}
				estimates[ei] += n
			}
			if startingMax == slices.Max(jl) {
				break
			}
		}
		estimates[ei] += slices.Max(jl)

	}
	// slog.Debug("computeLowerBound", "joltagesLeft", joltagesLeft, "estimates", estimates)
	return slices.Min(estimates)
}

func (self *day10b) step() error {
	if self.heap.Len() == 0 {
		return fmt.Errorf("empty heap")
	}
	current := heap.Pop(&self.heap).(d10bState)
	slog.Debug("step", "current", current)
	encoded := encode(current.joltagesLeft)
	otherPressesMade, ok := self.visited[encoded]
	if !ok {
		return fmt.Errorf("unexpected not visited")
	}
	if otherPressesMade < current.pressesMade {
		return nil
	}
	for _, button := range self.buttonIndices {
		nextLeft := slices.Clone(current.joltagesLeft)
		pressesMade := current.pressesMade + 1
		for _, buttonindex := range button {
			nextLeft[buttonindex] -= 1
		}
		if slices.Min(nextLeft) < 0 {
			continue
		}
		if slices.Max(nextLeft) == 0 {
			self.targetDistance = pressesMade
			return nil
		}
		encoded = encode(nextLeft)
		otherPressesMade, ok = self.visited[encoded]
		if ok && otherPressesMade <= pressesMade {
			continue
		}
		self.visited[encoded] = pressesMade

		heap.Push(&self.heap, d10bState{
			joltagesLeft: nextLeft,
			pressesMade:  pressesMade,
			lowerBound:   pressesMade + self.computeLowerBound(nextLeft),
		})
	}
	return nil
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
		buttonsIndices := make([][]int, 0)
		for _, button := range elements[1 : len(elements)-1] {
			if button[0] != '(' || button[len(button)-1] != ')' {
				return fmt.Errorf("button not wrapped in () %q", button)
			}
			buttonNs := strings.Split(button[1:len(button)-1], ",")
			buttonMask := 0
			buttonIndices := make([]int, 0)
			for _, buttonN := range buttonNs {
				n, err := strconv.Atoi(buttonN)
				if err != nil {
					return err
				}
				buttonMask |= (1 << n)
				buttonIndices = append(buttonIndices, n)
			}
			buttonMasks = append(buttonMasks, int(buttonMask))
			buttonsIndices = append(buttonsIndices, buttonIndices)
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
		machines = append(machines, machine10{target: int(targetMask), buttonMasks: buttonMasks, joltages: joltages, buttonIndices: buttonsIndices})
	}
	slog.Debug("input", "machines", machines)
	totalPresses := 0
machineLoop1:
	for mi, machine := range machines {
		q := queue[[2]int]{}
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
					continue machineLoop1
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
	totalPresses = 0
	for mi, machine := range machines {
		state := day10b{}
		slog.Debug("starting machine", "machine", mi)
		state.init(machine.joltages, machine.buttonIndices)
		for {
			currentPressees, done := state.extractSolution()
			if done {
				totalPresses += currentPressees
				slog.Debug("solution", "currentPresses", currentPressees)
				break
			}
			err = state.step()
			if err != nil {
				return err
			}
		}
	}
	context.Solution("B", totalPresses)

	return nil
}
