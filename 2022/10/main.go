package main

import (
	"fmt"
	"strings"

	"github.com/CGA1123/aoc"
)

type cpu struct {
	toAdd        *int64
	x            int64
	instructions []string
	ip           int
	t            int
}

func newCpu(i []string) *cpu {
	return &cpu{
		x: 1, instructions: i,
	}
}

func (c *cpu) cycle() {
	if c.toAdd != nil {
		c.x = c.x + *c.toAdd
		c.toAdd = nil
		c.ip++
		return
	}

	i := c.instructions[c.ip]
	if i == "noop" {
		c.ip++
		return
	}

	addx := strings.Split(i, " ")
	n := aoc.MustParse(addx[1])
	c.toAdd = &n
}

func main() {
	instructions := []string{}

	aoc.EachLine("input.txt", func(s string) {
		instructions = append(instructions, s)
	})

	c := newCpu(instructions)
	isLit := func(pixel int64) bool {
		pixel = pixel % 40

		switch pixel - c.x {
		case -1, 0, 1:
			return true
		}

		return false
	}

	var sum int64
	for pixel := int64(0); pixel < 240; pixel++ {
		cycle := pixel + 1
		if cycle == 20 || ((cycle-20)%40) == 0 {
			sum = sum + (c.x * pixel)
		}

		if isLit(pixel) {
			fmt.Printf("#")
		} else {
			fmt.Printf(" ")
		}

		c.cycle()

		if (pixel+1)%40 == 0 {
			fmt.Print("\n")
		}
	}
	fmt.Printf("%v\n", sum)
}
