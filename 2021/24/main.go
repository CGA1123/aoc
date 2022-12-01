package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/CGA1123/aoc"
)

type Variable int

const (
	W Variable = iota
	X
	Y
	Z
)

type Reference struct {
	Value Variable
}

func (r *Reference) Val(a *ALU) int64 {
	return a.variables[r.Value]
}

type Literal struct {
	Value int64
}

func (l *Literal) Val(a *ALU) int64 {
	return l.Value
}

type Operand interface {
	Val(*ALU) int64
}

type ALU struct {
	variables map[Variable]int64
	input     []int64
	inputIdx  int
}

func (a *ALU) Inp(v Variable) {
	a.variables[v] = a.input[a.inputIdx]
	a.inputIdx++
}

func (a *ALU) Add(x Variable, y Operand) {
	a.variables[x] = a.variables[x] + y.Val(a)
}

func (a *ALU) Mul(x Variable, y Operand) {
	a.variables[x] = a.variables[x] * y.Val(a)
}

func (a *ALU) Div(x Variable, y Operand) {
	a.variables[x] = a.variables[x] / y.Val(a)
}

func (a *ALU) Mod(x Variable, y Operand) {
	a.variables[x] = a.variables[x] % y.Val(a)
}

func (a *ALU) Eql(x Variable, y Operand) {
	a.variables[x] = a.variables[x] % y.Val(a)
}

func (a *ALU) Read(x Variable) int64 {
	return a.variables[x]
}

func (a *ALU) Reset(in []int64) {
	a.input = in
	a.inputIdx = 0
	a.variables[W] = 0
	a.variables[X] = 0
	a.variables[Y] = 0
	a.variables[Z] = 0
}

func NewALU() *ALU {
	return &ALU{
		variables: map[Variable]int64{
			W: 0,
			X: 0,
			Y: 0,
			Z: 0,
		},
	}
}

func main() {
	generator(func(i []int64) {
		fmt.Printf("%v\n", asNumber(i))
	})
	program := []func(*ALU){}
	inp := func(v Variable) func(a *ALU) {
		return func(a *ALU) {
			a.Inp(v)
		}
	}
	ops := map[string]func(Variable, Operand) func(*ALU){
		"div": func(v Variable, o Operand) func(*ALU) {
			return func(a *ALU) {
				a.Div(v, o)
			}
		},
		"add": func(v Variable, o Operand) func(*ALU) {
			return func(a *ALU) {
				a.Add(v, o)
			}
		},
		"mod": func(v Variable, o Operand) func(*ALU) {
			return func(a *ALU) {
				a.Mod(v, o)
			}
		},
		"mul": func(v Variable, o Operand) func(*ALU) {
			return func(a *ALU) {
				a.Mul(v, o)
			}
		},
		"eql": func(v Variable, o Operand) func(*ALU) {
			return func(a *ALU) {
				a.Eql(v, o)
			}
		},
	}

	aoc.EachLine("input.txt", func(s string) {
		cmd := strings.Split(s, " ")
		operation := cmd[0]
		v, err := parseVariable(cmd[1])
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		if operation == "inp" {
			program = append(program, inp(v))
			return
		}

		fn, ok := ops[operation]
		if !ok {
			log.Fatalf("error: unknown operation %v", operation)
		}

		operand, err := parseOperand(cmd[2])
		if err != nil {
			log.Fatal("error: %v", err)
		}

		program = append(program, fn(v, operand))
	})

	results := make(chan int64)
	inputs := make(chan []int64)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			a := NewALU()

			for input := range inputs {
				a.Reset(input)
				for _, i := range program {
					i(a)
				}
				if a.Read(Z) == 0 {
					results <- asNumber(input)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		// generate inputs

		close(inputs)
	}()

	// read results
	max := int64(0)
	go func() {
		for result := range results {
			if result > max {
				max = result
			}
		}
	}()

	wg.Wait()
	close(results)

	fmt.Printf("%v\n", max)
}

func asNumber(input []int64) int64 {
	coeff := int64(1)
	res := int64(0)
	for i := len(input) - 1; i >= 0; i-- {
		res = res + (input[i] * coeff)
		coeff = coeff * 10
	}

	return res
}

func parseOperand(o string) (Operand, error) {
	v, varErr := parseVariable(o)
	if varErr == nil {
		return &Reference{Value: v}, nil
	}

	return &Literal{Value: aoc.MustParse(o)}, nil
}

func parseVariable(v string) (Variable, error) {
	switch v {
	case "w":
		return W, nil
	case "x":
		return X, nil
	case "y":
		return Y, nil
	case "z":
		return Z, nil
	}

	return W, fmt.Errorf("bad variable: %v", v)
}

var valid = []int64{9, 8, 7, 6, 5, 4, 3, 2, 1}

func generate(in []int64, i int, out func([]int64)) {
	if i == 14 {
		out(in)
		return
	}

	for _, val := range valid {
		in[i] = val

		generate(in, i+1, out)
	}
}

func generator(out func([]int64)) {
	v := make([]int64, 14)

	generate(v, 0, out)
}
