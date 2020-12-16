package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/CGA1123/aoc"
	"github.com/CGA1123/aoc/2019/intcode"
)

var tiles = map[int64]string{
	0: " ",
	1: "*",
	2: "#",
	3: "_",
	4: "o"}

func PartOne(mem []int64) int64 {
	grid := aoc.NewGrid()
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

	var c int64
	grid.Each(func(el interface{}) {
		if el != nil && el.(int64) == 2 {
			c++
		}
	})

	return c
}

type Game struct {
	paddleX int64
	ballX   int64
	grid    *aoc.Grid
	score   int64
	i       int64
	buf     [3]int64
}

func (j *Game) Read() int64 {

	PrintGrid(j.grid)
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

func PrintGrid(g *aoc.Grid) {
	fmt.Printf("\033[4;0H")

	g.EachLine(func(line []interface{}) {
		for _, el := range line {
			if el == nil {
				el = 0
			}

			fmt.Printf("%v", tiles[el.(int64)])
		}

		fmt.Printf("\n")
	})
}

func PartTwo(mem []int64) int64 {
	mem[0] = 2
	grid := aoc.NewGrid()

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
