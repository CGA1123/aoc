package main

import (
	"log"
	"strings"
	"unicode"

	"github.com/CGA1123/aoc"
)

func main() {
	g := map[string][]string{}
	aoc.EachLine("input.txt", func(s string) {
		path := strings.Split(s, "-")
		from, to := path[0], path[1]

		g[from] = append(g[from], to)
		g[to] = append(g[to], from)
	})

	log.Printf("%v", traverse("start", g, map[string]struct{}{}, 0, true))
	log.Printf("%v", traverse("start", g, map[string]struct{}{}, 0, false))
}

func traverse(node string, g map[string][]string, v map[string]struct{}, d int, double bool) int {
	// check if we've already seen this node.
	if _, ok := v[node]; ok {
		if node == "start" {
			return 0
		}

		if double {
			return 0
		} else {
			double = true
		}
	}

	if node == "end" {
		return 1
	}

	// mark this node as visited if it's a small cave
	char := []rune(node)[0]
	if unicode.IsLower(char) {
		v[node] = struct{}{}
	}

	i := 0
	for _, next := range g[node] {
		vn := make(map[string]struct{}, len(v))
		for visited := range v {
			vn[visited] = struct{}{}
		}

		i = i + traverse(next, g, vn, d+1, double)
	}

	return i
}
