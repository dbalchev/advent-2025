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
	originalTarget       []int
	targetKernelAddition map[string]int
	targetDidtance       int
	q                    queue[[]int]
	visited              map[string]int
	buttonIndices        [][]int
}

func (self *day10b) extractSolution() (int, bool) {

	if self.targetDidtance > 0 {
		return self.targetDidtance, true
	}
	return 0, false
}
func (self *day10b) init(originalTarget []int, buttonIndices [][]int) {
	self.originalTarget = originalTarget
	self.targetKernelAddition = make(map[string]int)
	self.targetKernelAddition[encode(originalTarget)] = 0
	self.buttonIndices = buttonIndices
	self.visited = make(map[string]int)
	initialState := make([]int, len(originalTarget), len(originalTarget)+1)
	self.visited[encode(initialState)] = 0
	initialState = append(initialState, 0)
	initialState = append(initialState, 0)
	self.q.push(initialState)
}
func (self *day10b) step() error {
	x, ok := self.q.pop()
	if !ok {
		return fmt.Errorf("empty queue")
	}
	slog.Debug("", "x", x)
	firstButton := x[len(x)-1]
	for bi, button := range self.buttonIndices[firstButton:] {
		next := slices.Clone(x)
		nextState := next[:len(next)-2]
		next[len(next)-2] += 1
		next[len(next)-1] = bi + firstButton
		d := next[len(next)-2]
		for _, buttonIndex := range button {
			nextState[buttonIndex] += 1
		}
		encoded := encode(nextState)
		additioalDistance, reachable := self.targetKernelAddition[encoded]
		if reachable {
			slog.Debug("reachable", "additioalDistance", additioalDistance, "encoded", encoded)
			self.targetDidtance = additioalDistance + d
		}
		_, isVisited := self.visited[encoded]
		if isVisited {
			continue
		}
		if slices.Max(nextState) == slices.Min(nextState) && len(self.targetKernelAddition) == 1 {
			slog.Debug("kernel found", "kernel", nextState, "d", d, "self.targetKernelAddition", self.targetKernelAddition)
			mk := slices.Clone(nextState)
			residual := make([]int, len(mk))
			kernelMul := 1
		mkLoop:
			for {
				slog.Debug("prefilling", "kernelMul", kernelMul)
				for i, v := range mk {
					residual[i] = self.originalTarget[i] - v
					if residual[i] < 0 {
						break mkLoop
					}
				}
				encodedResidual := encode(residual)
				rd, reachable := self.visited[encodedResidual]
				if reachable {
					slog.Debug("residual reachable", "rd", rd, "residual", residual)
					self.targetDidtance = rd + d*kernelMul
				}
				self.targetKernelAddition[encodedResidual] = d * kernelMul
				kernelMul += 1
				for i, v := range nextState {
					mk[i] += v
				}
			}
		}
		self.visited[encoded] = d
		self.q.push(next)
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
