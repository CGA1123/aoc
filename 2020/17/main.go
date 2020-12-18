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

func (p PointND) Add(o PointND) PointND {
	n := make([]int64, len(p))

	for i, pi := range p {
		n[i] = pi + o[i]
	}

	return PointND(n)
}

type GridND struct {
	grid map[string]element
}

type element struct {
	el    interface{}
	point PointND
}

func NewGridND(n int64) *GridND {
	return &GridND{grid: map[string]element{}}
}

func (g *GridND) Read(p PointND) interface{} {
	if elem, ok := g.grid[p.Key()]; ok {
		return elem.el
	}

	return nil
}

func (g *GridND) EachSparse(fn func(PointND, interface{})) {
	for _, v := range g.grid {
		fn(v.point, v.el)
	}
}

func (g *GridND) Write(point PointND, el interface{}) {
	g.grid[point.Key()] = element{el: el, point: point}
}

func (g *GridND) Remove(point PointND) {
	delete(g.grid, point.Key())
}

func (g *GridND) Size() int {
	return len(g.grid)
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

func ActiveNeighbours(grid *GridND, point PointND, neighbours []PointND) int {
	var c int

	for _, p := range neighbours {
		n := point.Add(p)
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

func Tick(grid, tmp *GridND, point PointND, el interface{}, neighbours []PointND, ticked map[string]bool) {
	if _, ok := ticked[point.Key()]; ok {
		return
	}

	ticked[point.Key()] = true

	activeNeighbours := ActiveNeighbours(grid, point, neighbours)
	if el == nil {
		el = Inactive
	}

	if next := NewState(el.(byte), activeNeighbours); next == Active {
		tmp.Write(point, next)
	} else {
		tmp.Remove(point)
	}
}

func Do(dimensions int64) int {
	var tmp *GridND
	grid := NewGridND(dimensions)
	neighbours := NeighboursND(make([]int64, dimensions))

	var y int64

	aoc.EachLine("input.txt", func(line string) {
		for x, b := range []byte(line) {
			coords := make([]int64, dimensions)
			coords[0] = int64(x)
			coords[1] = y

			if b == Active {
				grid.Write(PointND(coords), b)
			}
		}

		y++
	})

	for t := 0; t < 6; t++ {
		tmp = NewGridND(dimensions)
		ticked := map[string]bool{}

		grid.EachSparse(func(point PointND, el interface{}) {
			Tick(grid, tmp, point, el, neighbours, ticked)

			for _, n := range neighbours {
				neighbour := point.Add(n)

				Tick(grid, tmp, neighbour, grid.Read(neighbour), neighbours, ticked)
			}
		})

		tmp, grid = nil, tmp
	}

	return grid.Size()
}

func main() {
	log.Printf("pt(1): %v", Do(3))
	log.Printf("pt(2): %v", Do(4))
}
