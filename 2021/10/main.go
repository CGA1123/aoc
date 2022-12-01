package main

import (
	"log"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

var openingPairs = map[string]string{
	"[": "]",
	"(": ")",
	"{": "}",
	"<": ">",
}

var closingPairs = map[string]string{
	"]": "[",
	")": "(",
	"}": "{",
	">": "<",
}

var autoScores = map[string]int{
	")": 1,
	"]": 2,
	"}": 3,
	">": 4,
}

var scores = map[string]int{
	")": 3,
	"]": 57,
	"}": 1197,
	">": 25137,
}

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

func main() {
	scoreA := 0
	scoresB := []int{}
	aoc.EachLine("input.txt", func(line string) {
		s := NewStack()
		for _, c := range strings.Split(line, "") {
			if _, ok := openingPairs[c]; ok {
				s.Push(c)
				continue
			}

			closer := closingPairs[c]
			if closer == s.Peek() {
				s.Pop()
				continue
			}

			scoreA = scoreA + scores[c]

			return
		}

		// incomplete
		scoreB := 0
		size := s.Size()
		for i := 0; i < size; i++ {
			scoreB = scoreB * 5
			closer := openingPairs[s.Pop()]
			score := autoScores[closer]
			scoreB = scoreB + score
		}

		scoresB = append(scoresB, scoreB)
	})

	log.Printf("%v", scoreA)
	sort.Ints(scoresB)
	log.Printf("%v", scoresB[len(scoresB)/2])
}
