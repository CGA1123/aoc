package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	rules := map[string]string{}

	var template string
	aoc.EachLine("example.txt", func(s string) {
		if template == "" {
			template = s
			return
		}

		if s == "" {
			return
		}

		pair := strings.Split(s, " -> ")

		rules[pair[0]] = pair[1]
	})

	for i := 0; i < 40; i++ {
		template = step(template, rules)
	}

	countsMap := map[rune]int{}
	for _, c := range template {
		countsMap[c] = countsMap[c] + 1
	}
	counts := make([]int, 0, len(countsMap))
	for _, c := range countsMap {
		counts = append(counts, c)
	}

	sort.Ints(counts)

	fmt.Printf("%v\n", counts[len(counts)-1]-counts[0])
}

func step(template string, rules map[string]string) string {
	t := template[0:1]
	for _, pair := range pairs(t) {
		t = t + rules[pair] + pair[1:]
	}

	return t
}

func pairs(str string) []string {
	p := []string{}
	for i := 0; i < len(str)-1; i++ {
		p = append(p, (str[i : i+2]))
	}

	return p
}
