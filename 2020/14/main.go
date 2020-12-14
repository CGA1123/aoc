package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var insructionRegex = regexp.MustCompile(`^(?P<instruction>[a-z]+)(\[(?P<address>.*)\])? = (?P<value>[0-9X]+)$`)
var names = insructionRegex.SubexpNames()

const (
	ZeroMask = uint64(0)
	OneMask  = uint64(^ZeroMask)
)

func capture(s string) map[string]string {
	m := map[string]string{}
	for _, match := range insructionRegex.FindAllStringSubmatch(s, -1) {
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

func mask(s string) (uint64, uint64, error) {
	or, err := strconv.ParseUint(strings.ReplaceAll(s, "X", "0"), 2, 64)
	if err != nil {
		return 0, 0, err
	}

	and, err := strconv.ParseUint(strings.ReplaceAll(s, "X", "1"), 2, 64)
	if err != nil {
		return 0, 0, err
	}

	return ZeroMask | or, OneMask & and, nil
}

type PartOne struct {
	or  uint64
	and uint64
	mem map[uint64]uint64
}

func NewPartOne() *PartOne {
	return &PartOne{mem: map[uint64]uint64{}}
}

func (pt *PartOne) Do(s string) error {
	result := capture(s)
	if result["instruction"] == "mask" {
		or, and, err := mask(result["value"])
		if err != nil {
			return err
		}

		pt.or = or
		pt.and = and
	} else {
		address, err := strconv.ParseUint(result["address"], 10, 64)
		if err != nil {
			return err
		}

		value, err := strconv.ParseUint(result["value"], 10, 64)
		if err != nil {
			return err
		}

		pt.mem[address] = (value | pt.or) & pt.and
	}

	return nil
}

func (pt *PartOne) Sum() uint64 {
	var total uint64
	for _, i := range pt.mem {
		total += i
	}

	return total
}

type PartTwo struct {
	mem  map[uint64]uint64
	mask []byte
}

func NewPartTwo() *PartTwo {
	return &PartTwo{mem: map[uint64]uint64{}}
}

func floating(mask []byte, address uint64, i int) []uint64 {
	var addresses []uint64

	if i == len(mask) {
		return []uint64{address}
	}

	bit := uint64(1 << (35 - i))
	one := address | bit
	zero := address & (^bit)

	switch mask[i] {
	case 'X':
		addresses = append(addresses, floating(mask, one, i+1)...)
		addresses = append(addresses, floating(mask, zero, i+1)...)
	case '1':
		addresses = append(addresses, floating(mask, one, i+1)...)
	case '0':
		addresses = append(addresses, floating(mask, address, i+1)...)
	}

	return addresses
}

func (pt *PartTwo) Do(s string) error {
	result := capture(s)
	if result["instruction"] == "mask" {
		pt.mask = []byte(result["value"])
	} else {
		address, err := strconv.ParseUint(result["address"], 10, 64)
		if err != nil {
			return err
		}

		value, err := strconv.ParseUint(result["value"], 10, 64)
		if err != nil {
			return err
		}

		for _, add := range floating(pt.mask, address, 0) {
			pt.mem[add] = value
		}
	}

	return nil
}

func (pt *PartTwo) Sum() uint64 {
	var total uint64
	for _, i := range pt.mem {
		total += i
	}

	return total
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	one := NewPartOne()
	two := NewPartTwo()

	for s.Scan() {
		l := s.Text()

		one.Do(l)
		two.Do(l)
	}
	if s.Err() != nil {
		log.Printf("error scanning: %v", s.Err())
		return
	}

	log.Printf("pt(1): %v", one.Sum())
	log.Printf("pt(2): %v", two.Sum())
}
