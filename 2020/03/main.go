package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func solve(forest [][]byte, gradients [][]int) []int {
	counts := make([]int, len(gradients))
	for y, line := range forest {
		for c, gradient := range gradients {
			dx, dy := gradient[0], gradient[1]

			if (y % dy) != 0 {
				continue
			}

			x := ((y / dy) * dx) % len(line)

			if line[x] == byte('#') {
				counts[c] = counts[c] + 1
			}
		}
	}

	return counts
}

func main() {
	f, e := os.Open("input.txt")
	if e != nil {
		log.Printf("opening file: %v", e)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)

	var forest [][]byte
	var err error

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			break
		}

		line = line[:len(line)-1]
		forest = append(forest, line)
	}
	if err != nil && err != io.EOF {
		log.Printf("error: %v", err)
		return
	}

	log.Printf("length %v", len(forest))

	problemsOne := [][]int{{3, 1}}

	log.Printf("count (pt. 1): %v", solve(forest, problemsOne)[0])

	problemsTwo := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2}}

	result := 1
	for _, count := range solve(forest, problemsTwo) {
		result = result * count
	}

	log.Printf("count (pt. 2): %v", result)
}
