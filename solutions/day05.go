package solutions

import (
	"fmt"
	"io"
	"log/slog"
	"slices"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day05", &day05{})
}

type day05 struct{}

// A 8:48 - 8:56
// B 8:57 - 9:06
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
	context.Solution("B", solve05B(freshRanges))
	return nil
}

func solve05B(freshRanges [][2]int) int {
	events := make([][]int, 0)
	for _, fr := range freshRanges {
		events = append(events, []int{fr[0], -1})
		events = append(events, []int{fr[1], +1})
	}
	slices.SortFunc(events, slices.Compare)
	freshStart := -1
	acc := 0
	freshLen := 0
	for _, event := range events {
		increment := -event[1]
		if acc == 0 && increment < 0 {
			panic("unexpected decrement")
		}
		if acc == 0 {
			freshStart = event[0]
		}
		acc += increment
		if acc == 0 {
			freshLen += event[0] - freshStart + 1
		}
		if acc < 0 {
			panic("acc is negative")
		}
	}
	if acc != 0 {
		panic("acc is 0 at end")
	}
	return freshLen
}
