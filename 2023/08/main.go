package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

const (
	Left  = 0
	Right = 1
)

var Direction = map[string]int{
	"L": Left,
	"R": Right,
}

func main() {
	var instructions string
	tree := map[string][]string{}

	aoc.EachLine("input.txt", func(s string) {
		if instructions == "" {
			instructions = s
			return
		}
		if s == "" {
			return
		}

		parts := strings.Split(s, " = ")
		from := parts[0]
		leftRight := strings.Split(strings.TrimSuffix(strings.TrimPrefix(parts[1], "("), ")"), ", ")

		tree[from] = leftRight
	})

	i := 0
	current := "AAA"
	for {
		if current == "ZZZ" {
			log.Printf("%v", i)
			break
		}

		idx := i % len(instructions)
		direction := Direction[instructions[idx:idx+1]]
		current = tree[current][direction]

		i++
	}
}
