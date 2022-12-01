package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CGA1123/aoc"
)

func main() {
	var counts []int
	var bitlen int
	lines := []string{}
	aoc.EachLine("input.txt", func(line string) {
		lines = append(lines, line)

		if counts == nil {
			counts = make([]int, len(line))
			bitlen = len(line)
		}

		for i, c := range strings.Split(line, "") {
			switch c {
			case "1":
				counts[i] += 1
			case "0":
				counts[i] -= 1
			}
		}
	})

	binStr := ""
	for _, c := range counts {
		if c >= 0 {
			binStr += "1"
		} else {
			binStr += "0"
		}
	}

	mask := uint64((1 << bitlen) - 1)

	bin, _ := strconv.ParseUint(binStr, 2, 64)
	ibin := (^bin & mask)

	fmt.Printf("%v\n", bin*ibin)

	// PART 2
	cols := len(lines[0])
	co2 := lines
	oxygen := lines

	for i := 0; i < cols; i++ {
		oxygen = getRating(oxygen, i, true)
		if len(oxygen) == 1 {
			break
		}
	}
	for i := 0; i < cols; i++ {
		co2 = getRating(co2, i, false)
		if len(co2) == 1 {
			break
		}
	}

	co2Rating, _ := strconv.ParseUint(co2[0], 2, 64)
	oxygenRating, _ := strconv.ParseUint(oxygen[0], 2, 64)

	fmt.Printf("%v\n", co2Rating*oxygenRating)

}

func getRating(lines []string, col int, isOxygen bool) []string {
	res := []string{}
	i := 0

	for _, line := range lines {
		if line[col] == '0' {
			i = i - 1
		} else {
			i = i + 1
		}
	}

	var target byte
	if isOxygen {
		// 1 is equal or most common
		if i >= 0 {
			target = '1'
		} else {
			target = '0'
		}
	} else {
		if i >= 0 {
			target = '0'
		} else {
			target = '1'
		}
	}

	for _, line := range lines {
		if line[col] == target {
			res = append(res, line)
		}
	}

	return res
}
