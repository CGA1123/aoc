package main

import (
	"log"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	pos := map[int64]int64{}
	aoc.EachLine("input.txt", func(s string) {
		for _, n := range strings.Split(s, ",") {
			i := aoc.MustParse(n)
			pos[i] = pos[i] + 1
		}
	})

	keys := make([]int64, 0, len(pos))
	for k := range pos {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	minACost := int64(0)
	for i := keys[0]; i <= keys[len(keys)-1]; i++ {
		cost := int64(0)
		for p, n := range pos {
			cost = cost + (aoc.Abs(i-p) * n)
		}

		if minACost == 0 || cost < minACost {
			minACost = cost
		}
	}

	log.Printf("%v", minACost)

	minBCost := int64(0)
	for i := keys[0]; i <= keys[len(keys)-1]; i++ {
		cost := int64(0)
		for p, n := range pos {
			dist := aoc.Abs(i - p)
			moveCost := (dist * (dist + 1)) / 2
			cost = cost + (moveCost * n)
		}

		if minBCost == 0 || cost < minBCost {
			minBCost = cost
		}
	}

	log.Printf("%v", minBCost)
}
