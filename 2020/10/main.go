package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

func counter(a []int, memo map[int]int, from int) (int, map[int]int) {
	if from == len(a)-1 {
		return 1, memo
	}

	var count int
	for i := 1; i <= 3; i++ {
		idx := from + i
		if idx >= len(a) || a[idx]-a[from] > 3 {
			continue
		}

		var currentCount int
		if c, ok := memo[idx]; ok {
			currentCount = c
		} else {
			currentCount, memo = counter(a, memo, idx)
			memo[idx] = currentCount
		}

		count = count + currentCount
	}

	return count, memo
}

func PartOne(adapters []int) int {
	diffs := map[int]int{1: 1, 3: 1}

	for i := 0; i < len(adapters)-1; i++ {
		diff := adapters[i+1] - adapters[i]
		diffs[diff] = diffs[diff] + 1
	}

	return diffs[1] * diffs[3]
}

func PartTwo(adapters []int) int {
	var current int
	var count int

	memo := map[int]int{}
	for i := 0; i < 3; i++ {
		if adapters[i] > 3 {
			continue
		}

		current, memo = counter(adapters, memo, i)
		count = count + current
	}

	return count
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	var adapters []int
	for s.Scan() {
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		adapters = append(adapters, i)
	}
	if s.Err() != nil {
		log.Printf("error: %v", s.Err())
		return
	}

	sort.Ints(adapters)

	log.Printf("pt(1): %v", PartOne(adapters))
	log.Printf("pt(2): %v", PartTwo(adapters))
}
