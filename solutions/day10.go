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
		slog.Debug("starting machine", "mi", mi)
		type precomp struct {
			nPresses int
			joltages []int
		}
		precomps := make([][]precomp, 1<<len(machine.joltages))
		for pressMask := 0; pressMask < (1 << len(machine.buttonIndices)); pressMask += 1 {
			nPresses := 0
			joltages := make([]int, len(machine.joltages))
			for bi, button := range machine.buttonIndices {
				if pressMask&(1<<bi) == 0 {
					continue
				}
				nPresses += 1
				for _, j := range button {
					joltages[j] += 1
				}
			}
			jmask := 0
			for i, j := range joltages {
				if j&1 == 0 {
					continue
				}
				jmask |= (1 << i)
			}
			fi := -1
			for i, p := range precomps[jmask] {
				if slices.Equal(p.joltages, joltages) {
					fi = i
					break
				}
			}
			if fi == -1 {
				precomps[jmask] = append(precomps[jmask], precomp{nPresses: nPresses, joltages: joltages})
			} else {
				precomps[jmask][fi].nPresses = min(precomps[jmask][fi].nPresses, nPresses)
			}
		}
		type intermediate struct {
			nPresses int
			residual []int
		}
		m := 1
		solutions := make([]int, 0)
		ims := []intermediate{{nPresses: 0, residual: slices.Clone(machine.joltages)}}
		for len(ims) != 0 {
			newIms := make([]intermediate, 0)
			for _, im := range ims {
				rm := 0
				for i, r := range im.residual {
					if r&1 == 1 {
						rm |= (1 << i)
					}
				}
				for _, p := range precomps[rm] {
					newResidual := make([]int, len(im.residual))
					for i := range newResidual {
						r := im.residual[i] - p.joltages[i]
						if r&1 == 1 {
							return fmt.Errorf("bit is set?")
						}
						newResidual[i] = r >> 1
					}
					if slices.Min(newResidual) < 0 {
						continue
					}
					nPresses := im.nPresses + p.nPresses*m
					if slices.Max(newResidual) == 0 {
						solutions = append(solutions, nPresses)
						continue
					}

					newIms = append(newIms, intermediate{residual: newResidual, nPresses: nPresses})
				}
			}
			ims = newIms
			m *= 2
		}
		slog.Debug("done", "solutions", solutions)
		totalPresses += slices.Min(solutions)
		// for jm, p := range precomps {
		// 	slog.Debug("pc", "jm", jm, "p", p)
		// }
	}
	context.Solution("B", totalPresses)

	return nil
}
