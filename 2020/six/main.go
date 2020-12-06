package main

import (
	"bufio"
	"log"
	"math/bits"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	answers := map[int][]uint32{}
	var group int

	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 {
			group = group + 1
			continue
		}

		var answer uint32
		for _, b := range line {
			offset := b - 'a'

			answer = answer | uint32(1<<offset)
		}

		answers[group] = append(answers[group], answer)
	}
	if s.Err() != nil {
		log.Printf("error: %v", s.Err())
		return
	}

	var ors int
	var ands int

	for _, a := range answers {
		or := uint32(0)
		and := ^uint32(0)

		for _, i := range a {
			or = or | i
			and = and & i
		}

		ors = ors + bits.OnesCount32(or)
		ands = ands + bits.OnesCount32(and)
	}

	log.Printf("one: %v", ors)
	log.Printf("two: %v", ands)
}
