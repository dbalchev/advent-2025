package solutions

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day02", &day02{})
}

type day02 struct{}

func (*day02) Solve(context *aoclibrary.Context) error {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil
	}
	strRanges := strings.Split(string(inputBytes[:]), ",")
	// fmt.Printf("input = %v\n", strRanges)
	splitRanges := make([][2]string, 0)
	for _, strRange := range strRanges {
		strRange = strings.Trim(strRange, "\n")
		rangeComponents := strings.Split(strRange, "-")
		if len(rangeComponents) != 2 {
			return fmt.Errorf("unexpected range %v", rangeComponents)
		}
		startStr, endStr := rangeComponents[0], rangeComponents[1]
		splitRanges = append(splitRanges, [2]string{startStr, endStr})
	}
	// fmt.Printf("splitRanges = %v\n", splitRanges)
	processedRanges := make([][3]int, 0)
	appendProcessedRange := func(startStr, endStr string) error {
		if len(startStr) != len(endStr) {
			return fmt.Errorf("%q and %q should have the same len", startStr, endStr)
		}
		startInt, err := strconv.Atoi(startStr)
		if err != nil {
			return err
		}
		endInt, err := strconv.Atoi(endStr)
		if err != nil {
			return err
		}
		processedRanges = append(processedRanges, [3]int{startInt, endInt, len(startStr)})
		return nil
	}
	for _, splitRange := range splitRanges {
		startStr, endStr := splitRange[0], splitRange[1]
		for len(startStr) < len(endStr) {
			err := appendProcessedRange(startStr, strings.Repeat("9", len(startStr)))
			if err != nil {
				return err
			}
			startStr = "1" + strings.Repeat("0", len(startStr))
		}
		err := appendProcessedRange(startStr, endStr)
		if err != nil {
			return err
		}
	}
	// fmt.Printf("splitRanges = %v\n", processedRanges)
	sum := 0
	for _, processedRange := range processedRanges {
		start, end, nDigits := processedRange[0], processedRange[1], processedRange[2]
		if nDigits%2 != 0 {
			continue
		}
		h := int(math.Pow10(nDigits / 2))
		m := 1 + h
		q1 := (start + m - 1) / m
		q2 := (end + m) / m
		minQ := h / 10
		if q1 < minQ {
			q1 = minQ
		}
		if q2 < minQ {
			q2 = minQ
		}
		for i := q1; i < q2; i += 1 {
			// fmt.Printf("s %v\n", i*m)
			sum += i * m
		}
		// fmt.Printf("%v %v %v %v %v\n", start, end, m, q1, q2)
		// sum += m * (q1 + q2) * (q2 - q1) / 2
	}
	context.Solution("A", sum)
	return nil
}
