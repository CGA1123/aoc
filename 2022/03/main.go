package main

import (
	"fmt"

	"github.com/CGA1123/aoc"
)

func main() {
	var part1 int
	var part2 int
	idx := 0
	group := make([]string, 3)
	aoc.EachLine("input.txt", func(s string) {
		part1 = part1 + PartOne(s)
		group[idx] = s
		if idx == 2 {
			part2 = part2 + PartTwo(group)
		}

		idx = (idx + 1) % 3
	})

	fmt.Printf("%v\n", part1)
	fmt.Printf("%v\n", part2)

}

func PartTwo(group []string) int {
	var items *aoc.Set
	for _, s := range group {
		groupItems := aoc.NewSet()

		for _, b := range []byte(s) {
			groupItems.Add(b)
		}

		if items == nil {
			items = groupItems
		} else {
			items = aoc.Intersection(items, groupItems)
		}
	}

	return priority(items.Elements()[0].(byte))
}

func PartOne(s string) int {
	half := len(s) / 2
	firstHalf := map[byte]int{}

	for i, b := range []byte(s) {
		if i < half {
			firstHalf[b]++
			continue
		}

		if _, ok := firstHalf[b]; ok {
			return priority(b)
		}
	}

	panic("fail part2")
}

func priority(b byte) int {
	if isUpper(b) {
		return int(b) - 65 + 26 + 1
	} else {
		return int(b) - 97 + 1
	}
}

func isUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}
