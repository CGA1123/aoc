package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Intcode struct {
	pc  int
	mem []int
	err error
}

func NewIntcode(mem []int) *Intcode {
	return &Intcode{mem: mem}
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
	fmt.Printf("> ")
	var i int
	fmt.Scanf("%d", &i)

	dest := ic.mem[ic.pc+1]
	ic.mem[dest] = i
	ic.pc = ic.pc + 2
}

func (ic *Intcode) out(modes []bool) {
	log.Printf("%v", ic.val(ic.pc+1, modes[0]))
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

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer f.Close()

	s, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	var ins []int
	for _, number := range strings.Split(string(s[:len(s)-1]), ",") {
		i, err := strconv.Atoi(number)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		ins = append(ins, i)
	}

	ic := NewIntcode(ins)

	for ic.Do() {
	}
	if ic.Err() != nil {
		log.Printf("error: %v", ic.Err())
		return
	}
}
