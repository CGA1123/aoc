package main

import (
	"log"
	"regexp"

	"github.com/CGA1123/aoc"
)

var PlayerRegex = regexp.MustCompile(`Player (1|2):`)

type Deck struct {
	player int64
	cards  []int64
}

func NewDeck(player int64, size int) *Deck {
	return &Deck{player: player, cards: make([]int64, size, 100)}
}

func CopyDeck(d *Deck) *Deck {
	slice := make([]int64, d.Size(), 100)
	copy(slice, d.cards)

	return &Deck{player: d.Player(), cards: slice}
}

func (d *Deck) Player() int64 {
	return d.player
}

func (d *Deck) Draw() int64 {
	c := d.cards[0]
	d.cards = d.cards[1:]

	return c
}

func (d *Deck) Add(card int64) {
	d.cards = append(d.cards, card)
}

func (d *Deck) Size() int64 {
	return int64(len(d.cards))
}

func (d *Deck) Cards() []int64 {
	return d.cards
}

func (d *Deck) Score() int64 {
	var result int64

	cards := d.Cards()
	length := len(cards)

	for i, c := range cards {
		result += int64(length-i) * c
	}

	return result
}

type CrabCombat struct {
	a, b *Deck
}

func NewCrabCombat(a, b *Deck) *CrabCombat {
	return &CrabCombat{a: a, b: b}
}

func (cc *CrabCombat) Play() bool {
	if cc.Winner() != nil {
		return false
	}
	ac, bc := cc.a.Draw(), cc.b.Draw()

	if ac > bc {
		cc.a.Add(ac)
		cc.a.Add(bc)
	} else {
		cc.b.Add(bc)
		cc.b.Add(ac)
	}

	return cc.a.Size() != 0 && cc.b.Size() != 0
}

func (cc *CrabCombat) Winner() *Deck {
	as, bs := cc.a.Size(), cc.b.Size()
	if as != 0 && bs != 0 {
		return nil
	}

	if as == 0 {
		return cc.b
	} else {
		return cc.a
	}
}

type RecursiveCombat struct {
	a, b     *Deck
	winner   *Deck
	previous *aoc.Set
}

func NewRecursiveCombat(a, b *Deck) *RecursiveCombat {
	return &RecursiveCombat{a: a, b: b, previous: aoc.NewSetWithSize(750)}
}

func (rc *RecursiveCombat) Winner() *Deck {
	return rc.winner
}

type state struct {
	a, b int64
}

func (rc *RecursiveCombat) state() state {
	var astate int64
	for _, card := range rc.a.Cards() {
		astate = (astate * 100) + card
	}

	var bstate int64
	for _, card := range rc.b.Cards() {
		bstate = (bstate * 100) + card
	}

	return state{a: astate, b: bstate}
}

func (rc *RecursiveCombat) subGame(sa, sb int64) *Deck {
	ac, bc := rc.a.Cards()[:sa], rc.b.Cards()[:sb]

	ad, bd := NewDeck(rc.a.Player(), len(ac)), NewDeck(rc.b.Player(), len(bc))
	copy(ad.cards, ac)
	copy(bd.cards, bc)

	sg := NewRecursiveCombat(ad, bd)
	for sg.Play() {
	}

	winner := sg.Winner()

	if winner.Player() == rc.a.Player() {
		return rc.a
	}

	return rc.b
}

func (rc *RecursiveCombat) Play() bool {
	if rc.winner != nil {
		return false
	}

	state := rc.state()
	if rc.previous.Contains(state) {
		rc.winner = rc.a

		return false
	}
	rc.previous.Add(state)

	var roundWin *Deck
	ac, bc := rc.a.Draw(), rc.b.Draw()

	if rc.a.Size() >= ac && rc.b.Size() >= bc {
		roundWin = rc.subGame(ac, bc)
	} else {
		if ac > bc {
			roundWin = rc.a
		} else {
			roundWin = rc.b
		}
	}

	if roundWin.Player() == rc.a.Player() {
		roundWin.Add(ac)
		roundWin.Add(bc)
	} else {
		roundWin.Add(bc)
		roundWin.Add(ac)
	}

	if rc.a.Size() == 0 {
		rc.winner = rc.b
		return false
	}

	if rc.b.Size() == 0 {
		rc.winner = rc.a
		return false
	}

	return true
}

func PartOne(a, b *Deck) int64 {
	cc := NewCrabCombat(CopyDeck(a), CopyDeck(b))

	for cc.Play() {
	}

	return cc.Winner().Score()
}

func PartTwo(a, b *Deck) int64 {
	rc := NewRecursiveCombat(CopyDeck(a), CopyDeck(b))

	for rc.Play() {
	}

	return rc.Winner().Score()
}

func main() {
	a, b := NewDeck(1, 0), NewDeck(2, 0)
	current := a

	aoc.EachLine("input.txt", func(l string) {
		if l == "" {
			current = b
			return
		}

		if PlayerRegex.MatchString(l) {
			return
		}

		current.Add(aoc.MustParse(l))
	})

	log.Printf("pt(1): %v", PartOne(a, b))
	log.Printf("pt(2): %v", PartTwo(a, b))
}
