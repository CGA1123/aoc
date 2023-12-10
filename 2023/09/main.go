package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	lines := [][]int64{}
	aoc.EachLine("input.txt", func(s string) {
		line := []int64{}
		for _, i := range strings.Split(s, " ") {
			line = append(line, aoc.MustParse(i))
		}
		lines = append(lines, line)
	})

	var totalNext, totalPrevious int64
	for _, line := range lines {
		totalNext += nextValue(line)
		totalPrevious += previousValue(line)
	}

	log.Printf("totalNext: %v", totalNext)
	log.Printf("totalPrevious: %v", totalPrevious)
}

func nextValue(seq []int64) int64 {
	if allZero(seq) {
		return 0
	}

	return seq[len(seq)-1] + nextValue(derive(seq))
}

func derive(seq []int64) []int64 {
	d := make([]int64, len(seq)-1)
	for i := 0; i < len(seq)-1; i++ {
		d[i] = seq[i+1] - seq[i]
	}

	return d
}

func previousValue(seq []int64) int64 {
	if allZero(seq) {
		return 0
	}

	return seq[0] - previousValue(derive(seq))
}

func allZero(seq []int64) bool {
	for _, s := range seq {
		if s != 0 {
			return false
		}
	}

	return true
}
