package main

import (
	"log"

	"github.com/CGA1123/aoc/2019/intcode"
)

const (
	Black = 0
	White = 1
)

type Point struct {
	X int64
	Y int64
}

type Hull struct {
	minx int64
	maxx int64
	miny int64
	maxy int64
	grid map[Point]int64
}

func NewHull() *Hull {
	return &Hull{grid: map[Point]int64{}}
}

func (h *Hull) Read(x, y int64) int64 {
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

	return h.grid[Point{X: x, Y: y}]
}

func (h *Hull) Write(x, y, i int64) {
	h.grid[Point{X: x, Y: y}] = i
}

func (h *Hull) Painted() int64 {
	return int64(len(h.grid))
}

func (h *Hull) Grid() string {
	str := ""

	for y := h.maxy; y >= h.miny; y-- {
		for x := h.minx; x <= h.maxx; x++ {
			if val := h.grid[Point{X: x, Y: y}]; val == White {
				str += "0"
			} else {
				str += " "
			}
		}

		str += "\n"
	}

	return str
}

type Robot struct {
	hull      *Hull
	direction int64
	x         int64
	y         int64
	ins       int64
}

func NewRobot(h *Hull) *Robot {
	return &Robot{hull: h}
}

func (r *Robot) directions() []func() {
	return []func(){
		r.north,
		r.east,
		r.south,
		r.west}
}

func (r *Robot) rotations() []func() {
	return []func(){r.left, r.right}
}

func (r *Robot) left() {
	r.direction = (r.direction + 3) % 4
}

func (r *Robot) right() {
	r.direction = (r.direction + 1) % 4
}

func (r *Robot) forward() {
	r.directions()[r.direction]()
}

func (r *Robot) north() {
	r.y += 1
}

func (r *Robot) south() {
	r.y -= 1
}

func (r *Robot) east() {
	r.x += 1
}

func (r *Robot) west() {
	r.x -= 1
}

func (r *Robot) Read() int64 {
	colour := r.hull.Read(r.x, r.y)

	return colour
}

func (r *Robot) Write(i int64) {
	if r.ins == 0 {
		r.hull.Write(r.x, r.y, i)
	} else {
		r.rotations()[i]()
		r.forward()
	}

	r.ins = (r.ins + 1) % 2
}

func main() {
	mem, err := intcode.FileMem("input.txt")
	if err != nil {
		log.Printf("err loading mem: %v", err)
		return
	}

	h := NewHull()
	r := NewRobot(h)
	ic := intcode.New(mem, r, r, 0)
	for ic.Do() {
	}

	log.Printf("pt(1): %v", h.Painted())

	h2 := NewHull()
	h2.Write(0, 0, White)
	r2 := NewRobot(h2)
	ic2 := intcode.New(mem, r2, r2, 0)
	for ic2.Do() {
	}

	log.Printf("pt(2):\n%v", h2.Grid())
}
