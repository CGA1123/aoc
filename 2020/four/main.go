package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var partOne = map[string](func(string) bool){
	"byr": Present,
	"iyr": Present,
	"eyr": Present,
	"hgt": Present,
	"hcl": Present,
	"ecl": Present,
	"pid": Present}

var partTwo = map[string](func(string) bool){
	"byr": IntBetween(1920, 2002),
	"iyr": IntBetween(2010, 2020),
	"eyr": IntBetween(2020, 2030),
	"hgt": ValidHeight,
	"hcl": Regex(`^#[0-9a-f]{6}$`),
	"ecl": Regex(`^(amb|blu|brn|gry|grn|hzl|oth)$`),
	"pid": Regex(`^[0-9]{9}$`)}

func parseIdentity(r *bufio.Reader) (map[string]string, error) {
	i := map[string]string{}
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			return i, err
		}

		l = strings.TrimSuffix(l, "\n")

		if l == "" {
			return i, nil
		}

		for _, kv := range strings.Split(l, " ") {
			field := strings.Split(kv, ":")
			if len(field) != 2 {
				return nil, fmt.Errorf("bad field: %v", kv)
			}

			i[field[0]] = field[1]
		}
	}
}

func IntBetween(min, max int) func(string) bool {
	return func(s string) bool {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false
		}

		return min <= i && i <= max
	}
}

func HasTrimSuffix(s, suf string) (string, bool) {
	if strings.HasSuffix(s, suf) {
		return strings.TrimSuffix(s, suf), true
	}

	return "", false
}

func ValidHeight(s string) bool {
	cm := IntBetween(150, 193)
	in := IntBetween(59, 76)

	if cmVal, ok := HasTrimSuffix(s, "cm"); ok {
		return cm(cmVal)
	}

	if inVal, ok := HasTrimSuffix(s, "in"); ok {
		return in(inVal)
	}

	return false
}

func Regex(exp string) func(string) bool {
	r := regexp.MustCompile(exp)

	return r.MatchString
}

func Present(s string) bool {
	return s != ""
}

func valid(i map[string]string, validators map[string]func(string) bool) bool {
	for k, f := range validators {
		v, ok := i[k]
		if !ok || !f(v) {
			return false
		}
	}

	return true
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
	var countOne int
	var countTwo int

	// TODO: error handling is a bit messy with EOF checks, Probably need to
	// separate parsing and reading a little :shrug:
	for {
		i, e := parseIdentity(r)
		if e != nil && e != io.EOF {
			break
		}

		if valid(i, partOne) {
			countOne = countOne + 1
		}

		if valid(i, partTwo) {
			countTwo = countTwo + 1
		}

		if e != nil {
			break
		}
	}
	if e != nil && e != io.EOF {
		return
	}

	log.Printf("count (pt 1): %v", countOne)
	log.Printf("count (pt 2): %v", countTwo)
}
