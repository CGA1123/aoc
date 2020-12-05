package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
)

type Row struct {
	back  int
	left  int
	front int
	right int
}

func NewRow() *Row {
	return &Row{front: 0, back: 127, left: 0, right: 7}
}

func F(r *Row) *Row {
	delta := (r.back - r.front - 1) / 2

	r.back = r.front + delta

	return r
}

func B(r *Row) *Row {
	delta := (r.back - r.front - 1) / 2

	r.front = r.back - delta

	return r
}

func L(r *Row) *Row {
	delta := (r.right - r.left - 1) / 2

	r.right = r.left + delta

	return r
}

func R(r *Row) *Row {
	delta := (r.right - r.left - 1) / 2

	r.left = r.right - delta

	return r
}

var f = map[byte](func(*Row) *Row){
	'F': F,
	'B': B,
	'L': L,
	'R': R}

func solve(l []byte) int {
	r := NewRow()

	for _, b := range l {
		r = f[b](r)
	}

	return (r.front * 8) + r.right
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)

	var l []byte
	var e error
	var max int
	var all []int

	for {
		l, e = r.ReadBytes('\n')
		if e != nil {
			break
		}
		l = l[:len(l)-1]

		i := solve(l)
		all = append(all, i)

		if i > max {
			max = i
		}

	}
	if e != nil && e != io.EOF {
		log.Printf("error: %v", e)
		return
	}

	log.Printf("max (pt1): %v", max)
	sort.Ints(all)

	for i := 1; i < len(all); i++ {
		if all[i]-all[i-1] != 1 {
			log.Printf("id (pt2): %v", all[i]-1)
			break
		}
	}
}
