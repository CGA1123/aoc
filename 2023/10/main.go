package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

var (
	North = aoc.Point{X: 0, Y: -1}
	South = aoc.Point{X: 0, Y: 1}
	East  = aoc.Point{X: 1, Y: 0}
	West  = aoc.Point{X: -1, Y: 0}
)

var connections = map[string][]aoc.Point{
	".": {},
	"S": {},
	"|": {North, South},
	"-": {East, West},
	"L": {North, East},
	"J": {North, West},
	"7": {South, West},
	"F": {South, East},
}

var printers = map[string]string{
	"|": "║",
	"-": "═",
	"F": "╔",
	"7": "╗",
	"L": "╚",
	"J": "╝",
	"S": "*",
	".": " ",
}

type Pipe struct {
	At       aoc.Point
	Type     string
	Connects []aoc.Point
}

func (p *Pipe) Empty() bool {
	return p == nil
}

func (p *Pipe) IsConnectedTo(pt aoc.Point) bool {
	for _, conn := range p.Connects {
		if conn == pt {
			return true
		}
	}

	return false
}

func (p *Pipe) Connection(from aoc.Point) aoc.Point {
	if from == p.Connects[0] {
		return p.Connects[1]
	}
	if from == p.Connects[1] {
		return p.Connects[0]
	}

	panic("not connected")
}

func Parse(p aoc.Point, s string) *Pipe {
	connects, ok := connections[s]
	if !ok {
		return nil
	}

	c := []aoc.Point{}
	for _, cnx := range connects {
		c = append(c, add(p, cnx))
	}

	return &Pipe{At: p, Type: s, Connects: c}
}

func add(a, b aoc.Point) aoc.Point {
	return aoc.Point{X: a.X + b.X, Y: a.Y + b.Y}
}

func main() {
	var S aoc.Point
	g := aoc.NewGrid()
	y := int64(0)
	aoc.EachLine("input.txt", func(s string) {
		for x, c := range strings.Split(s, "") {
			p := aoc.Point{X: int64(x), Y: y}
			if c == "S" {
				log.Printf("S: %+v", p)
				S = p
			}

			g.Write(int64(x), y, Parse(p, c))
		}
		y++
	})

	log.Printf("Start: %+v", S)

	g.EachLine(func(i []interface{}) {
		for _, pt := range i {
			p, ok := pt.(*Pipe)
			if !ok {
				continue
			}

			fmt.Printf(printers[p.Type])
		}

		fmt.Print("\n")
	})

	connections := map[*Pipe]bool{}
	for _, direction := range []aoc.Point{North, South, East, West} {
		point := add(S, direction)
		pipe, ok := g.Read(point.X, point.Y).(*Pipe)
		if !ok {
			continue
		}

		if !pipe.IsConnectedTo(S) {
			continue
		}

		connections[pipe] = false
	}

	from := S
	var steps int64
	for pipe := range connections {
		if connections[pipe] {
			continue
		}

		i, p, connects := traverse(g, 0, from, S, pipe)
		if !connects {
			continue
		}

		steps = i
		connections[pipe] = true
		connections[p] = true
	}

	log.Printf("result: %v", (steps/2)+1)
}

func traverse(g *aoc.Grid, i int64, from, to aoc.Point, p *Pipe) (int64, *Pipe, bool) {
	if !p.IsConnectedTo(from) {
		return i, nil, false
	}

	next := p.Connection(from)
	if next == to {
		return i, p, true
	}

	pp := g.Read(next.X, next.Y).(*Pipe)

	return traverse(g, i+1, p.At, to, pp)
}
