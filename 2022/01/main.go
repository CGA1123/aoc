package main

import (
	"fmt"
	"sort"

	"github.com/CGA1123/aoc"
)

func main() {
	elfs := []int64{}
	var currentElf int64
	aoc.EachLine("input.txt", func(s string) {
		if s == "" {
			elfs = append(elfs, currentElf)
			currentElf = 0
			return
		}

		currentElf = currentElf + aoc.MustParse(s)
	})

	sort.Slice(elfs, func(i, j int) bool {
		return elfs[j] < elfs[i]
	})

	fmt.Printf("%v\n", elfs[0])

	var topThree int64
	for i := 0; i < 3; i++ {
		topThree = topThree + elfs[i]
	}

	fmt.Printf("%v\n", topThree)
}
