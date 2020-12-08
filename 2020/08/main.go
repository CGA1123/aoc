package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var InstructionSet = map[string]func(*Program, int){
	"nop": func(p *Program, i int) {
		p.Counter = p.Counter + 1
	},
	"jmp": func(p *Program, i int) {
		p.Counter = p.Counter + i
	},
	"acc": func(p *Program, i int) {
		p.Acc = p.Acc + i
		p.Counter = p.Counter + 1
	},
}

type Instruction struct {
	Op  string
	Arg int
}

type Program struct {
	Acc          int
	Counter      int
	Instructions []Instruction
}

func (p *Program) Step() bool {
	if p.Counter == len(p.Instructions) {
		return false
	}

	i := p.Instructions[p.Counter]

	InstructionSet[i.Op](p, i.Arg)

	return true
}

func Run(p *Program) (bool, []int) {
	alreadyRan := map[int]struct{}{}
	var exec []int
	for {
		_, ok := alreadyRan[p.Counter]
		if ok {
			return false, exec
		}

		exec = append(exec, p.Counter)
		alreadyRan[p.Counter] = struct{}{}

		cont := p.Step()
		if !cont {
			return true, exec
		}
	}
}

func Flip(i Instruction) Instruction {
	var newOp string
	if i.Op == "jmp" {
		newOp = "nop"
	} else {
		newOp = "jmp"
	}

	return Instruction{Op: newOp, Arg: i.Arg}
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	var instructions []Instruction
	for s.Scan() {
		line := strings.Split(s.Text(), " ")
		argument, _ := strconv.Atoi(line[1])
		instructions = append(instructions, Instruction{Op: line[0], Arg: argument})
	}
	if s.Err() != nil {
		log.Printf("error: %v", s.Err())
	}

	// Part 1
	p := &Program{Instructions: instructions}
	_, ran := Run(p)
	log.Printf("pt(1): %v (at %v)", p.Acc, p.Counter)

	// Part 2

	// collect all jmp/nop operations in the initial cycle
	var contenders []int
	for _, i := range ran {
		if p.Instructions[i].Op != "acc" {
			contenders = append(contenders, i)
		}
	}

	// for every jmp/nop flip it and see if it terminates
	for _, i := range contenders {
		instructions[i] = Flip(instructions[i])

		attempt := &Program{Instructions: instructions}
		if terminated, _ := Run(attempt); terminated {
			log.Printf("pt(2): %v", attempt.Acc)
			return
		}

		// unclobber
		instructions[i] = Flip(instructions[i])
	}
}
