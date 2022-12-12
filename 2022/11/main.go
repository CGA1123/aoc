package main

import (
	"fmt"
	"sort"
)

type monkey struct {
	items     []int64
	operation func(int64) int64
	test      func(int64) int
	inspected int64
	mod       int64
}

func divisible(by int64, t, f int) func(int64) int {
	return func(i int64) int {
		if (i % by) == 0 {
			return t
		}

		return f
	}
}

var inputMonkeys = []*monkey{
	{
		items: []int64{92, 73, 86, 83, 65, 51, 55, 93},
		mod:   11,
		test:  divisible(11, 3, 4),
		operation: func(i int64) int64 {
			return i * 5
		},
	},
	{
		items: []int64{99, 67, 62, 61, 59, 98},
		test:  divisible(2, 6, 7),
		mod:   2,
		operation: func(i int64) int64 {
			return i * i
		},
	},
	{
		mod:   5,
		items: []int64{81, 89, 56, 61, 99},
		test:  divisible(5, 1, 5),
		operation: func(i int64) int64 {
			return i * 7
		},
	},
	{
		items: []int64{97, 74, 68},
		test:  divisible(17, 2, 5),
		mod:   17,
		operation: func(i int64) int64 {
			return i + 1
		},
	},
	{
		mod:   19,
		items: []int64{78, 73},
		test:  divisible(19, 2, 3),
		operation: func(i int64) int64 {
			return i + 3
		},
	},
	{
		items: []int64{50},
		test:  divisible(7, 1, 6),
		mod:   7,
		operation: func(i int64) int64 {
			return i + 5
		},
	},
	{
		items: []int64{95, 88, 53, 75},
		test:  divisible(3, 0, 7),
		mod:   3,
		operation: func(i int64) int64 {
			return i + 8
		},
	},
	{
		items: []int64{50, 77, 98, 85, 94, 56, 89},
		mod:   13,
		test:  divisible(13, 4, 0),
		operation: func(i int64) int64 {
			return i + 2
		},
	},
}

var exampleMonkeys = []*monkey{
	{
		items: []int64{79, 98},
		test:  divisible(23, 2, 3),
		mod:   23,
		operation: func(i int64) int64 {
			return i * 19
		},
	},
	{
		items: []int64{54, 65, 75, 74},
		test:  divisible(19, 2, 0),
		mod:   19,
		operation: func(i int64) int64 {
			return i + 6
		},
	},
	{
		items: []int64{79, 60, 97},
		test:  divisible(13, 1, 3),
		mod:   13,
		operation: func(i int64) int64 {
			return i * i
		},
	},
	{
		items: []int64{74},
		test:  divisible(17, 0, 1),
		mod:   17,
		operation: func(i int64) int64 {
			return i + 3
		},
	},
}

var monkeys = inputMonkeys

func printInspection() {
	for i, monkey := range monkeys {
		fmt.Printf("Monkey %v inspected items %v times.\n", i, monkey.inspected)
	}
}

func main() {
	mod := int64(1)
	for _, monkey := range monkeys {
		mod = mod * monkey.mod
	}

	fmt.Printf("mod: %v\n", mod)
	rounds := 10000
	for i := 0; i < rounds; i++ {
		if i == 1 || i == 20 || ((i % 1000) == 0) {
			fmt.Printf("== After round %v ==\n", i)
			printInspection()
		}

		for _, monkey := range monkeys {
			for _, item := range monkey.items {
				newValue := monkey.operation(item) % mod
				dest := monkeys[monkey.test(newValue)]
				dest.items = append(dest.items, newValue)
				monkey.inspected++
			}

			monkey.items = []int64{}
		}

	}

	fmt.Printf("== After round %v ==\n", rounds)
	printInspection()

	counts := make([]int64, len(monkeys))
	for i, monkey := range monkeys {
		counts[i] = monkey.inspected
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})

	fmt.Printf("%v * %v: %v\n", counts[0], counts[1], counts[0]*counts[1])
}
