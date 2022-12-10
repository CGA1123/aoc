package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

func right(p aoc.Point, n int64) aoc.Point {
	return aoc.Point{X: p.X + n, Y: p.Y}
}

func left(p aoc.Point, n int64) aoc.Point {
	return aoc.Point{X: p.X - n, Y: p.Y}
}

func up(p aoc.Point, n int64) aoc.Point {
	return aoc.Point{X: p.X, Y: p.Y - n}
}

func down(p aoc.Point, n int64) aoc.Point {
	return aoc.Point{X: p.X, Y: p.Y + n}
}

var moves = map[string]func(aoc.Point, int64) aoc.Point{
	"R": right,
	"L": left,
	"U": up,
	"D": down,
}

func minus(a, b aoc.Point) aoc.Point {
	return aoc.Point{X: a.X - b.X, Y: a.Y - b.Y}
}

func chase(h, t aoc.Point) aoc.Point {
	diff := minus(h, t)

	// we're still neighbouring head
	if aoc.Abs(diff.X) <= 1 && aoc.Abs(diff.Y) <= 1 {
		return t
	}

	if diff.X == 0 {
		if diff.Y > 1 {
			return aoc.Point{X: t.X, Y: h.Y - 1}
		} else {
			return aoc.Point{X: t.X, Y: h.Y + 1}
		}
	}

	if diff.Y == 0 {
		if diff.X > 1 {
			return aoc.Point{X: h.X - 1, Y: t.Y}
		} else {
			return aoc.Point{X: h.X + 1, Y: t.Y}
		}
	}

	var x, y int64
	if diff.X > 0 {
		x = 1
	} else {
		x = -1
	}
	if diff.Y > 0 {
		y = 1
	} else {
		y = -1
	}

	return aoc.Point{X: t.X + x, Y: t.Y + y}
}

func main() {
	doIt(2, "input.txt")
	doIt(10, "input.txt")
}

func doIt(knotsCount int, input string) {
	knots := make([]aoc.Point, knotsCount)
	visited := map[aoc.Point]struct{}{}

	aoc.EachLine(input, func(s string) {
		line := strings.Split(s, " ")
		amount := aoc.MustParse(line[1])

		move := moves[line[0]]
		for i := int64(0); i < amount; i++ {
			knots[0] = move(knots[0], 1)

			for j := 1; j < len(knots); j++ {
				knots[j] = chase(knots[j-1], knots[j])
			}

			visited[knots[len(knots)-1]] = struct{}{}
		}
	})

	log.Printf("%v", len(visited))
}
