package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

type Hand int

const (
	Rock Hand = iota
	Paper
	Scissors
)

type Result int

const (
	Win Result = iota
	Draw
	Loss
)

var handScore = map[Hand]int64{
	Rock:     1,
	Paper:    2,
	Scissors: 3,
}

var beatenBy = map[Hand]Hand{
	Rock:     Paper,
	Paper:    Scissors,
	Scissors: Rock,
}

var beats = map[Hand]Hand{
	Rock:     Scissors,
	Paper:    Rock,
	Scissors: Paper,
}

func Score(a, b Hand) int64 {
	hs := handScore[b]
	if a == b {
		return 3 + hs
	}

	if beatenBy[a] == b {
		return 6 + hs
	}

	return hs
}

func HandFromString(s string) Hand {
	switch s {
	case "A", "X":
		return Rock
	case "B", "Y":
		return Paper
	case "C", "Z":
		return Scissors
	default:
		panic("bad hand " + s)
	}
}

func ResultFromString(s string) Result {
	switch s {
	case "X":
		return Win
	case "Y":
		return Draw
	case "Z":
		return Loss
	default:
		panic("bad result " + s)
	}
}

func HandFromResult(a Hand, r Result) Hand {
	switch r {
	case Draw:
		return a
	case Win:
		return beats[a]
	case Loss:
		return beatenBy[a]
	default:
		panic("unknown result")
	}
}

func main() {
	var score int64
	aoc.EachLine("input.txt", func(s string) {
		hands := strings.Split(s, " ")

		them, me := HandFromString(hands[0]), HandFromString(hands[1])

		score = score + Score(them, me)
	})

	fmt.Printf("%v\n", score)

	score = 0
	aoc.EachLine("input.txt", func(s string) {
		hands := strings.Split(s, " ")

		them := HandFromString(hands[0])
		need := ResultFromString(hands[1])
		me := HandFromResult(them, need)
		score = score + Score(them, me)
	})

	fmt.Printf("%v\n", score)
}
