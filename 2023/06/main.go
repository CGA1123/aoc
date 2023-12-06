package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

func parseNumbers(in string) []int64 {
	out := []int64{}

	for _, n := range strings.Split(in, " ") {
		nn := strings.TrimSpace(n)
		if nn == "" {
			continue
		}

		out = append(out, aoc.MustParse(nn))
	}

	return out
}

func wins(t, r int64) int64 {
	var wins int64
	for i := int64(1); i < t; i++ {
		distance := (t - i) * i
		if distance > r {
			wins++
		}
	}

	return wins
}

func main() {
	r := map[string][]int64{}
	t := map[string]int64{}

	aoc.EachLine("input.txt", func(s string) {
		parts := strings.Split(s, ":  ")
		name := strings.TrimSpace(parts[0])
		r[name] = parseNumbers(parts[1])
		t[name] = aoc.MustParse(strings.ReplaceAll(parts[1], " ", ""))

	})

	total := int64(1)
	for i, t := range r["Time"] {
		total *= wins(t, r["Distance"][i])
	}

	log.Printf("%+v", total)
	log.Printf("%+v", wins(t["Time"], t["Distance"]))
}
