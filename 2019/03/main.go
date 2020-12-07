package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func Between(a, b, i int) bool {
	if a > b {
		a, b = b, a
	}

	return a <= i && i <= b
}

type Point struct {
	X int
	Y int
}

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

func (p *Point) Manhattan() int {
	return abs(p.X) + abs(p.Y)
}

func U(p Point, v int) []Point {
	var points []Point
	for i := 0; i < v; i++ {
		points = append(points, Point{X: p.X, Y: p.Y + 1 + i})
	}

	return points
}

func D(p Point, v int) []Point {
	var points []Point
	for i := 0; i < v; i++ {
		points = append(points, Point{X: p.X, Y: p.Y - 1 - i})
	}

	return points
}

func L(p Point, v int) []Point {
	var points []Point
	for i := 0; i < v; i++ {
		points = append(points, Point{X: p.X - 1 - i, Y: p.Y})
	}

	return points
}

func R(p Point, v int) []Point {
	var points []Point
	for i := 0; i < v; i++ {
		points = append(points, Point{X: p.X + 1 + i, Y: p.Y})
	}

	return points
}

func Parse(path []string) []Point {
	points := []Point{{0, 0}}
	current := points[0]

	for _, v := range path {
		direction, valueStr := v[0], v[1:]
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		var p []Point

		switch direction {
		case 'U':
			p = U(current, value)
		case 'D':
			p = D(current, value)
		case 'L':
			p = L(current, value)
		case 'R':
			p = R(current, value)
		}

		points = append(points, p...)
		current = p[len(p)-1]
	}

	return points
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	lines := make([]map[Point]int, 2)
	var i int

	s := bufio.NewScanner(bufio.NewReader(f))
	for s.Scan() {
		line := strings.Split(s.Text(), ",")

		lines[i] = map[Point]int{}

		for j, p := range Parse(line) {
			if _, ok := lines[i][p]; !ok {
				lines[i][p] = j
			}
		}

		i++
	}
	if s.Err() != nil {
		log.Printf("error: %v", s.Err())
	}

	a, b := lines[0], lines[1]
	var intersection int
	var combined int

	log.Printf("len(a): %v", len(a))
	log.Printf("len(b): %v", len(b))

	for ka, va := range a {
		if ka.X == 0 && ka.Y == 0 {
			continue
		}

		vb, ok := b[ka]
		if !ok {
			continue
		}

		if intersection == 0 || ka.Manhattan() < intersection {
			intersection = ka.Manhattan()
		}

		if combined == 0 || va+vb < combined {
			combined = va + vb
		}
	}

	log.Printf("pt(1): %v", intersection)
	log.Printf("pt(2): %v", combined)
}
