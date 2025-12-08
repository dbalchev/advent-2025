package solutions

import (
	"cmp"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strconv"

	aoclibrary "github.com/dbalchev/advent-2025/aoc-library"
)

func init() {
	aoclibrary.Register("day08", &day08{})
}

type day08 struct{}

type suf struct {
	parent []int
}

func makeSuf(n int) suf {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return suf{parent: parent}
}
func (self *suf) update(newRoot, leaf int) {
	for self.parent[leaf] != leaf {
		nextLeaf := self.parent[leaf]
		self.parent[leaf] = newRoot
		leaf = nextLeaf
	}
	self.parent[leaf] = newRoot
}
func (self *suf) find(query int) int {
	x := query
	for self.parent[x] != x {
		x = self.parent[x]
	}
	self.update(x, query)
	return x
}
func (self *suf) union(lh, rh int) {
	lhParent := self.find(lh)
	self.update(lhParent, rh)
}

func (*day08) Solve(context *aoclibrary.Context) error {
	boxes := make([][3]int, 0)
	for {
		newBox := [3]int{}
		n, err := fmt.Scanf("%d,%d,%d\n", &newBox[0], &newBox[1], &newBox[2])
		if n != 0 && n != 3 {
			return fmt.Errorf("unexpected n=%d", n)
		}
		if err != nil && err != io.EOF {
			return err
		}
		if n == 3 {
			boxes = append(boxes, newBox)
			continue
		}
		break
	}
	type edge struct {
		lh       int
		rh       int
		distance int
	}
	edges := make([]edge, 0)
	for i, lh := range boxes {
		for d, rh := range boxes[i+1:] {
			j := i + d + 1
			delta := []int{lh[0] - rh[0], lh[1] - rh[1], lh[2] - rh[2]}
			distance := delta[0]*delta[0] + delta[1]*delta[1] + delta[2]*delta[2]
			edges = append(edges, edge{lh: i, rh: j, distance: distance})
		}
	}
	slices.SortFunc(edges, func(lh edge, rh edge) int {
		return cmp.Compare(lh.distance, rh.distance)
	})
	slog.Debug("input", "boxes", boxes, "arg", context.GeneralArg())
	slog.Debug("edges", "edges", edges)
	uf := makeSuf(len(boxes))
	n, err := strconv.Atoi(context.GeneralArg())
	if err != nil {
		return err
	}
	for _, edge := range edges[:n] {
		uf.union(edge.lh, edge.rh)
	}
	circuitSizes := make(map[int]int)
	for i := range boxes {
		circuitSizes[uf.find(i)] += 1
	}
	slog.Debug("sizes", "circuitSizes", circuitSizes)
	justSizes := make([]int, 0)
	for _, s := range circuitSizes {
		justSizes = append(justSizes, s)
	}
	slices.Sort(justSizes)
	justSizes = justSizes[len(justSizes)-3:]
	context.Solution("A", justSizes[0]*justSizes[1]*justSizes[2])
	return nil
}
