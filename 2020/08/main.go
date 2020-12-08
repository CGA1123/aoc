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

func Cycle(p *Program) (bool, int) {

	return false, 0
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

	newInstructions := make([]Instruction, len(instructions))

	// for every jmp/nop flip it and see if it terminates
	for _, i := range contenders {
		copy(newInstructions, instructions)

		ins := newInstructions[i]
		var newOp string
		if ins.Op == "jmp" {
			newOp = "nop"
		} else {
			newOp = "jmp"
		}

		newInstructions[i] = Instruction{Op: newOp, Arg: ins.Arg}

		attempt := &Program{Instructions: newInstructions}

		terminated, _ := Run(attempt)
		if terminated {
			log.Printf("pt(2): %v", attempt.Acc)
		}
	}
}
