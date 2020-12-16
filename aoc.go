package aoc

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Set struct {
	set map[interface{}]bool
}

func NewSet() *Set {
	return &Set{set: map[interface{}]bool{}}
}

func (s *Set) Add(e interface{}) {
	s.set[e] = true
}

func (s *Set) Remove(e interface{}) {
	delete(s.set, e)
}

func (s *Set) Contains(e interface{}) bool {
	_, ok := s.set[e]

	return ok
}

func (s *Set) Size() int {
	return len(s.set)
}

func (s *Set) Elements() []interface{} {
	var el []interface{}

	for k := range s.set {
		el = append(el, k)
	}

	return el
}

func MustParse(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("bad input")
	}

	return i
}

func Capture(r *regexp.Regexp, s string) map[string]string {
	m := map[string]string{}
	names := r.SubexpNames()
	for _, match := range r.FindAllStringSubmatch(s, -1) {
		for i, submatch := range match {
			name := names[i]
			if name == "" {
				continue
			}

			m[name] = submatch
		}
	}

	return m
}

func Between(min, max int64) func(int64) bool {
	return func(i int64) bool {
		return min <= i && i <= max
	}
}

func Or(f, g func(int64) bool) func(int64) bool {
	return func(i int64) bool {
		return f(i) || g(i)
	}
}

func EachLine(input string, fn func(string)) error {
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("opening file (%v): %v", input, err)
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	for s.Scan() {
		fn(s.Text())
	}

	return s.Err()
}
