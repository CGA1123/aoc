package main

import (
	"log"
	"strconv"
)

func MatchOne(s string) bool {
	repeat := false

	for i := 1; i < len(s); i++ {
		if s[i] < s[i-1] {
			return false
		}

		repeat = repeat || (s[i] == s[i-1])
	}

	return repeat
}

func MatchTwo(s string) bool {
	repeat := false

	for i := 1; i < len(s); i++ {
		if s[i] < s[i-1] {
			return false
		}

		sequence := s[i] == s[i-1]
		notPrevious := i < 2 || s[i-2] != s[i]
		notNext := i == 5 || s[i+1] != s[i]

		repeat = repeat || (sequence && notPrevious && notNext)
	}

	return repeat
}

func main() {
	var countOne int
	var countTwo int
	for i := 347312; i <= 805915; i++ {
		if MatchOne(strconv.Itoa(i)) {
			countOne++
		}

		if MatchTwo(strconv.Itoa(i)) {
			countTwo++
		}
	}

	log.Printf("pt(1): %v", countOne)
	log.Printf("pt(2): %v", countTwo)
}
