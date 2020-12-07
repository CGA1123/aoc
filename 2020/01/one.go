package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func reduce(base int, f func(int, int) int) func([]int) int {
	return func(list []int) int {
		result := base

		for _, i := range list {
			result = f(result, i)
		}

		return result
	}
}

func sum(i []int) int {
	return reduce(0, func(a, b int) int {
		return a + b
	})(i)
}

func product(i []int) int {
	return reduce(1, func(a, b int) int {
		return a * b
	})(i)
}

func buildCombo(indices, list []int) []int {
	var result []int

	for _, i := range indices {
		result = append(result, list[i])
	}

	return result
}

func solve(list []int, n int) error {
	for _, comboIdx := range Combinations(len(list), n) {
		combo := buildCombo(comboIdx, list)
		if sum(combo) == 2020 {
			log.Printf("found pair [%v]", combo)
			log.Printf("product    [%v]", product(combo))

			return nil
		}
	}

	return errors.New("did not find solution!")
}

// Combination returns all combinations for n objects choosing r
// Objects are 0 indexed
func Combinations(n, r int) [][]int {
	var combinations func(int, int, int, []int) [][]int
	combinations = func(n, r, min int, prefix []int) [][]int {
		toBeSelected := r - len(prefix)

		// the buffer is full!
		if toBeSelected == 0 {
			return [][]int{prefix}
		}

		var combs [][]int

		for i := min; i <= (n - toBeSelected); i++ {
			combs = append(
				combs,
				combinations(n, r, i, append(prefix, i))...,
			)
		}

		return combs
	}

	return combinations(n, r, 0, []int{})
}

func readInt(r *bufio.Reader) (int, error) {
	str, err := r.ReadString('\n')
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(strings.TrimSuffix(str, "\n"))
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var list []int
	var i int
	var e error

	for {
		i, e = readInt(r)
		if e != nil {
			break
		}

		list = append(list, i)
	}
	if e != nil && e != io.EOF {
		log.Printf("reading input: %v", e)
		return
	}

	if err := solve(list, 2); err != nil {
		log.Printf("error: %v", err)
	}

	if err := solve(list, 3); err != nil {
		log.Printf("error: %v", err)
	}
}
