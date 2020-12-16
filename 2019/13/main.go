package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/CGA1123/aoc/2019/intcode"
)

var tiles = map[int64]string{
	0: " ",
	1: "*",
	2: "#",
	3: "_",
	4: "o"}

type Point struct {
	X int64
	Y int64
}

type Grid struct {
	minx int64
	maxx int64
	miny int64
	maxy int64
	grid map[Point]int64
}

func NewGrid() *Grid {
	return &Grid{grid: map[Point]int64{}}
}

func (h *Grid) Read(x, y int64) int64 {
	return h.grid[Point{X: x, Y: y}]
}

func (h *Grid) Write(x, y, i int64) {
	if x > h.maxx {
		h.maxx = x
	}

	if y > h.maxy {
		h.maxy = y
	}

	if x < h.minx {
		h.minx = x
	}

	if y < h.miny {
		h.miny = y
	}

	h.grid[Point{X: x, Y: y}] = i
}

func (h *Grid) Grid() [][]int64 {
	var grid [][]int64

	for y := h.maxy; y >= h.miny; y-- {
		var line []int64
		for x := h.minx; x <= h.maxx; x++ {
			line = append(line, h.grid[Point{X: x, Y: y}])
		}

		grid = append(grid, line)
	}

	return grid
}

func PartOne(mem []int64) int64 {
	grid := NewGrid()
	out := intcode.NewChanIO(3)
	ic := intcode.New(mem, intcode.NewNullIO(), out, 0)

	running := true
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for running {
			if len(out.Chan()) == 3 {
				grid.Write(
					out.Read(),
					out.Read(),
					out.Read())
			}
		}
		wg.Done()
	}()

	for ic.Do() {
	}
	running = false

	wg.Wait()
	close(out.Chan())

	g := grid.Grid()
	var c int64
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if g[y][x] == 2 {
				c++
			}
		}
	}

	return c
}

type Game struct {
	paddleX int64
	ballX   int64
	grid    *Grid
	score   int64
	i       int64
	buf     [3]int64
}

func (j *Game) Read() int64 {

	PrintGrid(j.grid.Grid())
	fmt.Printf("score: %v\n", j.score)

	if j.ballX == j.paddleX {
		return 0
	}

	if j.ballX > j.paddleX {
		return 1
	}

	return -1
}

func (j *Game) Write(i int64) {
	j.buf[j.i%3] = i
	j.i += 1

	if j.i%3 == 0 {
		x, y, t := j.buf[0], j.buf[1], j.buf[2]
		if x == -1 && y == 0 {
			j.score = t
		} else {
			j.grid.Write(x, y, t)

			if t == 3 {
				j.paddleX = x
			}

			if t == 4 {
				j.ballX = x
			}
		}
	}
}

func PrintGrid(g [][]int64) {
	fmt.Printf("\033[2;0H")

	for y := len(g) - 1; y >= 0; y-- {
		for x := 0; x < len(g[y]); x++ {
			fmt.Printf("%v", tiles[g[y][x]])
		}
		fmt.Printf("\n")
	}
}

func PartTwo(mem []int64) int64 {
	mem[0] = 2
	grid := NewGrid()

	in := &Game{grid: grid, buf: [3]int64{}}
	ic := intcode.New(mem, in, in, 0)

	for ic.Do() {
	}

	return in.score
}

func main() {
	mem, err := intcode.FileMem("input.txt")
	if err != nil {
		log.Printf("err loading mem: %v", err)
		return
	}

	log.Printf("pt(1): %v", PartOne(mem))
	log.Printf("pt(2): %v", PartTwo(mem))
}
