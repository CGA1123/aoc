package main

import (
	"log"
	"strings"
	"unicode"

	"github.com/CGA1123/aoc"
)

func calibrationValue(s string) int64 {
	var val int64

	for _, c := range s {
		if unicode.IsDigit(c) {
			val += aoc.MustParse(string(c)) * 10

			break
		}
	}
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		c := runes[i]
		if unicode.IsDigit(c) {
			val += aoc.MustParse(string(c))

			break
		}
	}

	return val
}

var numbers = map[string]int64{
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func fixedCalibrationValue(s string) int64 {
	var value int64

First:
	for i := 0; i < len(s); i++ {
		sub := s[i:]

		for str, val := range numbers {
			if strings.HasPrefix(sub, str) {
				value += (val * 10)
				break First
			}
		}
	}
Last:
	for i := len(s) - 1; i >= 0; i-- {
		sub := s[i:]

		for str, val := range numbers {
			if strings.HasPrefix(sub, str) {
				value += val
				break Last
			}
		}
	}

	return value

}

func main() {
	var totalOne int64
	var totalTwo int64
	aoc.EachLine("input.txt", func(s string) {
		totalOne += calibrationValue(s)
		totalTwo += fixedCalibrationValue(s)
	})

	log.Printf("%v", totalOne)
	log.Printf("%v", totalTwo)
}
