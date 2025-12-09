package solutions

import (
	"fmt"
	"io"
	"log/slog"
	"slices"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day09", &day09{})
}

type day09 struct{}

func (*day09) Solve(context *aoclibrary.Context) error {
	inputTiles := make([][]int, 0)
	m := 16
	for {
		tile := [2]int{}
		n, err := fmt.Scanf("%d,%d\n", &tile[0], &tile[1])
		if err != nil && err != io.EOF {
			return err
		}
		if n != 2 {
			break
		}
		inputTiles = append(inputTiles, tile[:])
	}
	xs := make([]int, 0)
	ys := make([]int, 0)
	exs := make([]int, 0)
	eys := make([]int, 0)
	mxs := make([]int, 0)
	mys := make([]int, 0)
	for i, inputTile := range inputTiles {
		nextIndex := (i + 1) % len(inputTiles)
		prevTile := inputTiles[(i+len(inputTiles)-1)%len(inputTiles)]
		nextTile := inputTiles[nextIndex]
		dx := prevTile[0] + nextTile[0] - 2*inputTile[0]
		dy := prevTile[1] + nextTile[1] - 2*inputTile[1]
		if dx == 0 || dy == 0 {
			return fmt.Errorf("unexpected dx or dy == 0")
		}
		if dx < 0 {
			dx = -1
		} else {
			dx = 1
		}
		if dy < 0 {
			dy = -1
		} else {
			dy = 1
		}
		xs = append(xs, inputTile[0])
		ys = append(ys, inputTile[1])
		// Scaled coordinates, so we can errode without floating point arithmetics
		mxs = append(mxs, inputTile[0]*m)
		mys = append(mys, inputTile[1]*m)
		// Eroded polygon coordinates, to make easier compares
		exs = append(exs, inputTile[0]*m+dx)
		eys = append(eys, inputTile[1]*m+dy)
	}
	maxArea := 0
	maxBArea := 0
	for i := range xs {
	jLoop:
		for dj := range xs[i+1:] {
			j := i + dj + 1
			dx := aoclibrary.Iabs(xs[i] - xs[j])
			dy := aoclibrary.Iabs(ys[i] - ys[j])
			area := (dx + 1) * (dy + 1)
			maxArea = max(maxArea, area)
			cmx := []int{mxs[i], mxs[j]}
			cmy := []int{mys[i], mys[j]}
			slices.Sort(cmx[:])
			slices.Sort(cmy[:])
			if cmx[0] > cmx[1] {
				panic("foo")
			}
			for csi := range xs {
				nsi := (csi + 1) % len(xs)
				esxs := []int{exs[csi], exs[nsi]}
				esys := []int{eys[csi], eys[nsi]}
				slices.Sort(esxs[:])
				slices.Sort(esys[:])
				// the segment is incident to a candidate corners
				if i == csi || i == nsi || j == csi || j == nsi {
					continue
				}
				// the errorded polygon segment doesn't touch the candidate horizontally
				if cmx[0] > esxs[1] || cmx[1] < esxs[0] {
					continue
				}
				// the errorded polygon segment doesn't touch the candidate vertically
				if cmy[0] > esys[1] || cmy[1] < esys[0] {
					continue
				}
				slog.Debug(
					"skipping",
					"area", area,
					"tile1", []int{xs[i], ys[i]},
					"tile2", []int{xs[j], ys[j]},
					"seg1", []int{xs[csi], ys[csi]},
					"seg2", []int{xs[nsi], ys[nsi]},
					"cmx", cmx,
					"cmy", cmy,
					"esxs", esxs,
					"esys", esys,
				)
				continue jLoop
			}
			maxBArea = max(maxBArea, area)
		}
	}
	context.Solution("A", maxArea)
	context.Solution("B", maxBArea)
	return nil
}
