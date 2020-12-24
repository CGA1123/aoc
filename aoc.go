package aoc

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
)

func Abs(i int64) int64 {
	if i < 0 {
		return -1
	}

	return i
}

type Set struct {
	set map[interface{}]struct{}
}

func Intersection(sets ...*Set) *Set {
	n := NewSet()

	allContain := func(el interface{}) bool {
		for _, set := range sets {
			if !set.Contains(el) {
				return false
			}
		}

		return true
	}

	for _, seta := range sets {
		for _, el := range seta.Elements() {
			if allContain(el) {
				n.Add(el)
			}
		}
	}

	return n
}

func Union(sets ...*Set) *Set {
	n := NewSet()

	for _, set := range sets {
		for _, el := range set.Elements() {
			n.Add(el)
		}
	}

	return n
}

var set struct{}

func NewSetWithSize(i int) *Set {
	m := make(map[interface{}]struct{}, i)

	return &Set{set: m}
}

func NewSet() *Set {
	return &Set{set: map[interface{}]struct{}{}}
}

func (s *Set) Add(e interface{}) {
	s.set[e] = set
}

func (s *Set) Remove(e interface{}) {
	delete(s.set, e)
}

func (s *Set) Contains(e interface{}) bool {
	_, ok := s.set[e]

	return ok
}

func (s *Set) Size() int {
	return len(s.set)
}

func (s *Set) Elements() []interface{} {
	var el []interface{}

	for k := range s.set {
		el = append(el, k)
	}

	return el
}

func MustParse(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("bad input")
	}

	return i
}

func Capture(r *regexp.Regexp, s string) map[string]string {
	m := map[string]string{}
	names := r.SubexpNames()
	for _, match := range r.FindAllStringSubmatch(s, -1) {
		for i, submatch := range match {
			name := names[i]
			if name == "" {
				continue
			}

			m[name] = submatch
		}
	}

	return m
}

func Between(min, max int64) func(int64) bool {
	return func(i int64) bool {
		return min <= i && i <= max
	}
}

func Or(f, g func(int64) bool) func(int64) bool {
	return func(i int64) bool {
		return f(i) || g(i)
	}
}

func EachLine(input string, fn func(string)) error {
	return Scanner(input, func(s *bufio.Scanner) {
		for s.Scan() {
			fn(s.Text())
		}
	})
}

func Scanner(input string, fn func(*bufio.Scanner)) error {
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("opening file (%v): %v", input, err)
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	fn(s)

	return s.Err()
}

type Point struct {
	X int64
	Y int64
}

type Grid struct {
	minx int64
	maxx int64
	miny int64
	maxy int64
	grid map[Point]interface{}
}

func NewGrid() *Grid {
	return &Grid{grid: map[Point]interface{}{}}
}

func (h *Grid) Read(x, y int64) interface{} {
	return h.grid[Point{X: x, Y: y}]
}

func (h *Grid) Write(x, y int64, i interface{}) {
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

func (h *Grid) Grid() [][]interface{} {
	var grid [][]interface{}

	h.EachLine(func(line []interface{}) {
		grid = append(grid, line)
	})

	return grid
}

func (h *Grid) EachLine(fn func([]interface{})) {
	for y := h.miny; y <= h.maxy; y++ {
		var line []interface{}

		for x := h.minx; x <= h.maxx; x++ {
			line = append(line, h.grid[Point{X: x, Y: y}])
		}

		fn(line)
	}
}

func (h *Grid) Each(fn func(interface{})) {
	for y := h.miny; y <= h.maxy; y++ {
		for x := h.minx; x <= h.maxx; x++ {
			fn(h.grid[Point{X: x, Y: y}])
		}
	}
}

func (h *Grid) EachSparse(fn func(Point, interface{})) {
	for point, element := range h.grid {
		fn(point, element)
	}
}

func (h *Grid) Height() int64 {
	return h.maxy - h.miny + 1
}

func (h *Grid) Width() int64 {
	return h.maxx - h.minx + 1
}

func (h *Grid) Count() int64 {
	return int64(len(h.grid))
}

func Profile() (func(), error) {
	f, err := os.Create("profile.cpu")
	if err != nil {
		return nil, fmt.Errorf("could not create CPU profile: %v", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("could not create CPU profile: %v", err)
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}

type PointND []int64

func CantorPair(a, b int64) int64 {
	return ((a + b) * (a + b + 1) / 2) + b
}

func CantorTuple(p []int64) int64 {
	t := CantorPair(ZN(p[0]), ZN(p[1]))

	for i := 2; i < len(p); i++ {
		t = CantorPair(t, ZN(p[i]))
	}

	return t
}

func ZN(i int64) int64 {
	if i < 0 {
		return i * -2
	}

	return (i * 2) - 1
}

func (p PointND) Key() interface{} {
	return CantorTuple(p)
}

func (p PointND) Equal(o PointND) bool {
	if len(p) != len(o) {
		return false
	}

	for i, n := range p {
		if o[i] != n {
			return false
		}
	}

	return true
}

func (p PointND) Add(o PointND) PointND {
	n := make([]int64, len(p))

	for i, pi := range p {
		n[i] = pi + o[i]
	}

	return PointND(n)
}

type GridND struct {
	grid map[interface{}]element
}

type element struct {
	el    interface{}
	point PointND
}

func NewGridND(n int64) *GridND {
	return &GridND{grid: map[interface{}]element{}}
}

func (g *GridND) Read(p PointND) interface{} {
	if elem, ok := g.grid[p.Key()]; ok {
		return elem.el
	}

	return nil
}

func (g *GridND) EachSparse(fn func(PointND, interface{})) {
	for _, v := range g.grid {
		fn(v.point, v.el)
	}
}

func (g *GridND) Write(point PointND, el interface{}) {
	g.grid[point.Key()] = element{el: el, point: point}
}

func (g *GridND) Remove(point PointND) {
	delete(g.grid, point.Key())
}

func (g *GridND) Size() int {
	return len(g.grid)
}

func genNeighbour(original, current PointND, idx int) []PointND {
	var neighbours []PointND

	if idx == len(current) {
		if original.Equal(current) {
			return []PointND{}
		} else {
			return []PointND{current}
		}
	}

	for i := int64(-1); i <= 1; i++ {
		candidate := make([]int64, len(current))
		copy(candidate, current)

		candidate[idx] += i

		neighbours = append(
			neighbours,
			genNeighbour(original, PointND(candidate), idx+1)...,
		)
	}

	return neighbours
}

func NeighboursND(point PointND) []PointND {
	candidate := make([]int64, len(point))
	copy(candidate, point)

	return genNeighbour(point, candidate, 0)
}
