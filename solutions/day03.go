package solutions

import (
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day03", &day03{})
}

type day03 struct{}

func (*day03) Solve(context *aoclibrary.Context) error {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.Trim(string(inputBytes[:]), "\n"), "\n")
	slog.Info("Size stats", "nlines", len(lines), "fistrLineLen", len(lines[0]))

	sum := 0
	for _, line := range lines {
		firstIndex := findLargestIndex(line[:len(line)-1])
		secondIndex := findLargestIndex(line[firstIndex+1:])
		bestNumberStr := string([]byte{line[firstIndex], line[firstIndex+secondIndex+1]})
		bestNumber, err := strconv.Atoi(bestNumberStr)
		if err != nil {
			return err
		}
		sum += bestNumber
		slog.Debug("best number", "num", bestNumberStr, "line", line)
	}
	context.Solution("A", sum)
	return nil
}

func findLargestIndex(line string) int {
	largestIndex := 0
	for i, c := range line {
		if c > rune(line[largestIndex]) {
			largestIndex = i
		}
	}
	return largestIndex
}
