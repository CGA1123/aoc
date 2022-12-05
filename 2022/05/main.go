package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

type element struct {
	previous *element
	value    string
}

type stack struct {
	head *element
	size int
}

func NewStack() *stack {
	return &stack{}
}

func (s *stack) Pop() string {
	v := s.head.value

	s.head = s.head.previous
	s.size = s.size - 1

	return v
}

func (s *stack) Push(el string) {
	n := &element{
		value:    el,
		previous: s.head,
	}

	s.head = n
	s.size = s.size + 1
}

func (s *stack) Size() int {
	return s.size
}

func (s *stack) Peek() string {
	return s.head.value
}

type move struct {
	from   int64
	to     int64
	amount int64
}

func main() {
	var parseMoves bool
	moves := []move{}
	stackLines := []string{}

	aoc.EachLine("input.txt", func(s string) {
		if s == "" {
			parseMoves = true
			return
		}

		if parseMoves {
			parts := strings.Split(s, " ")

			qty, from, to := aoc.MustParse(parts[1]), aoc.MustParse(parts[3]), aoc.MustParse(parts[5])

			moves = append(moves, move{from: from - 1, to: to - 1, amount: qty})

			return
		}

		stackLines = append(stackLines, s)
	})

	stacks9000 := parseStacks(stackLines)
	for _, m := range moves {
		execMove9000(stacks9000, m)
	}

	for _, s := range stacks9000 {
		fmt.Print(s.Peek())
	}
	fmt.Println("")

	stacks9001 := parseStacks(stackLines)
	for _, m := range moves {
		execMove9001(stacks9001, m)
	}

	for _, s := range stacks9001 {
		fmt.Print(s.Peek())
	}
	fmt.Println("")
}

func execMove9000(stacks []*stack, m move) {
	for i := 0; i < int(m.amount); i++ {
		el := stacks[m.from].Pop()
		stacks[m.to].Push(el)
	}
}

func execMove9001(stacks []*stack, m move) {
	containers := make([]string, m.amount)
	for i := int(m.amount) - 1; i >= 0; i-- {
		containers[i] = stacks[m.from].Pop()
	}

	for _, c := range containers {
		stacks[m.to].Push(c)
	}
}

func parseStacks(lines []string) []*stack {
	cols := (len(lines[0]) + 1) / 4
	stacks := make([]*stack, cols)
	for i := 0; i < cols; i++ {
		stacks[i] = NewStack()
	}

	for i := len(lines) - 2; i >= 0; i-- {
		for col := 0; col < cols; col++ {
			idx := (col * 4) + 1
			char := lines[i][idx : idx+1]
			if char == " " {
				continue
			}

			stacks[col].Push(char)
		}
	}

	return stacks
}
