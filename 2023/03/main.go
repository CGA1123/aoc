package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

var numbers = map[string]int64{
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"0": 0,
}

func main() {
	g := aoc.NewGrid()
	row := int64(0)

	aoc.EachLine("input.txt", func(s string) {
		for column, chr := range strings.Split(s, "") {
			col := int64(column)

			g.Write(col, row, chr)
		}

		row++
	})

	currentNumber := &Number{}
	numbers := []*Number{}

	for row := int64(0); row < g.Height(); row++ {
		for col := int64(0); col < g.Width(); col++ {
			v := g.Read(col, row).(string)
			sym := &Symbol{Value: v, Point: aoc.Point{X: col, Y: row}}

			if sym.IsNumber() {
				currentNumber.Symbols = append(currentNumber.Symbols, sym)
			} else {
				if currentNumber.Present() {
					numbers = append(numbers, currentNumber)
					currentNumber = &Number{}
				}
			}
		}
		if currentNumber.Present() {
			numbers = append(numbers, currentNumber)
			currentNumber = &Number{}
		}
	}

	var total int64

	gears := map[aoc.Point][]*Number{}

	for _, number := range numbers {
		number.Calculate()

		neighbours := number.Neighbours(g)
		if len(neighbours) == 0 {
			continue
		}

		total += number.Value

		for _, neighbour := range neighbours {
			if !neighbour.IsGear() {
				continue
			}

			if _, ok := gears[neighbour.Point]; !ok {
				gears[neighbour.Point] = []*Number{}
			}

			gears[neighbour.Point] = append(gears[neighbour.Point], number)
		}
	}

	log.Printf("part 1: %v", total)

	var gearTotal int64

	for _, nums := range gears {
		if len(nums) != 2 {
			continue
		}

		gearTotal += nums[0].Value * nums[1].Value
	}

	log.Printf("part 2: %v", gearTotal)
}

type Symbol struct {
	Point aoc.Point
	Value string
}

func (s *Symbol) IsGear() bool {
	return s.Value == "*"
}

func (s *Symbol) IsBlank() bool {
	return s.Value == "."
}

func (s *Symbol) IsNumber() bool {
	_, ok := numbers[s.Value]
	return ok
}

type Number struct {
	Symbols []*Symbol
	Value   int64
}

func (n *Number) Calculate() int64 {
	v, c := int64(0), int64(1)
	for i := len(n.Symbols) - 1; i >= 0; i-- {
		v += aoc.MustParse(n.Symbols[i].Value) * c
		c *= 10
	}
	n.Value = v

	return n.Value
}

func (n *Number) Present() bool {
	return len(n.Symbols) > 0
}

func (n *Number) Points() map[aoc.Point]struct{} {
	m := make(map[aoc.Point]struct{}, len(n.Symbols))
	for _, s := range n.Symbols {
		m[s.Point] = struct{}{}
	}

	return m
}

func (n *Number) Neighbours(g *aoc.Grid) []*Symbol {
	p := n.Points()

	result := map[aoc.Point]*Symbol{}

	for _, sym := range n.Symbols {
		for _, neighbour := range g.Neighbours(sym.Point) {
			if _, ok := p[neighbour]; ok {
				continue
			}

			nsym := &Symbol{Value: g.Read(neighbour.X, neighbour.Y).(string), Point: neighbour}
			if nsym.IsNumber() || nsym.IsBlank() {
				continue
			}

			result[neighbour] = nsym
		}
	}

	syms := make([]*Symbol, 0, len(result))
	for _, s := range result {
		syms = append(syms, s)
	}

	return syms
}
