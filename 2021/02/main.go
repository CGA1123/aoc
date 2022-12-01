package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	horizontal := int64(0)
	depth := int64(0)

	aoc.EachLine("input.txt", func(line string) {
		splits := strings.SplitN(line, " ", 2)
		cmd, arg := splits[0], aoc.MustParse(splits[1])

		switch cmd {
		case "forward":
			horizontal += arg
		case "down":
			depth += arg
		case "up":
			depth -= arg
		}
	})

	fmt.Println(horizontal * depth)

	horizontal2 := int64(0)
	depth2 := int64(0)
	aim := int64(0)

	aoc.EachLine("input.txt", func(line string) {
		splits := strings.SplitN(line, " ", 2)
		cmd, arg := splits[0], aoc.MustParse(splits[1])

		switch cmd {
		case "forward":
			horizontal2 += arg
			depth2 += (arg * aim)
		case "down":
			aim += arg
		case "up":
			aim -= arg
		}
	})

	fmt.Println(horizontal2 * depth2)
}
