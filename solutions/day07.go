package solutions

import (
	"bytes"
	"io"
	"os"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day07", &day07{})
}

type day07 struct{}

// A - 16:06-16:10
// B - 16:14-16:15
func (*day07) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := bytes.Split(bytes.Trim(input, "\n"), []byte{'\n'})
	sourceIndex := bytes.Index(lines[0], []byte{'S'})
	tIndexes := make(map[int]int64)
	tIndexes[sourceIndex] = 1
	nSplits := 0
	for _, currentLine := range lines[1:] {
		newIndexes := make(map[int]int64)
		for index, nTimelines := range tIndexes {
			if currentLine[index] == '^' {
				newIndexes[index-1] += nTimelines
				newIndexes[index+1] += nTimelines
				nSplits += 1
			} else {
				newIndexes[index] += nTimelines
			}
		}
		tIndexes = newIndexes
	}
	var nTimelines int64 = 0
	for _, cnt := range tIndexes {
		nTimelines += cnt
	}
	context.Solution("A", nSplits)
	context.Solution("B", nTimelines)
	return nil
}
