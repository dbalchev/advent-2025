package solutions

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"
	"unicode"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day06", &day06{})
}

type day06 struct{}

// A - 9:33-9:47
// B - 9:51-10:09, 12:00
func (*day06) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	lines := bytes.Split(bytes.Trim(input, "\n"), []byte("\n"))
	numbersA, opsA, err := parse06A(lines)
	if err != nil {
		return err
	}
	resultA, err := processMath06(opsA, numbersA)
	if err != nil {
		return err
	}
	context.Solution("A", resultA)
	numbersB, opsB, err := parse06B(lines)
	if err != nil {
		return err
	}
	resultB, err := processMath06(opsB, numbersB)
	if err != nil {
		return err
	}
	context.Solution("B", resultB)
	return nil
}

func transpose(m [][]*big.Int) [][]*big.Int {
	r := make([][]*big.Int, 0)
	nElements := len(m[0])
	for i := range m[0] {
		newRow := make([]*big.Int, 0)
		for _, oldRow := range m {
			if len(oldRow) != nElements {
				panic("not rectangular matrix")
			}
			newRow = append(newRow, oldRow[i])
		}
		r = append(r, newRow)
	}
	return r
}

func parse06B(lines [][]byte) ([][]*big.Int, []string, error) {
	opColumns := make([]int, 0)
	lastLine := lines[len(lines)-1]
	for i, c := range lastLine {
		if !unicode.IsSpace(rune(c)) {
			opColumns = append(opColumns, i)
		}
	}
	opColumns = append(opColumns, len(lastLine)+1)
	slog.Debug("opColumns", "", opColumns)
	for opColIndex, opColumn := range opColumns[1 : len(opColumns)-1] {
		for lineIndex, line := range lines {
			if !unicode.IsSpace(rune(line[opColumn-1])) {
				return nil, nil, fmt.Errorf("unexpected non-space on %d:%d", lineIndex, opColIndex)
			}

		}
	}
	ops := make([]string, 0)
	numbers := make([][]*big.Int, 0)

	numersLens := make([]int, 0)
	for i, columnStart := range opColumns[:len(opColumns)-1] {
		columnEnd := opColumns[i+1] - 1
		ops = append(ops, string(lastLine[columnStart:columnStart+1]))
		colNumbers := make([]*big.Int, 0)
		for j := columnEnd - 1; j >= columnStart; j -= 1 {
			numStr := make([]byte, 0)
			for _, line := range lines[0 : len(lines)-1] {
				c := line[j]
				numStr = append(numStr, c)
			}
			if len(numStr) == 0 {
				return nil, nil, fmt.Errorf("unexpected empty number")
			}
			numStr = bytes.Trim(numStr, " ")
			number := big.Int{}
			_, ok := number.SetString(string(numStr), 10)
			if !ok {
				return nil, nil, fmt.Errorf("cannot parse %q", string(numStr))
			}
			colNumbers = append(colNumbers, &number)
		}
		numbers = append(numbers, colNumbers)
		numersLens = append(numersLens, len(colNumbers))
	}
	// slog.Debug("numbersLen", "numbersLen", numersLens)
	return numbers, ops, nil
}

func parse06A(lines [][]byte) ([][]*big.Int, []string, error) {
	numbers := make([][]*big.Int, 0)
	for _, line := range lines[:len(lines)-1] {
		lineNumbers := make([]*big.Int, 0)
		for _, field := range bytes.Fields(line) {
			number := big.Int{}
			_, ok := number.SetString(string(field), 10)
			if !ok {
				return nil, nil, fmt.Errorf("cannot parse %q", string(field))
			}
			lineNumbers = append(lineNumbers, &number)
		}
		numbers = append(numbers, lineNumbers)
	}
	ops := make([]string, 0)
	for _, field := range bytes.Fields(lines[len(lines)-1]) {
		ops = append(ops, string(field))
	}
	numbers = transpose(numbers)

	return numbers, ops, nil
}

func processMath06(ops []string, numbers [][]*big.Int) (*big.Int, error) {
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
			return nil, fmt.Errorf("cannot parse op %q", op)
		}
		for _, x := range numbers[i] {
			update(x)
		}
		slog.Debug("op result", "result", opAccum)
		accum = accum.Add(accum, opAccum)
	}
	return accum, nil
}
