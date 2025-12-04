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
	aoclibrary.Register("day04", &day04{})
}

type day04 struct{}

func (*day04) Solve(context *aoclibrary.Context) error {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	mazeStr := strings.Split(strings.Trim(string(inputBytes), "\n"), "\n")
	maze := make([][]byte, 0)
	for _, row := range mazeStr {
		maze = append(maze, []byte(row))
	}
	nRows := len(maze)
	nCols := len(maze[0])
	for i, row := range maze {
		if len(row) != nCols {
			return fmt.Errorf("row %d has %d cols instead of %d", i, len(row), nCols)
		}
	}
	slog.Info("Maze Size", "nRows", nRows, "nCols", nCols)
	solutionA := 0
	solutionB := 0
	for {
		removable := findRemovable(nRows, nCols, maze)
		if solutionB == 0 {
			solutionA = len(removable)
		}
		if len(removable) == 0 {
			break
		}
		slog.Debug("Removable step", "len(removable)", len(removable))
		solutionB += len(removable)
		for _, p := range removable {
			maze[p[0]][p[1]] = '.'
		}
	}
	context.Solution("A", solutionA) // 9 min
	context.Solution("B", solutionB) // 6 min
	return nil
}

func findRemovable(nRows int, nCols int, maze [][]byte) [][2]int {
	removable := make([][2]int, 0)
	for i := range nRows {
		for j := range nCols {
			if maze[i][j] != '@' {
				continue
			}
			aiMin := max(i-1, 0)
			aiMax := min(i+2, nRows)
			ajMin := max(j-1, 0)
			ajMax := min(j+2, nCols)
			adjRolls := 0
			for ai := aiMin; ai < aiMax; ai += 1 {
				for aj := ajMin; aj < ajMax; aj += 1 {
					if ai == i && aj == j {
						continue
					}
					if maze[ai][aj] == '@' {
						adjRolls += 1
					}
				}
			}
			if adjRolls < 4 {
				removable = append(removable, [2]int{i, j})
			}
		}
	}
	return removable
}
