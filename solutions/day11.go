package solutions

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day11", &day11{})
}

type day11 struct{}

func solve11(
	children, parents map[string][]string,
	allNodes map[string]bool,
	startNode, endNode string,
	intermediateNodes []string,
) (int, error) {
	type state struct {
		node        string
		imNodesMask int
	}
	pathCount := make(map[state]int)
	pathCount[state{node: endNode, imNodesMask: 0}] = 1
	stack := []state{{node: endNode, imNodesMask: 0}}
	nRemainingChildren := make(map[string]int)
	for node := range allNodes {
		nRemainingChildren[node] = len(children[node])
	}
	maxNodeMask := make(map[string]int)
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, parent := range parents[x.node] {
			newMaks := x.imNodesMask
			for i, n := range intermediateNodes {
				if n == parent {
					newMaks |= 1 << i
					slog.Debug("visiting", "n", n)
				}
			}
			newState := state{node: parent, imNodesMask: newMaks}
			pathCount[newState] += pathCount[x]
			slog.Debug("adding", "newState", newState, "x", x, "pathCount[x]", pathCount[x])
			nRemainingChildren[parent] -= 1
			maxNodeMask[parent] = maxNodeMask[parent] | newMaks
			if nRemainingChildren[parent] == 0 {
				stack = append(stack, state{node: parent, imNodesMask: maxNodeMask[parent]})
			}
		}
	}
	solution, ok := pathCount[state{node: startNode, imNodesMask: (1 << len(intermediateNodes)) - 1}]
	if !ok {
		slog.Debug("no counts", "pathCount", pathCount, "nRemainingChildren", nRemainingChildren)
		return 0, fmt.Errorf("no %q counts", startNode)
	}
	return solution, nil
}

func (*day11) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.Trim(string(input), "\n"), "\n")
	children := make(map[string][]string)
	parents := make(map[string][]string)
	allNodes := make(map[string]bool)
	for _, line := range lines {
		fields := strings.Fields(line)
		parent := fields[0]
		if parent[len(parent)-1] != ':' {
			return fmt.Errorf("missing ':'")
		}
		parent = parent[:len(parent)-1]
		allNodes[parent] = true
		for _, child := range fields[1:] {
			allNodes[child] = true
			children[parent] = append(children[parent], child)
			parents[child] = append(parents[child], parent)
		}
	}
	slog.Debug("input", "children", children, "parents", parents, "allNodes", allNodes)
	solutionA, err := solve11(children, parents, allNodes, "you", "out", []string{})
	if err != nil {
		return err
	}
	context.Solution("A", solutionA)
	solutionB, err := solve11(children, parents, allNodes, "svr", "out", []string{"fft", "dac"})
	if err != nil {
		return err
	}
	context.Solution("B", solutionB)
	return nil
}
