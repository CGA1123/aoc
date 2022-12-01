package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	p := map[aoc.Point]int{}

	aoc.EachLine("input.txt", func(input string) {
		line := strings.Split(input, " -> ")
		from, to := strings.Split(line[0], ","), strings.Split(line[1], ",")

		x1, y1, x2, y2 := aoc.MustParse(from[0]), aoc.MustParse(from[1]), aoc.MustParse(to[0]), aoc.MustParse(to[1])
		for _, pp := range points(x1, y1, x2, y2) {
			p[pp] = p[pp] + 1
		}
	})

	score := 0
	for _, i := range p {
		if i >= 2 {
			score = score + 1
		}
	}

	log.Printf("%v", score)
}

func points(x1, y1, x2, y2 int64) []aoc.Point {
	if x1 != x2 && y1 != y2 {
		return diagonal(x1, y1, x2, y2)
	}

	p := []aoc.Point{}
	if x1 == x2 {
		for _, y := range between(y1, y2) {
			p = append(p, aoc.Point{X: x1, Y: y})
		}
	} else {
		for _, x := range between(x1, x2) {
			p = append(p, aoc.Point{X: x, Y: y1})
		}
	}

	return p
}

func between(a, b int64) []int64 {
	if a > b {
		a, b = b, a
	}

	n := []int64{}
	for i := a; i <= b; i++ {
		n = append(n, i)
	}

	return n
}

func diagonal(x1, y1, x2, y2 int64) []aoc.Point {
	xd := (x2 - x1) / aoc.Abs(x2-x1)
	yd := (y2 - y1) / aoc.Abs(y2-y1)

	p := []aoc.Point{}
	for i := int64(0); i <= aoc.Abs(x2-x1); i++ {
		p = append(p, aoc.Point{
			X: x1 + (i * xd),
			Y: y1 + (i * yd),
		})
	}

	return p
}
