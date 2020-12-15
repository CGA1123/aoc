package main

import "log"

type Game struct {
	turn     int
	last     int
	lastSaid map[int]int
}

func NewGame(init []int) *Game {
	g := &Game{lastSaid: map[int]int{}}

	for _, n := range init {
		g.last = n
		g.turn += 1
		g.lastSaid[g.last] = g.turn
	}

	return g
}

func (g *Game) Do() {
	saidOn, ok := g.lastSaid[g.last]

	g.lastSaid[g.last] = g.turn

	if !ok {
		g.last = 0
	} else {
		g.last = g.turn - saidOn
	}

	g.turn += 1
}

func (g *Game) LastSaid() int {
	return g.last
}

func (g *Game) Turn() int {
	return g.turn
}

func main() {
	input := []int{2, 1, 10, 11, 0, 6}

	g := NewGame(input)
	for g.Turn() < 30000000 {
		g.Do()
	}

	log.Printf("%v", g.LastSaid())
}
