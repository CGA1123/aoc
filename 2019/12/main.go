package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

var moon = regexp.MustCompile(`<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`)

func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int64, integers ...int64) int64 {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func MustParse(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("bad input")
	}

	return i
}

func Combinations(n, r int) [][]int {
	var combinations func(int, int, int, []int) [][]int
	combinations = func(n, r, min int, prefix []int) [][]int {
		toBeSelected := r - len(prefix)

		// the buffer is full!
		if toBeSelected == 0 {
			return [][]int{prefix}
		}

		var combs [][]int

		for i := min; i <= (n - toBeSelected); i++ {
			combs = append(
				combs,
				combinations(n, r, i, append(prefix, i))...,
			)
		}

		return combs
	}

	return combinations(n, r, 0, []int{})
}

type Coord struct {
	X int64
	Y int64
	Z int64
}

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}

	return i
}

func (c *Coord) Sum() int64 {
	return abs(c.X) + abs(c.Y) + abs(c.Z)
}

type Object struct {
	Pos *Coord
	Vec *Coord
}

func (o *Object) EPot() int64 {
	return o.Pos.Sum()
}

func (o *Object) EKin() int64 {
	return o.Vec.Sum()
}

func (o *Object) ETotal() int64 {
	return o.EKin() * o.EPot()
}

func (o *Object) EqDimension(other *Object, f func(*Coord) int64) bool {
	return (f(o.Pos) == f(other.Pos)) && (f(o.Vec) == f(other.Vec))
}

type Space struct {
	time    int64
	objects []*Object
	pairs   [][2]*Object
}

func NewSpace(obj []*Object) *Space {
	var pairs [][2]*Object
	for _, pair := range Combinations(len(obj), 2) {
		pairs = append(pairs, [2]*Object{
			obj[pair[0]],
			obj[pair[1]],
		})
	}

	return &Space{objects: obj, pairs: pairs}
}

func increment(a, b int64) (int64, int64) {
	if a == b {
		return 0, 0
	}

	if a > b {
		return -1, 1
	} else {
		return 1, -1
	}
}

func (s *Space) updateVelocity(a, b *Object) {
	ax, bx := increment(a.Pos.X, b.Pos.X)
	a.Vec.X += ax
	b.Vec.X += bx

	ay, by := increment(a.Pos.Y, b.Pos.Y)
	a.Vec.Y += ay
	b.Vec.Y += by

	az, bz := increment(a.Pos.Z, b.Pos.Z)
	a.Vec.Z += az
	b.Vec.Z += bz
}

func (s *Space) updatePosition(b *Object) {
	b.Pos.X += b.Vec.X
	b.Pos.Y += b.Vec.Y
	b.Pos.Z += b.Vec.Z
}

func (s *Space) Step() int64 {
	for _, pair := range s.pairs {
		s.updateVelocity(pair[0], pair[1])
	}

	for _, object := range s.objects {
		s.updatePosition(object)
	}

	s.time += 1

	return s.time
}

func (s *Space) ETotal() int64 {
	var i int64
	for _, o := range s.objects {
		i += o.ETotal()
	}

	return i
}

func (s *Space) Objects() []*Object {
	return s.objects
}

func AllEq(a, b []*Object, f func(*Coord) int64) bool {
	for i, ai := range a {
		if !ai.EqDimension(b[i], f) {
			return false
		}
	}

	return true
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	var obj []*Object
	var original []*Object
	s := bufio.NewScanner(bufio.NewReader(f))
	for s.Scan() {
		m := moon.FindStringSubmatch(s.Text())

		x, y, z := MustParse(m[1]), MustParse(m[2]), MustParse(m[3])

		obj = append(obj, &Object{Pos: &Coord{X: x, Y: y, Z: z}, Vec: &Coord{}})
		original = append(original, &Object{Pos: &Coord{X: x, Y: y, Z: z}, Vec: &Coord{}})
	}
	if s.Err() != nil {
		log.Printf("error scanning: %v", s.Err())
		return
	}

	space := NewSpace(obj)

	var x int64
	var y int64
	var z int64

	for {
		t := space.Step()
		if t == 1000 {
			log.Printf("pt(1): %v", space.ETotal())
		}

		if x != 0 && y != 0 && z != 0 {
			break
		}

		o := space.Objects()
		if x == 0 && AllEq(o, original, func(c *Coord) int64 { return c.X }) {
			x = t
		}

		if y == 0 && AllEq(o, original, func(c *Coord) int64 { return c.Y }) {
			y = t
		}

		if z == 0 && AllEq(o, original, func(c *Coord) int64 { return c.Z }) {
			z = t
		}
	}

	log.Printf("pt(2) %v", lcm(x, y, z))
}
