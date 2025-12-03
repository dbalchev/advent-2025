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

	sumA := 0
	var sumB int64 = 0
	for _, line := range lines {
		solutionA := solveLine03(line, 2)
		solutionB := solveLine03(line, 12)
		slog.Debug("best number", "solutionA", solutionA, "solutionB", solutionB, "line", line)
		numA, err := strconv.Atoi(solutionA)
		if err != nil {
			return err
		}
		numB, err := strconv.ParseInt(solutionB, 10, 64)
		if err != nil {
			return err
		}
		sumA += numA
		sumB += numB
	}
	context.Solution("A", sumA)
	context.Solution("B", sumB)
	return nil
}

func solveLine03(line string, nBateries int) string {
	resultBytes := make([]byte, 0)
	lineBytes := []byte(line)
	startOffset := 0
	for batteriesLeft := nBateries; batteriesLeft > 0; batteriesLeft -= 1 {
		currentIndex := findLargestIndex(lineBytes[startOffset : len(line)-batteriesLeft+1])
		resultBytes = append(resultBytes, lineBytes[startOffset+currentIndex])
		startOffset += currentIndex + 1
	}
	return string(resultBytes)
}

func findLargestIndex(lineBytes []byte) int {
	largestIndex := 0
	for i, c := range lineBytes {
		if c > lineBytes[largestIndex] {
			largestIndex = i
		}
	}
	return largestIndex
}
