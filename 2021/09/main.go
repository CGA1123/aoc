package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	g := aoc.NewGrid()

	y := int64(0)
	aoc.EachLine("input.txt", func(s string) {
		for x, v := range strings.Split(s, "") {
			n := aoc.MustParse(v)

			g.Write(int64(x), y, n)
		}
		y = y + 1
	})

	risk := int64(0)
	g.EachSparse(func(p aoc.Point, i interface{}) {
		h := i.(int64)

		isLow := true
		for _, n := range g.Neighbours(p) {
			if diagonal(n, p) {
				continue
			}

			nh := g.Read(n.X, n.Y).(int64)
			if nh < h {
				isLow = false
				break
			}
		}

		if isLow {
			risk = risk + int64(1) + h
		}
	})

	log.Printf("%v", risk)

	explored := make(map[aoc.Point]int, g.Count())
	basin := 1
	g.EachSparse(func(p aoc.Point, i interface{}) {
		// check if we've already seen this point.
		if _, ok := explored[p]; ok {
			return
		}

		explore(g, p, basin, explored)

		// we've explored this basin, move on.
		basin = basin + 1
	})

	basinSizes := map[int]int{}
	for _, id := range explored {
		basinSizes[id] = basinSizes[id] + 1
	}

	sizes := []int{}
	for id, size := range basinSizes {
		if id == -1 {
			continue
		}
		sizes = append(sizes, size)
	}

	sort.Ints(sizes)

	total := 1
	for _, i := range sizes[len(sizes)-3:] {
		total = total * i
	}
	fmt.Printf("%v\n", total)
}

func explore(cave *aoc.Grid, p aoc.Point, basin int, explored map[aoc.Point]int) {
	if _, ok := explored[p]; ok {
		return
	}

	h := cave.Read(p.X, p.Y).(int64)

	// check if this point is a boundary
	if h == 9 {
		// this point is not in a basin.
		explored[p] = -1
		return
	}

	// this point is in the current basin.
	explored[p] = basin

	// explore each neighbour
	for _, n := range cave.Neighbours(p) {
		if diagonal(p, n) {
			continue
		}

		explore(cave, n, basin, explored)
	}
}

func diagonal(a, b aoc.Point) bool {
	return a.X != b.X && a.Y != b.Y
}
