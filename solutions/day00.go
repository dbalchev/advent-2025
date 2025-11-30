package solutions

import aoclibrary "github.com/dbalchev/advent-2025/aoc-library"

func init() {
	aoclibrary.Register(0, &Day00{})
}

type Day00 struct{}

func (*Day00) Solve(context *aoclibrary.Context) error {
	context.Solution(1, "foo")
	return nil
}
