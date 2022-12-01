package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type board struct {
	cells   []*cell
	mapping map[int64]int
}

func newBoard(nums []int64) *board {
	b := &board{
		cells:   make([]*cell, 25),
		mapping: map[int64]int{},
	}

	for i, n := range nums {
		b.cells[i] = &cell{value: n, marked: false}
		b.mapping[n] = i
	}

	return b
}

func (b *board) mark(val int64) bool {
	if i, ok := b.mapping[val]; ok {
		b.cells[i].marked = true

		return b.isWin(i)
	}

	return false
}

func (b *board) isWin(idx int) bool {
	col := idx % 5
	row := idx / 5

	rowWin := true
	colWin := true

	for i := 0; i < 5; i++ {
		rowIdx := (row * 5) + i
		colIdx := (5 * i) + col

		colWin = colWin && b.cells[colIdx].marked
		rowWin = rowWin && b.cells[rowIdx].marked
	}

	return colWin || rowWin
}

func (b *board) score() int64 {
	s := int64(0)
	for _, c := range b.cells {
		if c.marked {
			continue
		}

		s = s + c.value
	}

	return s
}

type cell struct {
	value  int64
	marked bool
}

func main() {
	lines := []string{}
	aoc.EachLine("input.txt", func(line string) {
		lines = append(lines, line)
	})

	numbers := []int64{}
	for _, n := range strings.Split(lines[0], ",") {
		numbers = append(numbers, aoc.MustParse(n))
	}

	boardNums := []int64{}
	for i := 2; i < len(lines); i++ {
		for _, f := range strings.Fields(lines[i]) {
			boardNums = append(boardNums, aoc.MustParse(f))
		}
	}

	numberOfBoards := len(boardNums) / 25
	boards := make([]*board, numberOfBoards)
	for i := 0; i < numberOfBoards; i++ {
		offset := (i * 25)
		boards[i] = newBoard(boardNums[offset : offset+25])
	}

	winners := map[int]int64{}
	winOrder := []int{}

	for _, n := range numbers {
		for i, b := range boards {
			if _, ok := winners[i]; ok {
				continue
			}

			if b.mark(n) {
				winners[i] = b.score() * n
				winOrder = append(winOrder, i)
			}
		}
	}

	loser := winOrder[len(winOrder)-1]
	winner := winOrder[0]
	log.Printf("winner: %v, score: %v", winner, winners[winner])
	log.Printf("loser: %v, score: %v", loser, winners[loser])
}
