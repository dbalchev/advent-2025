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
	maze := strings.Split(strings.Trim(string(inputBytes), "\n"), "\n")
	nRows := len(maze)
	nCols := len(maze[0])
	for i, row := range maze {
		if len(row) != nCols {
			return fmt.Errorf("row %d has %d cols instead of %d", i, len(row), nCols)
		}
	}
	slog.Info("Maze Size", "nRows", nRows, "nCols", nCols)
	totalCount := 0
	for i := range nRows {
		for j := range nCols {
			if maze[i][j] != '@' {
				continue
			}
			aiMin := i - 1
			if aiMin < 0 {
				aiMin = 0
			}
			aiMax := i + 2
			if aiMax > nRows {
				aiMax = nRows
			}
			ajMin := j - 1
			if ajMin < 0 {
				ajMin = 0
			}
			ajMax := j + 2
			if ajMax > nCols {
				ajMax = nCols
			}
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
				totalCount += 1
			}
		}
	}
	context.Solution("A", totalCount) // 9 min
	return nil
}
