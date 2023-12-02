package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type Game struct {
	ID     int64
	Rounds []Round
}

type Round struct {
	Red   int64
	Blue  int64
	Green int64
}

func ParseGame(s string) Game {
	gameAndRounds := strings.Split(s, ": ")
	rounds := []Round{}

	for _, roundStr := range strings.Split(gameAndRounds[1], "; ") {
		round := Round{}
		for _, countColorStr := range strings.Split(roundStr, ", ") {
			countColor := strings.Split(countColorStr, " ")
			count := aoc.MustParse(countColor[0])
			switch countColor[1] {
			case "blue":
				round.Blue = count
			case "red":
				round.Red = count
			case "green":
				round.Green = count
			}
		}

		rounds = append(rounds, round)
	}

	return Game{
		ID:     aoc.MustParse(strings.TrimPrefix(gameAndRounds[0], "Game ")),
		Rounds: rounds,
	}
}

func RoundPossible(maxes Round, round Round) bool {
	return round.Red <= maxes.Red && round.Green <= maxes.Green && round.Blue <= maxes.Blue
}

func (game Game) Minimum() Round {
	round := Round{}

	for _, r := range game.Rounds {
		if round.Blue < r.Blue {
			round.Blue = r.Blue
		}
		if round.Green < r.Green {
			round.Green = r.Green
		}
		if round.Red < r.Red {
			round.Red = r.Red
		}
	}

	return round
}

func (game Game) Power() int64 {
	min := game.Minimum()

	return min.Red * min.Blue * min.Green
}

func (game Game) Possible(maxes Round) bool {
	for _, round := range game.Rounds {
		if !RoundPossible(maxes, round) {
			return false
		}
	}

	return true
}

func main() {
	maxes := Round{Red: 12, Green: 13, Blue: 14}
	var total int64
	var totalPower int64
	aoc.EachLine("input.txt", func(s string) {
		game := ParseGame(s)

		if game.Possible(maxes) {
			total += game.ID
		}

		totalPower += game.Power()
	})

	log.Printf("total: %v", total)
	log.Printf("totalPower: %v", totalPower)
}
