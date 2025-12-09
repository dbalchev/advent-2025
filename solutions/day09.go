package solutions

import (
	"fmt"
	"io"
	"log/slog"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day09", &day09{})
}

type day09 struct{}

func (*day09) Solve(context *aoclibrary.Context) error {
	tiles := make([][2]int, 0)
	for {
		tile := [2]int{}
		n, err := fmt.Scanf("%d,%d\n", &tile[0], &tile[1])
		if err != nil && err != io.EOF {
			return err
		}
		if n != 2 {
			break
		}
		tiles = append(tiles, tile)
	}
	maxArea := 0
	for i, tile1 := range tiles {
		for _, tile2 := range tiles[i+1:] {
			area := aoclibrary.Iabs(tile1[0]-tile2[0]+1) * aoclibrary.Iabs(tile1[1]-tile2[1]+1)
			maxArea = max(maxArea, area)
		}
	}
	slog.Debug("input", "tiles", tiles)
	context.Solution("A", maxArea)
	return nil
}
