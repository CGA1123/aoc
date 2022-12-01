package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

// uniques is a mapping of the number of segments to the unique number the must
// display
var uniques = map[int]int{
	2: 1,
	3: 7,
	4: 4,
	7: 8,
}

func solve(s string) int {
	allNumbers := strings.Split(strings.Split(s, " | ")[0], " ")
	for i, n := range allNumbers {
		s := strings.Split(n, "")
		sort.Strings(s)
		allNumbers[i] = strings.Join(s, "")
	}

	solution := map[string]int{}
	numbers := map[int]string{}
	twoThreeFive := []string{}
	zeroSixNine := []string{}

	for _, number := range allNumbers {
		l := len(number)
		if i, ok := uniques[l]; ok {
			numbers[i] = number
			solution[number] = i
			continue
		}
		if l == 5 {
			twoThreeFive = append(twoThreeFive, number)
		}
		if l == 6 {
			zeroSixNine = append(zeroSixNine, number)
		}
	}
	for _, v := range twoThreeFive {
		if len(minus(v, numbers[1])) == 3 {
			numbers[3] = v
			solution[v] = 3
			continue
		}

		if len(minus(v, numbers[4])) == 3 {
			numbers[2] = v
			solution[v] = 2
		} else {
			numbers[5] = v
			solution[v] = 5
		}
	}
	for _, v := range zeroSixNine {
		if len(minus(v, numbers[1])) == 5 {
			solution[v] = 6
			numbers[6] = v
			continue
		}
		if len(minus(v, numbers[4])) == 2 {
			solution[v] = 9
			numbers[9] = v
		} else {
			solution[v] = 0
			numbers[0] = v
		}
	}

	answer := strings.Split(strings.Split(s, " | ")[1], " ")
	for i, n := range answer {
		s := strings.Split(n, "")
		sort.Strings(s)
		answer[i] = strings.Join(s, "")
	}

	number := 0
	for i, n := range answer {
		number = number + (solution[n] * pow(10, 3-i))
	}

	return number
}

func pow(a, b int) int {
	n := 1
	for i := 0; i < b; i++ {
		n = n * a
	}

	return n
}

func minus(a, b string) string {
	s := a

	for _, c := range strings.Split(b, "") {
		s = strings.ReplaceAll(s, c, "")
	}

	return s
}

func main() {
	c := 0
	aoc.EachLine("input.txt", func(s string) {
		result := strings.Split(s, " | ")[1]
		for _, num := range strings.Fields(result) {
			l := len(num)
			if l == 4 || l == 2 || l == 3 || l == 7 {
				c = c + 1
			}
		}
	})

	fmt.Printf("%v\n", c)

	total := 0
	aoc.EachLine("input.txt", func(s string) {
		total = total + solve(s)
	})
	fmt.Printf("%v\n", total)
}
