package main

import (
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type Scratchcard struct {
	ID      int64
	Numbers []int64
	Winners []int64
	Score   int64
	Count   int64
}

func parseNumbers(in string) []int64 {
	out := []int64{}

	for _, n := range strings.Split(in, " ") {
		nn := strings.TrimSpace(n)
		if nn == "" {
			continue
		}

		out = append(out, aoc.MustParse(nn))
	}

	return out
}

func parseScratchCard(line string) *Scratchcard {
	idNumbers := strings.Split(line, ": ")

	id := aoc.MustParse(strings.TrimSpace(strings.TrimPrefix(idNumbers[0], "Card ")))

	numbersGroups := strings.Split(idNumbers[1], " | ")

	return &Scratchcard{
		ID:      id,
		Numbers: parseNumbers(numbersGroups[0]),
		Winners: parseNumbers(numbersGroups[1]),
	}
}

func (sc *Scratchcard) Compute() {
	winners := []int64{}
	wins := map[int64]struct{}{}

	for _, w := range sc.Winners {
		wins[w] = struct{}{}
	}

	for _, n := range sc.Numbers {
		if _, ok := wins[n]; ok {
			winners = append(winners, n)
		}
	}
	if len(winners) == 0 {
		sc.Score = 0
		return
	}

	score := int64(1)

	for i := 0; i < len(winners)-1; i++ {
		score *= 2
	}

	sc.Score = score
	sc.Count = int64(len(winners))
}

func main() {
	cards := map[int64]*Scratchcard{}

	aoc.EachLine("input.txt", func(s string) {
		c := parseScratchCard(s)
		cards[c.ID] = c
	})

	var total int64
	for _, c := range cards {
		c.Compute()
		total += c.Score
	}

	log.Printf("total: %v", total)

	totalScore := int64(len(cards))
	scores := make(map[int64]int64, len(cards))
	for i := int64(1); i <= int64(len(cards)); i++ {
		scores[i] = computeScore(cards, scores, i)
		totalScore += scores[i]
	}

	log.Printf("total: %v", totalScore)
}

func computeScore(cards map[int64]*Scratchcard, scores map[int64]int64, id int64) int64 {
	if score, ok := scores[id]; ok {
		return score
	}

	score := cards[id].Count

	for i := int64(0); i < cards[id].Count; i++ {
		score += computeScore(cards, scores, id+i+1)
	}

	return score
}
