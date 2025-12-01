package solutions

import (
	"fmt"
	"io"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day01", &day01{})
}

type day01 struct{}

type instruction struct {
	direction rune
	count     int
}

func (*day01) Solve(context *aoclibrary.Context) error {
	instructions := make([]instruction, 0)
	for {
		x := instruction{}
		n, err := fmt.Scanf("%c%d\n", &x.direction, &x.count)
		if n != 2 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		instructions = append(instructions, x)
	}
	// fmt.Println("instructions = %v", instructions)
	current := 50
	zeroCounts := 0
	for _, instruction := range instructions {
		rots := instruction.count
		if instruction.direction == 'L' {
			rots = 100 - instruction.count
		}
		current = (current + rots) % 100
		if current == 0 {
			zeroCounts += 1
		}
	}
	context.Solution("A", zeroCounts)
	return nil
}
