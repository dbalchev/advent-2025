package solutions

import (
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sync"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day09", &day09{})
}

type day09 struct{}

func (*day09) Solve(context *aoclibrary.Context) error {
	nThreads := 8
	inputTiles := make([][]int, 0)
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
	}
	candidates := make([][]int, 0)
	for i := range xs {
		for dj := range xs[i+1:] {
			j := i + dj + 1
			dx := aoclibrary.Iabs(xs[i] - xs[j])
			dy := aoclibrary.Iabs(ys[i] - ys[j])
			area := (dx + 1) * (dy + 1)
			candidates = append(candidates, []int{area, i, j})
		}
	}
	slices.SortFunc(candidates, slices.Compare)
	slices.Reverse(candidates)
	candidateChannel := make(chan []int, 64)
	solutionDiscovered := make(chan any, 1)
	clsoeSolutionDiscovered := sync.OnceFunc(func() {
		close(solutionDiscovered)
	})
	solutions := make(chan int, nThreads+1)
	wg := sync.WaitGroup{}
	go func() {
		for _, candcandidate := range candidates {
			select {
			case candidateChannel <- candcandidate:
			case _, ok := <-solutionDiscovered:
				slog.Info("exiting the candidates filling channel", "ok", ok)
				return
			}
		}
	}()
	for threadIndex := range nThreads {
		wg.Go(func() {
		candidatesLoop:
			for {
				select {
				case _, ok := <-solutionDiscovered:
					slog.Info("exiting worker", "workerId", threadIndex, "ok", ok)
					return
				case candidate, ok := <-candidateChannel:
					if !ok {
						slog.Info("ran out of candidates", "workerId", threadIndex)
						return
					}
					area := candidate[0]
					i := candidate[1]
					j := candidate[2]
					cmx := []int{xs[i], xs[j]}
					cmy := []int{ys[i], ys[j]}
					slices.Sort(cmx[:])
					slices.Sort(cmy[:])
					for csi := range xs {
						nsi := (csi + 1) % len(xs)
						esxs := []int{xs[csi], xs[nsi]}
						esys := []int{ys[csi], ys[nsi]}
						slices.Sort(esxs[:])
						slices.Sort(esys[:])
						// the segment is incident to a candidate corners
						if i == csi || i == nsi || j == csi || j == nsi {
							continue
						}
						// the shifted polygon segment doesn't touch the candidate horizontally
						if cmx[0] >= esxs[1] || cmx[1] <= esxs[0] {
							continue
						}
						// the shifted polygon segment doesn't touch the candidate vertically
						if cmy[0] >= esys[1] || cmy[1] <= esys[0] {
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
						continue candidatesLoop
					}
					clsoeSolutionDiscovered()
					solutions <- area
				}
			}
		})
	}
	wg.Wait()
	close(solutions)
	solutionsSlice := make([]int, 0)
	for solution := range solutions {
		solutionsSlice = append(solutionsSlice, solution)
	}
	slog.Debug("solutionSlice", "solutionSlice", solutionsSlice)
	context.Solution("A", candidates[0][0])
	context.Solution("B", slices.Max(solutionsSlice))
	return nil
}
