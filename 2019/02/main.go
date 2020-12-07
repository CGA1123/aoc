package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func compute(mem []int) int {
	var pc int
	for {
		switch mem[pc] {
		case 1:
			mem[mem[pc+3]] = mem[mem[pc+1]] + mem[mem[pc+2]]
			pc = pc + 4
		case 2:
			mem[mem[pc+3]] = mem[mem[pc+1]] * mem[mem[pc+2]]
			pc = pc + 4
		case 99:
			return mem[0]
		}
	}
}

func main() {
	var orig []int

	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer f.Close()

	s, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	for _, number := range strings.Split(string(s[:len(s)-1]), ",") {
		i, err := strconv.Atoi(number)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		orig = append(orig, i)
	}

	mem := make([]int, len(orig))
	copy(mem, orig)

	mem[1], mem[2] = 12, 2
	log.Printf("pt1: %v", compute(mem))

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			copy(mem, orig)

			mem[1], mem[2] = i, j

			if compute(mem) == 19690720 {
				log.Printf("pt2: %v", i*100+j)
			}
		}
	}
}
