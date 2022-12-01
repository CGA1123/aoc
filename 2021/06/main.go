package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	fish := map[int64]int64{}
	aoc.EachLine("input.txt", func(l string) {
		for _, f := range strings.Split(l, ",") {
			n := aoc.MustParse(f)
			fish[n] = fish[n] + 1
		}
	})

	for i := 0; i < 256; i++ {
		fish = tick(fish)
	}

	fmt.Printf("%v\n", total(fish))
}

func total(fish map[int64]int64) int64 {
	t := int64(0)
	for _, n := range fish {
		t = t + n
	}

	return t
}

func tick(fish map[int64]int64) map[int64]int64 {
	next := map[int64]int64{}

	for t, n := range fish {
		if t == 0 {
			next[int64(8)] = n
			next[int64(6)] = next[int64(6)] + n
			continue
		}

		next[t-1] = next[t-1] + n
	}

	return next
}
