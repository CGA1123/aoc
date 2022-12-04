package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	var containCount, overlapCount int
	aoc.EachLine("input.txt", func(s string) {
		a := strings.Split(s, ",")
		aL, aR := strings.Split(a[0], "-"), strings.Split(a[1], "-")

		fromL, toL := aoc.MustParse(aL[0]), aoc.MustParse(aL[1])
		fromR, toR := aoc.MustParse(aR[0]), aoc.MustParse(aR[1])

		betweenL := aoc.Between(fromL, toL)
		betweenR := aoc.Between(fromR, toR)

		if (betweenR(fromL) && betweenR(toL)) || (betweenL(fromR) && betweenL(toR)) {
			containCount++
		}
		if betweenR(fromL) || betweenR(toL) || betweenL(fromR) || betweenL(toR) {
			overlapCount++
		}

	})

	fmt.Printf("%v\n", containCount)
	fmt.Printf("%v\n", overlapCount)
}
