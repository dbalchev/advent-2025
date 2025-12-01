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
	context.Solution("A", solve01A(instructions))
	context.Solution("B", solve02A(instructions))
	return nil
}

func solve01A(instructions []instruction) int {
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
	return zeroCounts
}

func moddiv(a, n int) (int, int) {
	r := a % n
	d := a / n
	if r < 0 {
		r += 100
		d -= 1
	}
	if d*n+r != a {
		panic(fmt.Sprintf("sanitycheck failed %d * %d + %d != %d", d, n, r, a))
	}
	return r, d
}

func solve02A(instructions []instruction) int {
	current := 50
	clickCounts := 0

	for _, instruction := range instructions {
		increaseCount := func() {
			clickCounts += 1
			// fmt.Printf("Click after %c %d\n", instruction.direction, instruction.count)
		}
		rots := instruction.count
		if instruction.direction == 'L' {
			rots = -instruction.count
		}
		shouldCount := current != 0
		current += rots
		for current < 0 {
			if shouldCount {
				increaseCount()
			} else {
				shouldCount = true
			}
			current += 100
		}
		if current == 0 {
			increaseCount()
		}
		for current >= 100 {
			increaseCount()
			current -= 100
		}
	}
	return clickCounts
}
