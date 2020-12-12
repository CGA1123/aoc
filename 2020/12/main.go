package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func Abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

type WaypointShip struct {
	posEast  int
	posNorth int
	wayEast  int
	wayNorth int
}

func NewWaypointShip() *WaypointShip {
	return &WaypointShip{wayEast: 10, wayNorth: 1}
}

func (s *WaypointShip) Do(ins string, n int) {
	map[string]func(int){
		"F": s.Forward,
		"L": s.Left,
		"R": s.Right,
		"N": s.North,
		"E": s.East,
		"S": s.South,
		"W": s.West}[ins](n)
}

func (s *WaypointShip) Manhattan() int {
	return Abs(s.posEast) + Abs(s.posNorth)
}

func (s *WaypointShip) West(n int) {
	s.East(-n)
}

func (s *WaypointShip) East(n int) {
	s.wayEast += n
}

func (s *WaypointShip) North(n int) {
	s.wayNorth += n
}

func (s *WaypointShip) South(n int) {
	s.North(-n)
}

func (s *WaypointShip) Right(n int) {
	i := (n / 90) % 4

	for j := 0; j < i; j++ {
		s.wayEast, s.wayNorth = s.wayNorth, -s.wayEast
	}
}

func (s *WaypointShip) Left(n int) {
	s.Right(360 - n)
}

func (s *WaypointShip) Forward(n int) {
	s.posEast += (s.wayEast * n)
	s.posNorth += (s.wayNorth * n)
}

type Ship struct {
	east      int
	north     int
	direction int
}

func NewShip() *Ship {
	return &Ship{direction: 1}
}

func (s *Ship) Do(ins string, n int) {
	map[string]func(int){
		"F": s.Forward,
		"L": s.Left,
		"R": s.Right,
		"N": s.North,
		"E": s.East,
		"S": s.South,
		"W": s.West}[ins](n)
}

func (s *Ship) Manhattan() int {
	return Abs(s.east) + Abs(s.north)
}

func (s *Ship) directions() []func(int) {
	return []func(int){
		s.North,
		s.East,
		s.South,
		s.West}
}

func (s *Ship) West(n int) {
	s.East(-n)
}

func (s *Ship) East(n int) {
	s.east += n
}

func (s *Ship) North(n int) {
	s.north += n
}

func (s *Ship) South(n int) {
	s.North(-n)
}

func (s *Ship) Right(n int) {
	i := (n/90)%len(s.directions()) + len(s.directions())

	s.direction = (s.direction + i) % len(s.directions())
}

func (s *Ship) Left(n int) {
	s.Right(-n)
}

func (s *Ship) Forward(n int) {
	s.directions()[s.direction](n)
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	ship := NewShip()
	wship := NewWaypointShip()
	for s.Scan() {
		l := s.Text()
		instruction := l[0:1]
		number, err := strconv.Atoi(l[1:])
		if err != nil {
			log.Printf("error strconv: %v", err)
			return
		}

		ship.Do(instruction, number)
		wship.Do(instruction, number)
	}
	if s.Err() != nil {
		log.Printf("error scan: %v", s.Err())
		return
	}

	log.Printf("pt(1): %v", ship.Manhattan())
	log.Printf("pt(2): %v", wship.Manhattan())
}
