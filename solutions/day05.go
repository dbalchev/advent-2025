package solutions

import (
	"fmt"
	"io"
	"log/slog"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day05", &day05{})
}

type day05 struct{}

// A 8:48 - 8:56
func (*day05) Solve(context *aoclibrary.Context) error {
	freshRanges := make([][2]int, 0)
	ingredients := make([]int, 0)
	for {
		s, e := 0, 0
		n, err := fmt.Scanf("%d-%d\n", &s, &e)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		if n != 2 {
			return fmt.Errorf("scanf n=%d", n)
		}
		freshRanges = append(freshRanges, [2]int{s, e})
	}
	for {
		x := 0
		n, err := fmt.Scanf("%d\n", &x)
		if err != nil && err != io.EOF {
			return err
		}
		if n != 1 {
			break
		}
		ingredients = append(ingredients, x)
	}
	slog.Debug("input", "freshRanges", freshRanges, "ingredients", ingredients)
	freshCount := 0
	for _, ingredient := range ingredients {
		for _, freshRange := range freshRanges {
			if freshRange[0] <= ingredient && ingredient <= freshRange[1] {
				freshCount += 1
				break
			}
		}
	}
	context.Solution("A", freshCount)
	return nil
}
