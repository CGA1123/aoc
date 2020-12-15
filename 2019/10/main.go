package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
)

const (
	Empty    = '.'
	Asteroid = '#'
)

func angle(x, y, xi, yi int) float64 {
	dx, dy := xi-x, yi-y

	deg := (math.Atan2(float64(dy), float64(dx)) * 180) / math.Pi
	if deg < 0 {
		deg += 360
	}

	// munge to match coordinate system of grid
	deg += 90

	if deg >= 360 {
		deg -= 360
	}

	return deg
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

type Point struct {
	X int
	Y int
}

func distance(x, y, xi, yi int) float64 {
	a, b := abs(x-xi), abs(y-yi)

	return math.Sqrt(float64((a * a) + (b * b)))
}

func detect(grid [][]byte, x, y int) map[float64]Point {
	detected := map[float64]Point{}

	for yi := 0; yi < len(grid); yi++ {
		for xi := 0; xi < len(grid[yi]); xi++ {
			if (xi == x && yi == y) || grid[yi][xi] == Empty {
				continue
			}

			a := angle(x, y, xi, yi)
			d := distance(x, y, xi, yi)

			if p, ok := detected[a]; !ok || d < distance(x, y, p.X, p.Y) {
				detected[a] = Point{X: xi, Y: yi}
			}
		}
	}

	return detected
}

func PartOne(grid [][]byte) (int, Point) {
	var most int
	var point Point

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == Empty {
				continue
			}

			if c := len(detect(grid, x, y)); c > most {
				point = Point{X: x, Y: y}
				most = c
			}
		}
	}

	return most, point
}

func PartTwo(grid [][]byte, point Point) Point {
	var i int
	for {
		points := detect(grid, point.X, point.Y)
		keys := make([]float64, 0, len(points))
		for k := range points {
			keys = append(keys, k)
		}

		sort.Float64s(keys)

		for _, k := range keys {
			p := points[k]
			grid[p.Y][p.X] = Empty

			i++
			if i == 200 {
				return p
			}
		}
	}
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	var grid [][]byte
	s := bufio.NewScanner(bufio.NewReader(f))
	for s.Scan() {
		line := []byte(s.Text())
		grid = append(grid, line)
	}
	if s.Err() != nil {
		log.Printf("error scanning: %v", s.Err())
		return
	}

	most, point := PartOne(grid)
	log.Printf("pt(1) (%v, %v): %v", point.X, point.Y, most)

	point2 := PartTwo(grid, point)
	log.Printf("pt(2) (%v, %v): %v", point2.X, point2.Y, (point2.X*100)+point2.Y)
}
