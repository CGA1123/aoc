package main

import (
	"bufio"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func timestamp(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

type Bus struct {
	x int64
	n int64
}

func buses(s string) ([]Bus, error) {
	var i []Bus

	for offset, bus := range strings.Split(s, ",") {
		if bus == "x" {
			continue
		}

		b, err := strconv.Atoi(bus)
		if err != nil {
			return nil, err
		}

		i = append(i, Bus{n: int64(b), x: int64(offset)})
	}

	return i, nil
}

func PartOne(timestamp int64, buses []Bus) int64 {
	var waiting int64
	var onbus int64

	for _, bus := range buses {
		val := (bus.n - (timestamp % bus.n))
		if waiting == 0 || val < waiting {
			waiting = val
			onbus = bus.n
		}
	}

	return waiting * onbus
}

func product(buses []Bus, f func(Bus) int64) int64 {
	p := int64(1)
	for _, b := range buses {
		p *= f(b)
	}

	return p
}

func inverseModulo(g, n int64) int64 {
	bigN := big.NewInt(n)
	bigG := big.NewInt(g)

	return (&big.Int{}).ModInverse(bigG, bigN).Int64()
}

func PartTwo(buses []Bus) int64 {
	productOfModulos := product(buses, func(b Bus) int64 { return int64(b.n) })

	var t int64
	for _, bus := range buses {
		nullifyingCoeff := (productOfModulos / bus.n)
		offsetCoeff := inverseModulo(nullifyingCoeff, bus.n) * (bus.n - bus.x)

		t += nullifyingCoeff * offsetCoeff
	}

	return t % productOfModulos
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		log.Printf("error scanning: %v", s.Err())
		return
	}

	timestamp, err := timestamp(lines[0])
	if err != nil {
		log.Printf("error timestamp: %v", err)
	}

	b, err := buses(lines[1])
	if err != nil {
		log.Printf("error buses: %v", err)
	}

	log.Printf("pt(1): %v", PartOne(timestamp, b))
	log.Printf("pt(2): %v", PartTwo(b))
}
