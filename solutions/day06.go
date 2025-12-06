package solutions

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day06", &day06{})
}

type day06 struct{}

// A - 9:33-9:47
func (*day06) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := bytes.Split(bytes.Trim(input, "\n"), []byte("\n"))
	numbers := make([][]*big.Int, 0)
	for _, line := range lines[:len(lines)-1] {
		lineNumbers := make([]*big.Int, 0)
		for _, field := range bytes.Fields(line) {
			number := big.Int{}
			_, ok := number.SetString(string(field), 10)
			if !ok {
				return fmt.Errorf("Cannot parse %q", string(field))
			}
			lineNumbers = append(lineNumbers, &number)
		}
		numbers = append(numbers, lineNumbers)
	}
	ops := make([]string, 0)
	for _, field := range bytes.Fields(lines[len(lines)-1]) {
		ops = append(ops, string(field))
	}
	slog.Debug("ops", "ops", ops)
	slog.Debug("numbers", "numbers", numbers)
	accum := big.NewInt(0)
	for i, op := range ops {
		var opAccum *big.Int
		var update func(*big.Int)
		switch op {
		case "+":
			opAccum = big.NewInt(0)
			update = func(x *big.Int) {
				opAccum = opAccum.Add(opAccum, x)
			}
		case "*":
			opAccum = big.NewInt(1)
			update = func(x *big.Int) {
				opAccum = opAccum.Mul(opAccum, x)
			}
		default:
			return fmt.Errorf("cannot parse op %q", op)
		}
		for _, line := range numbers {
			update(line[i])
		}
		slog.Debug("op result", "result", opAccum)
		accum = accum.Add(accum, opAccum)
	}
	context.Solution("A", accum)
	return nil
}
