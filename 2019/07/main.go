package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Perm calls f with each permutation of a.
func Perm(a []int, f func([]int)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, f func([]int), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

type Intcode struct {
	pc      int
	id      int
	mem     []int
	err     error
	inChan  chan int
	outChan chan int
}

func NewIntcode(mem []int, in chan int, out chan int, id int) *Intcode {
	safeMem := make([]int, len(mem))
	copy(safeMem, mem)

	return &Intcode{
		mem:     safeMem,
		id:      id,
		inChan:  in,
		outChan: out}
}

func (ic *Intcode) In() chan<- int {
	return ic.inChan
}

func (ic *Intcode) Out() <-chan int {
	return ic.outChan
}

func (ic *Intcode) Err() error {
	return ic.err
}

func (ic *Intcode) Do() bool {
	code := ic.mem[ic.pc]
	modes := parameterModes(code)

	switch op := opcode(code); op {
	case 8:
		ic.eq(modes)
		return true
	case 7:
		ic.lt(modes)
		return true
	case 6:
		ic.jmpf(modes)
		return true
	case 5:
		ic.jmpt(modes)
		return true
	case 4:
		ic.out(modes)
		return true
	case 3:
		ic.inp(modes)
		return true
	case 2:
		ic.mul(modes)
		return true
	case 1:
		ic.add(modes)
		return true
	case 99:
		return false
	default:
		ic.err = fmt.Errorf("bad opcode: %v (op %v) (pc %v)", code, op, ic.pc)
		return false
	}
}

func opcode(i int) int {
	return i % 100
}

func parameterModes(i int) []bool {
	var modes []bool
	current := i / 100

	for j := 0; j < 3; j++ {
		mode := current % 10
		current = current / 10

		modes = append(modes, mode == 1)
	}

	return modes
}

func (ic *Intcode) val(i int, immediate bool) int {
	value := ic.mem[i]
	if !immediate {
		value = ic.mem[value]
	}

	return value
}

func (ic *Intcode) eq(modes []bool) {
	a, b := ic.val(ic.pc+1, modes[0]), ic.val(ic.pc+2, modes[1])
	var result int
	if a == b {
		result = 1
	}

	dest := ic.mem[ic.pc+3]
	ic.mem[dest] = result
	ic.pc = ic.pc + 4
}

func (ic *Intcode) lt(modes []bool) {
	a, b := ic.val(ic.pc+1, modes[0]), ic.val(ic.pc+2, modes[1])

	var result int
	if a < b {
		result = 1
	}

	dest := ic.mem[ic.pc+3]
	ic.mem[dest] = result
	ic.pc = ic.pc + 4
}

func (ic *Intcode) jmpf(modes []bool) {
	v := ic.val(ic.pc+1, modes[0])
	if v == 0 {
		ic.pc = ic.val(ic.pc+2, modes[1])
		return
	}

	ic.pc = ic.pc + 3
}

func (ic *Intcode) jmpt(modes []bool) {
	v := ic.val(ic.pc+1, modes[0])
	if v != 0 {
		ic.pc = ic.val(ic.pc+2, modes[1])
		return
	}

	ic.pc = ic.pc + 3
}

func (ic *Intcode) inp(modes []bool) {
	in := <-ic.inChan
	dest := ic.mem[ic.pc+1]
	ic.mem[dest] = in
	ic.pc = ic.pc + 2
}

func (ic *Intcode) out(modes []bool) {
	out := ic.val(ic.pc+1, modes[0])

	ic.outChan <- out
	ic.pc = ic.pc + 2
}

func (ic *Intcode) add(modes []bool) {
	a, b := ic.val(ic.pc+1, modes[0]), ic.val(ic.pc+2, modes[1])

	dest := ic.mem[ic.pc+3]
	ic.mem[dest] = a + b
	ic.pc = ic.pc + 4
}

func (ic *Intcode) mul(modes []bool) {
	a, b := ic.val(ic.pc+1, modes[0]), ic.val(ic.pc+2, modes[1])

	dest := ic.mem[ic.pc+3]
	ic.mem[dest] = a * b
	ic.pc = ic.pc + 4
}

func Attempt(input, mem []int, chans []chan int) int {
	var amplifiers []*Intcode
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		amp := NewIntcode(mem, chans[i], chans[i+1], i)
		amp.In() <- input[i]

		go func() {
			for amp.Do() {
			}

			wg.Done()
		}()

		amplifiers = append(amplifiers, amp)
	}

	amplifiers[0].In() <- 0

	wg.Wait()

	return <-amplifiers[len(amplifiers)-1].Out()
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	s.Scan()

	var mem []int
	for _, n := range strings.Split(s.Text(), ",") {
		i, _ := strconv.Atoi(n)
		mem = append(mem, i)
	}

	var n int
	Perm([]int{0, 1, 2, 3, 4}, func(config []int) {
		var chans []chan int
		for i := 0; i < 6; i++ {
			chans = append(chans, make(chan int, 2))
		}

		a := Attempt(config, mem, chans)
		if a > n {
			n = a
		}
	})

	log.Printf("pt(1): %v", n)

	var m int

	Perm([]int{5, 6, 7, 8, 9}, func(config []int) {
		var chans []chan int
		for i := 0; i < 5; i++ {
			chans = append(chans, make(chan int, 2))
		}
		chans = append(chans, chans[0])

		a := Attempt(config, mem, chans)
		if a > m {
			m = a
		}
	})

	log.Printf("pt(2): %v", m)
}
