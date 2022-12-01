package main

import (
	"fmt"

	"github.com/CGA1123/aoc"
)

func main() {
	measures := []int64{}
	previous := int64(0)
	counter := -1

	aoc.EachLine("input.txt", func(line string) {
		current := aoc.MustParse(line)
		measures = append(measures, current)

		if current > previous {
			counter++
		}

		previous = current
	})

	fmt.Println(counter)

	previousSum := measures[0] + measures[1] + measures[2]
	counterSum := 0

	for i := 3; i < len(measures); i++ {
		currentSum := measures[i] + measures[i-1] + measures[i-2]

		if currentSum > previousSum {
			counterSum++
		}

		previousSum = currentSum
	}

	fmt.Println(counterSum)
}
