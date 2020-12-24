package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type Tile int
type Direction int

const (
	White Tile = iota
	Black
)

const (
	East Direction = iota
	West
	NorthWest
	SouthEast
	NorthEast
	SouthWest
)

func (d Direction) String() string {
	return map[Direction]string{
		East:      "e",
		West:      "w",
		NorthWest: "nw",
		SouthEast: "se",
		NorthEast: "ne",
		SouthWest: "sw"}[d]
}

var Tokens = map[string]Direction{
	"e":  East,
	"w":  West,
	"nw": NorthWest,
	"se": SouthEast,
	"ne": NorthEast,
	"sw": SouthWest}

var Origin = aoc.PointND([]int64{0, 0, 0})
var Vectors = map[Direction]aoc.PointND{
	East:      aoc.PointND([]int64{1, 1, 0}),
	West:      aoc.PointND([]int64{-1, -1, 0}),
	NorthWest: aoc.PointND([]int64{-1, 0, 1}),
	SouthEast: aoc.PointND([]int64{1, 0, -1}),
	NorthEast: aoc.PointND([]int64{0, 1, 1}),
	SouthWest: aoc.PointND([]int64{0, -1, -1})}

func Step(point aoc.PointND, directions ...Direction) aoc.PointND {
	p := point

	for _, direction := range directions {
		p = p.Add(Vectors[direction])
	}

	return p
}

func Lex(line string) []Direction {
	var directions []Direction

	str := line
	for str != "" {
		for k, v := range Tokens {
			if strings.HasPrefix(str, k) {
				directions = append(directions, v)
				str = strings.TrimPrefix(str, k)
			}
		}
	}

	return directions
}

var neighbours = map[interface{}][]aoc.PointND{}

func Neighbours(point aoc.PointND) []aoc.PointND {
	key := point.Key()

	if n, ok := neighbours[key]; ok {
		return n
	}

	points := make([]aoc.PointND, 0, 6)

	for _, v := range Vectors {
		points = append(points, point.Add(v))
	}

	neighbours[key] = points

	return points
}

func BlackNeighbours(grid *aoc.GridND, point aoc.PointND) int {
	var total int

	for _, v := range Neighbours(point) {
		if grid.Read(v) != nil {
			total++
		}
	}

	return total
}

func NewState(t Tile, n int) Tile {
	if t == Black && (n == 0 || n > 2) {
		return White
	}

	if t == White && n == 2 {
		return Black
	}

	return t
}

func Tick(grid, tmp *aoc.GridND, p aoc.PointND, ticked map[interface{}]bool) {
	if _, ok := ticked[p.Key()]; ok {
		return
	}

	ticked[p.Key()] = true

	activeNeighbours := BlackNeighbours(grid, p)

	el := grid.Read(p)
	if el == nil {
		el = White
	}

	if next := NewState(el.(Tile), activeNeighbours); next == Black {
		tmp.Write(p, next)
	} else {
		tmp.Remove(p)
	}
}

func main() {
	cleanup, _ := aoc.Profile()
	defer cleanup()

	grid := aoc.NewGridND(3)

	aoc.EachLine("input.txt", func(line string) {
		path := Lex(line)
		point := Step(Origin, path...)

		if grid.Read(point) == nil {
			grid.Write(point, Black)
		} else {
			grid.Remove(point)
		}
	})

	log.Printf("pt(1): %v", grid.Size())

	var tmp *aoc.GridND

	for i := 0; i < 100; i++ {
		tmp = aoc.NewGridND(3)
		ticked := map[interface{}]bool{}

		grid.EachSparse(func(point aoc.PointND, el interface{}) {
			Tick(grid, tmp, point, ticked)

			for _, n := range Neighbours(point) {
				Tick(grid, tmp, n, ticked)
			}
		})

		tmp, grid = nil, tmp
	}

	log.Printf("pt(2): %v", grid.Size())
}
