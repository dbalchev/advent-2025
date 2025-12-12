package solutions

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day12", &day12{})
}

type day12 struct{}

func (*day12) Solve(context *aoclibrary.Context) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	chunks := bytes.Split(bytes.Trim(input, "\n"), []byte{'\n', '\n'})
	presentUnits := make([]int, 0)
	for i, present := range chunks[:len(chunks)-1] {
		presentLines := strings.Split(string(present), "\n")
		if presentLines[0] != fmt.Sprintf("%d:", i) {
			return fmt.Errorf("unexpected first present line")
		}
		counter := make(map[rune]int)
		for _, line := range presentLines[1:] {
			for _, r := range line {
				counter[r] += 1
			}
		}
		if counter['.']+counter['#'] != 9 {
			return fmt.Errorf("unexpected present size")
		}
		presentUnits = append(presentUnits, counter['#'])
	}
	slog.Debug("presentUnits", "presentUnits", presentUnits)
	sureFalseRegions, sureTrueRegions, undecidedRegions, withEnoughArea := 0, 0, 0, 0
	for _, region := range strings.Split(string(chunks[len(chunks)-1]), "\n") {
		fields := strings.Fields(region)
		presentCounts := make([]int, 0)
		nPresents := 0
		totalPresentUnits := 0
		for i, f := range fields[1:] {
			n, err := strconv.Atoi(f)
			if err != nil {
				return err
			}
			presentCounts = append(presentCounts, n)
			nPresents += n
			totalPresentUnits += presentUnits[i] * n
		}
		if len(presentCounts) != len(presentUnits) {
			return fmt.Errorf("present count is different")
		}
		sizesStr := strings.Split(strings.Trim(fields[0], ":"), "x")
		if len(sizesStr) != 2 {
			return fmt.Errorf("wrong size len")
		}
		w, err := strconv.Atoi(sizesStr[0])
		if err != nil {
			return err
		}
		h, err := strconv.Atoi(sizesStr[1])
		if err != nil {
			return err
		}
		slog.Debug("region", "w", w, "h", h, "presentCounts", presentCounts, "nPresents", nPresents, "totalPresentUnits", totalPresentUnits)
		if totalPresentUnits > w*h {
			sureFalseRegions += 1
		} else if ((w+2)/3)*((h+2)/3) <= nPresents {
			sureTrueRegions += 1
		} else if w*h >= totalPresentUnits {
			withEnoughArea += 1
		} else {
			undecidedRegions += 1
		}
	}
	slog.Debug("heuristics", "sureFalseRegions", sureFalseRegions, "sureTrueRegions", sureTrueRegions, "undecidedRegions", undecidedRegions, "withEnoughArea", withEnoughArea)
	return nil
}
