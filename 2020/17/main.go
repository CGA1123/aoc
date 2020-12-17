package main

import (
	"fmt"
	"log"

	"github.com/CGA1123/aoc"
)

const (
	Active   = byte('#')
	Inactive = byte('.')
)

type PointND []int64

func (p PointND) Key() string {
	var str string
	for _, x := range p {
		str = fmt.Sprintf("%v,%v", str, x)
	}

	return str
}

func (p PointND) Equal(o PointND) bool {
	if len(p) != len(o) {
		return false
	}

	for i, n := range p {
		if o[i] != n {
			return false
		}
	}

	return true
}

type GridND struct {
	grid map[string]interface{}
	min  []int64
	max  []int64
	n    int64
}

func NewGridND(n int64) *GridND {
	return &GridND{grid: map[string]interface{}{},
		n:   n,
		min: make([]int64, n),
		max: make([]int64, n)}
}

func (g *GridND) Read(p PointND) interface{} {
	return g.grid[p.Key()]
}

func (g *GridND) genEach(current []int64, buf int64, fn func(PointND, interface{})) {
	idx := len(current)

	if int64(idx) == g.n {
		point := PointND(current)

		fn(point, g.grid[point.Key()])
		return
	}

	for i := g.min[idx] - buf; i <= g.max[idx]+buf; i++ {
		candidate := make([]int64, len(current))
		copy(candidate, current)

		g.genEach(append(candidate, i), buf, fn)
	}
}

func (g *GridND) Each(buf int64, fn func(PointND, interface{})) {
	g.genEach([]int64{}, buf, fn)
}

func (g *GridND) Write(point PointND, el interface{}) {
	for i, c := range point {
		if c < g.min[i] {
			g.min[i] = c
		}

		if c > g.max[i] {
			g.max[i] = c
		}

	}

	g.grid[point.Key()] = el
}

func genNeighbour(original, current PointND, idx int) []PointND {
	var neighbours []PointND

	if idx == len(current) {
		if original.Equal(current) {
			return []PointND{}
		} else {
			return []PointND{current}
		}
	}

	for i := int64(-1); i <= 1; i++ {
		candidate := make([]int64, len(current))
		copy(candidate, current)

		candidate[idx] += i

		neighbours = append(
			neighbours,
			genNeighbour(original, PointND(candidate), idx+1)...,
		)
	}

	return neighbours
}

func NeighboursND(point PointND) []PointND {
	candidate := make([]int64, len(point))
	copy(candidate, point)

	return genNeighbour(point, candidate, 0)
}

func ActiveNeighbours(grid *GridND, point PointND) int {
	var c int

	for _, n := range NeighboursND(point) {
		el := grid.Read(n)
		if el == nil || el.(byte) == Inactive {
			continue
		}

		c++
	}

	return c
}

func NewState(el byte, count int) byte {
	if el == Active && 2 <= count && count <= 3 {
		return Active
	}

	if el == Inactive && count == 3 {
		return Active
	}

	return Inactive
}

func CountActive(g *GridND) int {
	var c int

	g.Each(0, func(point PointND, el interface{}) {
		if el != nil && el.(byte) == Active {
			c++
		}
	})

	return c
}

func Do(dimensions int64) int {
	var tmp *GridND
	grid := NewGridND(dimensions)

	var y int64

	aoc.EachLine("input.txt", func(line string) {
		for x, b := range []byte(line) {
			coords := make([]int64, dimensions)
			coords[0] = int64(x)
			coords[1] = y

			grid.Write(PointND(coords), byte(b))
		}

		y++
	})

	for t := 0; t < 6; t++ {
		tmp = NewGridND(dimensions)

		grid.Each(1, func(point PointND, el interface{}) {
			activeNeighbours := ActiveNeighbours(grid, point)
			if el == nil {
				el = Inactive
			}

			next := NewState(el.(byte), activeNeighbours)

			tmp.Write(point, next)
		})

		tmp, grid = nil, tmp
	}

	return CountActive(grid)
}

func main() {
	log.Printf("pt(1): %v", Do(3))
	log.Printf("pt(2): %v", Do(4))
}
