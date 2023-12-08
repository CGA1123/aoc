package main

import (
	"log"
	"sort"
	"strings"

	"github.com/CGA1123/aoc"
)

const part = 1

var cards = map[string]int64{
	"A": 14,
	"K": 13,
	"Q": 12,
	"J": 11,
	"T": 10,
	"9": 9,
	"8": 8,
	"7": 7,
	"6": 6,
	"5": 5,
	"4": 4,
	"3": 3,
	"2": 2,
}

type HandType func([]string) (int64, bool)

func tally(cards []string) (map[string]int64, map[int64]int64) {
	m := map[string]int64{}
	for _, c := range cards {
		m[c] = m[c] + 1
	}

	mm := map[int64]int64{}
	for _, i := range m {
		mm[i] = mm[i] + 1
	}

	return m, mm
}

var hands = []HandType{
	// five of a kind
	func(s []string) (int64, bool) {
		_, c := tally(s)
		if c[5] == 1 {
			return 7, true
		}

		return 0, false

	},
	// four of a kind
	//
	// J: 1 => five-of-a-kind
	// J: 4 => five-of-a-kind
	func(s []string) (int64, bool) {
		i, c := tally(s)
		if c[4] == 1 {
			if part == 2 {
				if i["J"] == 1 || i["J"] == 4 {
					return 7, true
				}
			}

			return 6, true
		}

		return 0, false
	},
	// full house
	func(s []string) (int64, bool) {
		i, c := tally(s)
		if c[3] == 1 && c[2] == 1 {
			if part == 2 {
				if i["J"] == 2 || i["J"] == 3 {
					return 7, true
				}
			}

			return 5, true
		}

		return 0, false
	},
	// three-of-a-kind
	//
	// J: 1 => four-of-a-kind
	// J: 2 => five-of-a-kind
	// J: 3 => four-of-a-kind
	func(s []string) (int64, bool) {
		i, c := tally(s)
		if c[3] == 1 {
			if part == 2 {
				if i["J"] == 1 {
					return 6, true
				}
				if i["J"] == 2 {
					return 7, true
				}
				if i["J"] == 3 {
					return 6, true
				}
			}

			return 4, true
		}

		return 0, false
	},
	// two-pair
	//
	// J: 1 => full-house
	// J: 2 => four-of-a-kind
	func(s []string) (int64, bool) {
		i, c := tally(s)
		if c[2] == 2 {
			if part == 2 {
				if i["J"] == 1 {
					return 5, true
				}
				if i["J"] == 2 {
					return 6, true
				}
			}

			return 3, true
		}

		return 0, false
	},
	// one-pair
	//
	// J: 1 => three-of-a-kind
	// J: 2 => three-of-a-kind
	func(s []string) (int64, bool) {
		i, c := tally(s)
		if c[2] == 1 {
			if part == 2 {
				if i["J"] == 1 {
					return 4, true
				}
				if i["J"] == 2 {
					return 4, true
				}
			}

			return 2, true
		}

		return 0, false
	},
	// high-card
	func(s []string) (int64, bool) {
		i, c := tally(s)
		if c[1] == 5 {
			if part == 2 {
				if i["J"] == 1 {
					return 2, true
				}
			}
			return 1, true
		}

		return 0, false
	},
}

type Hand struct {
	Bid   int64
	Cards []string
}

func (h *Hand) Type() int64 {
	for _, hand := range hands {
		if score, ok := hand(h.Cards); ok {
			return score
		}
	}

	return 0
}

func (h *Hand) Less(i *Hand) bool {
	hs, is := h.Type(), i.Type()
	if hs < is {
		return true
	}
	if hs > is {
		return false
	}

	for idx, c := range h.Cards {
		hc, ic := cards[c], cards[i.Cards[idx]]
		if hc == ic {
			continue
		}

		if c == "J" && part == 2 {
			hc -= 10
		}
		if i.Cards[idx] == "J" && part == 2 {
			ic -= 10
		}

		if hc < ic {
			return true
		} else {
			return false
		}
	}

	return false
}

func main() {
	var hands []*Hand
	aoc.EachLine("input.txt", func(s string) {
		parts := strings.Split(s, " ")

		hands = append(hands, &Hand{
			Bid:   aoc.MustParse(parts[1]),
			Cards: strings.Split(parts[0], ""),
		})
	})

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Less(hands[j])
	})

	var total int64
	for i, h := range hands {
		rank := int64(i + 1)
		total += rank * h.Bid
	}

	log.Printf("total: %v", total)
}
