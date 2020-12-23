package main

import (
	"log"
	"strconv"

	"github.com/CGA1123/aoc"
)

const Million = 1000000

var Input = []int{9, 5, 2, 4, 3, 8, 7, 1, 6}
var Example = []int{3, 8, 9, 1, 2, 5, 4, 6, 7}

type CrabCups struct {
	cups    []int
	current int
	length  int
	grabbed [3]int
}

func NewCrabCups(initialCups []int) *CrabCups {
	total := len(initialCups)
	cups := make([]int, len(initialCups)+1)

	for i, v := range initialCups {
		cups[v] = initialCups[(i+1)%total]
	}

	return &CrabCups{cups: cups, current: initialCups[0], length: total, grabbed: [3]int{}}
}

func (cc *CrabCups) decr(i int) int {
	n := i - 1

	if n == 0 {
		return cc.length
	} else {
		return n
	}
}

func (cc *CrabCups) contain(n int) bool {
	for _, v := range cc.grabbed {
		if v == n {
			return true
		}
	}

	return false
}

func (cc *CrabCups) next() int {
	d := cc.decr(cc.current)

	for i := 0; i < 3; i++ {
		if !cc.contain(d) {
			break
		}

		d = cc.decr(d)
	}

	return d
}

func (cc *CrabCups) Play() {
	cc.grabbed[0] = cc.cups[cc.current]
	cc.grabbed[1] = cc.cups[cc.grabbed[0]]
	cc.grabbed[2] = cc.cups[cc.grabbed[1]]

	next := cc.next()

	first, last := cc.grabbed[0], cc.grabbed[2]

	cc.cups[cc.current] = cc.cups[last]
	cc.cups[last] = cc.cups[next]
	cc.cups[next] = first

	cc.current = cc.cups[cc.current]
}

func PartOne() string {
	cc := NewCrabCups(Input)

	for i := 0; i < 100; i++ {
		cc.Play()
	}

	var str string

	current := 1
	for i := 0; i < len(cc.cups)-1; i++ {
		str += strconv.Itoa(cc.cups[current])
		current = cc.cups[current]
	}

	return str
}

func PartTwo() int {
	list := make([]int, Million)
	input := Input

	for i := 0; i < Million; i++ {
		if i < len(input) {
			list[i] = input[i]
		} else {
			list[i] = i + 1
		}
	}

	ncc := NewCrabCups(list)

	for i := 0; i < 10*Million; i++ {
		ncc.Play()
	}

	return ncc.cups[1] * ncc.cups[ncc.cups[1]]
}

func main() {
	cleanup, err := aoc.Profile()
	if err != nil {
		log.Printf("couldn't profile: %v", err)
		return
	}
	defer cleanup()

	log.Printf("pt(1): %v", PartOne())
	log.Printf("pt(2): %v", PartTwo())
}
