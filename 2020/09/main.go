package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

type XMASEncryption struct {
	data     []int
	valid    map[int]map[int]struct{}
	preamble int
	idx      int
}

func NewXMASEncryption(data []int, preamble int) *XMASEncryption {
	x := &XMASEncryption{
		data:     data,
		preamble: preamble,
		idx:      preamble,
		valid:    map[int]map[int]struct{}{}}

	return x.init()
}

func (x *XMASEncryption) init() *XMASEncryption {
	for i := 0; i < x.idx; i++ {
		x.valid[i] = map[int]struct{}{}

		for j := i + 1; j < x.idx; j++ {
			xal := x.data[i] + x.data[j]

			x.valid[i][xal] = struct{}{}
		}
	}

	return x
}

func (x *XMASEncryption) contains(i int) bool {
	for _, x := range x.valid {
		if _, ok := x[i]; ok {
			return true
		}
	}

	return false
}

func (x *XMASEncryption) next() {
	delete(x.valid, x.idx-x.preamble)

	x.idx = x.idx + 1

	for k := range x.valid {
		x.valid[k][x.data[k]+x.data[x.idx]] = struct{}{}
	}

	x.valid[x.idx] = map[int]struct{}{}
}

func (x *XMASEncryption) Crack() (int, bool) {
	for i := x.preamble; i < len(x.data); i++ {
		val := x.data[i]

		if !x.contains(val) {
			return val, true
		}

		x.next()
	}

	return 0, false
}

// Attempt tries to find a contiguous set of integers in the sequence that sum
// to the targer, returning the min and max integer within the sequence
func Attempt(c []int, i, t int) (bool, int, int) {
	count := c[i]
	sequence := []int{c[i]}
	for j := i + 1; j < len(c); j++ {
		count = count + c[j]
		if count > t {
			break
		}

		sequence = append(sequence, c[j])

		if count == t {
			sort.Ints(sequence)
			return true, sequence[0], sequence[len(sequence)-1]
		}
	}

	return false, 0, 0
}

func main() {
	preamble := 25
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	var ciphertext []int
	for s.Scan() {
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		ciphertext = append(ciphertext, i)
	}
	if s.Err() != nil {
		log.Printf("error: %v", s.Err())
		return
	}

	toFind, ok := NewXMASEncryption(ciphertext, preamble).Crack()
	if !ok {
		log.Printf("couldn't crack encryption!")
		return
	}
	log.Printf("pt(1): %v", toFind)

	for i := 0; i < len(ciphertext)-1; i++ {
		ok, start, end := Attempt(ciphertext, i, toFind)
		if ok {
			log.Printf("pt(2): %v", start+end)
		}
	}
}
