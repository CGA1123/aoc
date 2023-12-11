package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	universe := parseUniverse("input.txt")
	galaxies := galaxies(universe)
	rows, cols := empties(universe.EachLine), empties(universe.EachCol)

	var sumOfDistances, sumOfDistances2 int64
	for _, pair := range aoc.Combinations(len(galaxies), 2) {
		i, j := pair[0], pair[1]
		if i == j {
			continue
		}

		gi, gj := galaxies[i], galaxies[j]
		distance := manhattan(gi, gj)
		adjustment := adjustDistance(gi, gj, rows, cols)

		sumOfDistances += distance + adjustment
		sumOfDistances2 += distance + (adjustment * (1000000 - 1))
	}

	log.Printf("%v", sumOfDistances)
	log.Printf("%v", sumOfDistances2)
}

func emptiesBetween(a, b int64, empties []int64) int64 {
	var i int64
	if a > b {
		a, b = b, a
	}

	for _, e := range empties {
		if e < a {
			continue
		}
		if e > b {
			break
		}

		i++
	}

	return i
}

func showGrid(g *aoc.Grid) {
	g.EachLine(func(els []interface{}) {
		for _, x := range els {
			fmt.Printf(x.(string))
		}
		fmt.Printf("\n")
	})
}

func adjustDistance(a, b aoc.Point, rows, cols []int64) int64 {
	return emptiesBetween(a.X, b.X, cols) + emptiesBetween(a.Y, b.Y, rows)
}

func manhattan(a, b aoc.Point) int64 {
	return aoc.Abs(a.X-b.X) + aoc.Abs(a.Y-b.Y)
}

func empties(gen func(func([]interface{}))) []int64 {
	idxs := []int64{}
	var i int64

	gen(func(els []interface{}) {
		if All(els, func(el interface{}) bool {
			return el.(string) == "."
		}) {
			idxs = append(idxs, i)
		}

		i++
	})

	return idxs
}

func Any[A any](els []A, f func(A) bool) bool {
	for _, e := range els {
		if f(e) {
			return true
		}
	}

	return false
}

func All[A any](els []A, f func(A) bool) bool {
	for _, e := range els {
		if !f(e) {
			return false
		}
	}

	return true
}

func parseUniverse(input string) *aoc.Grid {
	g := aoc.NewGrid()

	var y int64
	aoc.EachLine(input, func(s string) {
		for x, e := range strings.Split(s, "") {
			g.Write(int64(x), y, e)
		}

		y++
	})

	return g
}

func galaxies(og *aoc.Grid) []aoc.Point {
	galaxies := []aoc.Point{}
	var y int64
	og.EachLine(func(els []interface{}) {
		for i, el := range els {
			x := int64(i)
			e := el.(string)

			if e == "#" {
				galaxies = append(galaxies, aoc.Point{X: x, Y: y})
			}
		}

		y++
	})

	return galaxies
}
