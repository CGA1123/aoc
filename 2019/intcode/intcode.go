package intcode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type In interface {
	Read() int64
}

type Out interface {
	Write(int64)
}

// NullIO
type NullIO struct{}

func NewNullIO() *NullIO {
	return &NullIO{}
}

func (r *NullIO) Read() int64 {
	panic("NullIO was read.")
}

func (r *NullIO) Write(i int64) {
	panic("NullIO was writter.")
}

// ChanIO
type ChanIO struct {
	ch chan int64
}

func NewChanIO(buffer int) *ChanIO {
	return &ChanIO{ch: make(chan int64, buffer)}
}

func (c *ChanIO) Chan() chan int64 {
	return c.ch
}

func (c *ChanIO) Read() int64 {
	return <-c.ch
}

func (c *ChanIO) Write(i int64) {
	c.ch <- i
}

// Intcode
type Intcode struct {
	pc      int64
	id      int64
	mem     []int64
	err     error
	relBase int64
	in      In
	out     Out
}

func FileMem(s string) ([]int64, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewScanner(bufio.NewReader(f))
	buf.Scan()

	var mem []int64
	for _, n := range strings.Split(buf.Text(), ",") {
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return nil, err
		}

		mem = append(mem, i)
	}

	return mem, nil
}

func New(mem []int64, in In, out Out, id int64) *Intcode {
	safeMem := make([]int64, len(mem)*100)
	copy(safeMem, mem)

	return &Intcode{
		mem: safeMem,
		id:  id,
		in:  in,
		out: out}
}

func (ic *Intcode) Err() error {
	return ic.err
}

func (ic *Intcode) Do() bool {
	code := ic.mem[ic.pc]

	modes := parameterModes(code)

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
		ic.outp(modes)
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
	in := ic.in.Read()

	ic.write(0, modes, in)
	ic.pc += 2
}

func (ic *Intcode) outp(modes []int64) {
	out := ic.val(0, modes)
	ic.out.Write(out)
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
