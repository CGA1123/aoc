package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func additional(i int) int {
	amount := (i / 3) - 2
	if amount <= 0 {
		return 0
	}

	return amount + additional(amount)
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var e error
	var fuel int
	var additionalFuel int

	for {
		l, e := r.ReadString('\n')
		if e != nil {
			break
		}

		i, e := strconv.Atoi(strings.TrimSuffix(l, "\n"))
		if e != nil {
			break
		}

		f := (i / 3) - 2
		fuel = fuel + f
		additionalFuel = additionalFuel + additional(f)
	}
	if e != nil {
		log.Printf("error: %v", e)
		return
	}

	log.Printf("fuel: %v", fuel)
	log.Printf("fuel (pt2): %v", fuel+additionalFuel)
}
