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
func (*day07) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := bytes.Split(bytes.Trim(input, "\n"), []byte{'\n'})
	sourceIndex := bytes.Index(lines[0], []byte{'S'})
	tIndexes := make(map[int]bool)
	tIndexes[sourceIndex] = true
	nSplits := 0
	for _, currentLine := range lines[1:] {
		newIndexes := make(map[int]bool)
		for index := range tIndexes {
			if currentLine[index] == '^' {
				newIndexes[index-1] = true
				newIndexes[index+1] = true
				nSplits += 1
			} else {
				newIndexes[index] = true
			}
		}
		tIndexes = newIndexes
	}
	context.Solution("A", nSplits)
	return nil
}
