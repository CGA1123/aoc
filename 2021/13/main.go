package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

type instruction struct {
	axis   string
	target int64
}

func main() {
	g := aoc.NewGrid()

	atInstructions := false
	instructions := []instruction{}
	aoc.EachLine("input.txt", func(s string) {
		if s == "" {
			atInstructions = true
			return
		}

		if atInstructions {
			ins := strings.Split(strings.Replace(s, "fold along ", "", 1), "=")

			instructions = append(instructions, instruction{axis: ins[0], target: aoc.MustParse(ins[1])})

			return
		}

		coords := strings.Split(s, ",")
		x, y := aoc.MustParse(coords[0]), aoc.MustParse(coords[1])
		g.Write(x, y, true)
	})

	for _, ins := range instructions {
		if ins.axis == "x" {
			g = foldX(g, ins.target)
		} else {
			g = foldY(g, ins.target)
		}
	}

	display(g)
}

func foldX(g *aoc.Grid, at int64) *aoc.Grid {
	x := int64(0)
	folded := false
	ng := aoc.NewGrid()

	g.EachCol(func(i []interface{}) {
		// skip the fold line
		if at == x {
			folded = true
			x = x - 1
			return
		}

		for y, v := range i {
			if ng.Read(int64(x), int64(y)) == nil {
				ng.Write(int64(x), int64(y), v)
			}
		}

		if folded {
			x = x - 1
		} else {
			x = x + 1
		}
	})

	return ng
}

func foldY(g *aoc.Grid, at int64) *aoc.Grid {
	y := int64(0)
	folded := false
	ng := aoc.NewGrid()

	g.EachLine(func(i []interface{}) {
		// skip the fold line
		if at == y {
			folded = true
			y = y - 1
			return
		}

		for x, v := range i {
			if ng.Read(int64(x), y) == nil {
				ng.Write(int64(x), y, v)
			}
		}

		if folded {
			y = y - 1
		} else {
			y = y + 1
		}
	})

	return ng
}

func display(g *aoc.Grid) {
	g.EachLine(func(i []interface{}) {
		for _, v := range i {
			if v == nil {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	})
}
