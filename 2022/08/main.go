package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	trees := aoc.NewGrid()

	y := int64(0)
	aoc.EachLine("input.txt", func(s string) {
		for x, t := range strings.Split(s, "") {
			trees.Write(int64(x), y, aoc.MustParse(t))
		}

		y++
	})

	c := 0
	s := int64(0)
	trees.EachSparse(func(p aoc.Point, i interface{}) {
		if visible(trees, p.X, p.Y) {
			c++
		}

		scenic := scenicScore(trees, p.X, p.Y)
		if scenic > s {
			s = scenic
		}

	})

	fmt.Printf("%v\n", c)
	fmt.Printf("%v\n", s)
}

func scenicScore(g *aoc.Grid, x, y int64) int64 {
	th := g.Read(x, y).(int64)
	row, col := g.Row(y), g.Col(x)
	leftr, rightr := row[:x], row[(x+1):]
	topc, bottomc := col[:y], col[(y+1):]

	sleftr := score(reverse(leftr), th)
	srightr := score(rightr, th)
	stopc := score(reverse(topc), th)
	sbotc := score(bottomc, th)

	return sleftr * srightr * stopc * sbotc
}

func reverse(x []interface{}) []interface{} {
	y := make([]interface{}, 0, len(x))

	for i := len(x) - 1; i >= 0; i-- {
		y = append(y, x[i])
	}

	return y
}

func score(heights []interface{}, height int64) int64 {
	s := int64(0)

	for _, h := range heights {
		s++

		if h.(int64) >= height {
			break
		}
	}

	return s
}

func visible(g *aoc.Grid, x, y int64) bool {
	th := g.Read(x, y).(int64)
	row, col := g.Row(y), g.Col(x)
	leftr, rightr := row[:x], row[(x+1):]
	topc, bottomc := col[:y], col[(y+1):]

	w, h := g.Width(), g.Height()
	if x == 0 || x == (w-1) || y == 0 || y == (h-1) {
		return true
	}

	return allLess(leftr, th) || allLess(rightr, th) || allLess(topc, th) || allLess(bottomc, th)
}

func allLess(heights []interface{}, height int64) bool {
	for _, h := range heights {
		if h.(int64) >= height {
			return false
		}
	}

	return true
}
