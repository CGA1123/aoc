package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type octopus struct {
	energy  int64
	flashed bool
}

func (o *octopus) Print() {
	if o.flashed {
		fmt.Print("\033[5m*\033[0m")
	} else {
		fmt.Printf("%v", o.energy)
	}
}

func main() {
	g := aoc.NewGrid()
	y := int64(0)
	aoc.EachLine("input.txt", func(s string) {
		for x, i := range strings.Split(s, "") {
			n := aoc.MustParse(i)
			g.Write(int64(x), y, &octopus{energy: n})
		}

		y = y + 1
	})

	flashes := 0
	for i := 0; i < 1000; i++ {
		// increase by 1 for each octopus
		eachOctopus(g, func(p aoc.Point, o *octopus) {
			o.energy = o.energy + 1
		})

		// start flashing
		eachOctopus(g, func(p aoc.Point, o *octopus) {
			newFlashes := flash(g, p, o)

			if i < 100 {
				flashes = flashes + newFlashes
			}
		})

		allFlashed := true
		// stop flashing
		eachOctopus(g, func(p aoc.Point, o *octopus) {
			allFlashed = allFlashed && o.flashed
			o.flashed = false
		})

		if allFlashed {
			log.Printf("part 2: %v", i+1)
			break
		}
	}

	log.Printf("part 1: %v\n", flashes)
}

func eachOctopus(g *aoc.Grid, f func(aoc.Point, *octopus)) {
	g.EachSparse(func(p aoc.Point, i interface{}) {
		o := i.(*octopus)
		f(p, o)
	})
}

func flash(g *aoc.Grid, p aoc.Point, o *octopus) int {
	if o.energy <= 9 || o.flashed {
		return 0
	}

	flashes := 0
	o.energy = 0
	o.flashed = true
	flashes = flashes + 1

	for _, n := range g.Neighbours(p) {
		no := g.Read(n.X, n.Y).(*octopus)
		if no.flashed {
			continue
		}

		no.energy = no.energy + 1
		flashes = flashes + flash(g, n, no)
	}

	return flashes
}

func display(g *aoc.Grid) {
	g.EachLine(func(i []interface{}) {
		for _, i := range i {
			i.(*octopus).Print()
		}
		fmt.Print("\n")
	})
}
