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
	pathCount := make(map[string]int)
	pathCount["out"] = 1
	stack := []string{"out"}
	nRemainingChildren := make(map[string]int)
	for node := range allNodes {
		nRemainingChildren[node] = len(children[node])
	}
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, parent := range parents[x] {
			pathCount[parent] += pathCount[x]
			nRemainingChildren[parent] -= 1
			if nRemainingChildren[parent] == 0 {
				stack = append(stack, parent)
			}
		}
	}
	solutionA, ok := pathCount["you"]
	if !ok {
		return fmt.Errorf("no you counts")
	}
	context.Solution("A", solutionA)
	return nil
}
