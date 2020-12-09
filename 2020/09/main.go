package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

// Init initialises the set of valid integers based on the preamble
func Init(ciphertext []int, preamble int) []map[int]struct{} {
	var valid []map[int]struct{}
	for i := 0; i < preamble; i++ {
		valid = append(valid, map[int]struct{}{})

		for j := 0; j < preamble; j++ {
			if i == j {
				continue
			}

			val := ciphertext[i] + ciphertext[j]

			valid[i][val] = struct{}{}
		}
	}

	return valid
}

// Valid checks whether an integer is valid
func Valid(valid []map[int]struct{}, val int) bool {
	for _, v := range valid {
		if _, ok := v[val]; ok {
			return true
		}
	}

	return false
}

// Next builds the set of valid integers based on the current index
func Next(ciphertext []int, preamble, i int) map[int]struct{} {
	next := map[int]struct{}{}
	val := ciphertext[i]
	for j := 1; j < preamble; j++ {
		other := ciphertext[i-j]
		next[other+val] = struct{}{}
	}

	return next
}

// Attempt tries to find a contiguous set of integers in the sequence that sum
// to the targer, returning the min and max integer within the sequence
func Attempt(c []int, i, t int) (bool, int, int) {
	count := c[i]
	sequence := []int{c[i]}
	for j := i + 1; j < len(c); j++ {
		count = count + c[j]
		if count > t {
			break
		}

		sequence = append(sequence, c[j])

		if count == t {
			sort.Ints(sequence)
			return true, sequence[0], sequence[len(sequence)-1]
		}
	}

	return false, 0, 0
}

func main() {
	preamble := 25
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))

	var ciphertext []int
	for s.Scan() {
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		ciphertext = append(ciphertext, i)
	}
	if s.Err() != nil {
		log.Printf("error: %v", s.Err())
		return
	}

	valid := Init(ciphertext, preamble)
	var toFind int
	for i := preamble; i < len(ciphertext); i++ {
		val := ciphertext[i]

		if !Valid(valid, val) {
			log.Printf("pt(1): %v", val)
			toFind = val
			break
		}

		valid = append(valid, Next(ciphertext, preamble, i))
		valid = valid[1:]
	}

	for i := 0; i < len(ciphertext)-1; i++ {
		ok, start, end := Attempt(ciphertext, i, toFind)
		if ok {
			log.Printf("pt(2): %v", start+end)
		}
	}
}
