package solutions

import (
	"fmt"
	"io"
	"slices"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day00", &day00{})
}

type day00 struct{}

func (*day00) Solve(context *aoclibrary.Context) error {
	lefts := make([]int, 0)
	rights := make([]int, 0)
	for {
		left := -1
		right := -1
		n, err := fmt.Scanf("%d %d\n", &left, &right)
		if err != nil && err != io.EOF {
			return err
		}
		if n != 2 {
			break
		}
		lefts = append(lefts, left)
		rights = append(rights, right)

	}
	fmt.Printf("lefts = %v; rights = %v\n", lefts, rights)
	context.Solution(1, solve1(slices.Clone(lefts), slices.Clone(rights)))

	sum := solve2(rights, lefts)
	context.Solution(2, sum)

	return nil
}

func iabs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func solve1(lefts []int, rights []int) int {
	slices.Sort(lefts)
	slices.Sort(rights)
	sum := 0
	for i := range lefts {
		sum += int(iabs(lefts[i] - rights[i]))
	}
	return sum
}

func solve2(rights []int, lefts []int) int {
	rightCounts := make(map[int]int)
	for _, x := range rights {
		c, exists := rightCounts[x]
		if exists {
			rightCounts[x] = c + 1
		} else {
			rightCounts[x] = 1
		}
	}
	sum := 0
	for _, x := range lefts {
		sum += rightCounts[x] * x
	}
	return sum
}
