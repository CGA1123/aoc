package main

import (
	"bufio"
	"log"
	"os"
)

func copyGrid(grid [][]byte) [][]byte {
	var newGrid [][]byte
	h := len(grid)

	for i := 0; i < h; i++ {
		w := len(grid[i])
		newGrid = append(newGrid, make([]byte, w))

		for j := 0; j < w; j++ {
			newGrid[i][j] = grid[i][j]
		}
	}

	return newGrid
}

// stepping fns
func west(y, x int) (int, int) {
	return y, x - 1
}

func east(y, x int) (int, int) {
	return y, x + 1
}

func north(y, x int) (int, int) {
	return y - 1, x
}

func south(y, x int) (int, int) {
	return y + 1, x
}

func southwest(y, x int) (int, int) {
	return south(west(y, x))
}

func southeast(y, x int) (int, int) {
	return south(east(y, x))
}

func northwest(y, x int) (int, int) {
	return north(west(y, x))
}

func northeast(y, x int) (int, int) {
	return north(east(y, x))
}

var directions = []func(int, int) (int, int){west, east, north, south, southwest, southeast, northeast, northwest}

const (
	Occupied = '#'
	Empty    = 'L'
	Floor    = '.'
)

type Life struct {
	w    int
	h    int
	grid [][]byte
	next [][]byte
}

func NewLife(grid [][]byte) *Life {
	w := len(grid[0])
	h := len(grid)
	l := &Life{
		w:    w,
		h:    h,
		next: copyGrid(grid),
		grid: grid,
	}

	return l
}

func (l *Life) Occupied() int {
	var c int
	for i := 0; i < l.h; i++ {
		for j := 0; j < l.w; j++ {
			if l.grid[i][j] == Occupied {
				c++
			}
		}
	}

	return c
}

func (l *Life) inBounds(y, x int) bool {
	return y >= 0 && x >= 0 && y < l.h && x < l.w
}

func (l *Life) closestIsOccupied(takeAStep func(int, int) (int, int), y, x int) bool {
	yi, xi := takeAStep(y, x)

	for l.inBounds(yi, xi) {
		current := l.grid[yi][xi]
		if current == Floor {
			yi, xi = takeAStep(yi, xi)
			continue
		}

		return l.grid[yi][xi] == Occupied
	}

	return false
}

func (l *Life) PartTwo(y, x int) (byte, bool) {
	current := l.grid[y][x]
	var occupiedNeighbours int

	for _, direction := range directions {
		if l.closestIsOccupied(direction, y, x) {
			occupiedNeighbours++
		}
	}

	if current == Empty && occupiedNeighbours == 0 {
		return Occupied, true
	}

	if current == Occupied && occupiedNeighbours >= 5 {
		return Empty, true
	}

	return current, false
}

func (l *Life) PartOne(y, x int) (byte, bool) {
	current := l.grid[y][x]
	var occupiedNeighbours int

	for _, direction := range directions {
		yi, xi := direction(y, x)
		if l.inBounds(yi, xi) && l.grid[yi][xi] == Occupied {
			occupiedNeighbours++
		}
	}

	if current == Empty && occupiedNeighbours == 0 {
		return Occupied, true
	}

	if current == Occupied && occupiedNeighbours >= 4 {
		return Empty, true
	}

	return current, false
}

func (l *Life) Tick(cellTick func(int, int) (byte, bool)) bool {
	changed := false

	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			cell, change := cellTick(y, x)

			changed = changed || change

			l.next[y][x] = cell
		}
	}

	l.grid, l.next = l.next, l.grid

	return changed
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	var grid [][]byte
	for s.Scan() {
		grid = append(grid, []byte(s.Text()))
	}

	l := NewLife(copyGrid(grid))
	for l.Tick(l.PartOne) {
	}

	m := NewLife(copyGrid(grid))
	for m.Tick(m.PartTwo) {
	}

	log.Printf("pt(1): %v", l.Occupied())
	log.Printf("pt(2): %v", m.Occupied())
}
