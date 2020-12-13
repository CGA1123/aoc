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

type Intcode struct {
	pc      int64
	id      int64
	mem     []int64
	err     error
	relBase int64
	inChan  chan int64
	outChan chan int64
}

func NewIntcode(mem []int64, in chan int64, out chan int64, id int64) *Intcode {
	safeMem := make([]int64, len(mem)*100)
	copy(safeMem, mem)

	return &Intcode{
		mem:     safeMem,
		id:      id,
		inChan:  in,
		outChan: out}
}

func (ic *Intcode) In() chan<- int64 {
	return ic.inChan
}

func (ic *Intcode) Out() <-chan int64 {
	return ic.outChan
}

func (ic *Intcode) Err() error {
	return ic.err
}

func (ic *Intcode) Do() bool {
	code := ic.mem[ic.pc]

	modes := parameterModes(code)

	// log.Print64f("%v", code)

	switch op := opcode(code); op {
	case 9:
		ic.rel(modes)
		return true
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

func opcode(i int64) int64 {
	return i % 100
}

func parameterModes(i int64) []int64 {
	var modes []int64
	current := i / 100

	for j := 0; j < 3; j++ {
		mode := current % 10
		current = current / 10

		modes = append(modes, mode)
	}

	return modes
}

func (ic *Intcode) val(i int64, modes []int64) int64 {
	var value int64

	idx := ic.pc + i + 1

	switch modes[i] {
	case 0: // position
		value = ic.mem[ic.mem[idx]]
	case 1: // immediate
		value = ic.mem[idx]
	case 2: // relative
		value = ic.mem[ic.relBase+ic.mem[idx]]
	}

	return value
}

func (ic *Intcode) write(i int64, modes []int64, value int64) {
	idx := ic.pc + i + 1
	var dest int64

	switch modes[i] {
	case 0: // position
		dest = ic.mem[idx]
	case 2: // relative
		dest = ic.relBase + ic.mem[idx]
	}

	ic.mem[dest] = value
}

func (ic *Intcode) rel(modes []int64) {
	value := ic.val(0, modes)

	ic.relBase += value
	ic.pc += 2
}

func (ic *Intcode) eq(modes []int64) {
	a, b := ic.val(0, modes), ic.val(1, modes)
	var result int64
	if a == b {
		result = 1
	}

	ic.write(2, modes, result)
	ic.pc += 4
}

func (ic *Intcode) lt(modes []int64) {
	a, b := ic.val(0, modes), ic.val(1, modes)

	var result int64
	if a < b {
		result = 1
	}

	ic.write(2, modes, result)
	ic.pc += 4
}

func (ic *Intcode) jmpf(modes []int64) {
	v, to := ic.val(0, modes), ic.val(1, modes)

	if v == 0 {
		ic.pc = to
		return
	}

	ic.pc += 3
}

func (ic *Intcode) jmpt(modes []int64) {
	v, to := ic.val(0, modes), ic.val(1, modes)

	if v != 0 {
		ic.pc = to
		return
	}

	ic.pc += 3
}

func (ic *Intcode) inp(modes []int64) {
	in := <-ic.inChan

	ic.write(0, modes, in)
	ic.pc += 2
}

func (ic *Intcode) out(modes []int64) {
	out := ic.val(0, modes)
	ic.outChan <- out
	ic.pc += 2
}

func (ic *Intcode) add(modes []int64) {
	a, b := ic.val(0, modes), ic.val(1, modes)

	ic.write(2, modes, a+b)
	ic.pc += 4
}

func (ic *Intcode) mul(modes []int64) {
	a, b := ic.val(0, modes), ic.val(1, modes)

	ic.write(2, modes, a*b)
	ic.pc += 4
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

	var mem []int64
	for _, n := range strings.Split(s.Text(), ",") {
		i, _ := strconv.ParseInt(n, 10, 64)
		mem = append(mem, i)
	}

	in := make(chan int64, 2)
	out := make(chan int64, 2)
	ic := NewIntcode(mem, in, out, 0)

	ic.In() <- 2

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for x := range ic.Out() {
			fmt.Printf("-> %v\n", x)
		}
		wg.Done()
	}()

	for ic.Do() {
	}
	close(in)
	close(out)

	wg.Wait()
}
