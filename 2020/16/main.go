package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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

var ruleregex = regexp.MustCompile(`(?P<field>.*): (?P<loa>\d+)-(?P<hia>\d+) or (?P<lob>\d+)-(?P<hib>\d+)`)

type Rule struct {
	Field     string
	Satisfied func(int64) bool
}

func Invalid(ticket []int64, rules []*Rule) []int64 {
	var invalid []int64
	for _, field := range ticket {
		fieldValid := false

		for _, rule := range rules {
			if rule.Satisfied(field) {
				fieldValid = true
				break
			}
		}

		if !fieldValid {
			invalid = append(invalid, field)
		}
	}

	return invalid
}

func PartOne(tickets [][]int64, rules []*Rule) int64 {
	var c int64
	for _, ticket := range tickets {
		for _, field := range Invalid(ticket, rules) {
			c += field
		}
	}

	return c
}

func ValidTickets(tickets [][]int64, rules []*Rule) [][]int64 {
	var v [][]int64

	for _, ticket := range tickets {
		if len(Invalid(ticket, rules)) != 0 {
			continue
		}

		v = append(v, ticket)
	}

	return v
}

func FindField(tickets [][]int64, rules []*Rule, i int) *Set {
	allTickets := func(t [][]int64, rule *Rule, i int) bool {
		for _, ticket := range t {
			if !rule.Satisfied(ticket[i]) {
				return false
			}
		}

		return true
	}

	set := NewSet()
	for _, rule := range rules {
		if allTickets(tickets, rule, i) {
			set.Add(rule.Field)
		}
	}

	return set
}

// TODO cleanup :thinking:

type solve struct {
	set *Set
	idx int
}

func Solve(fields []*Set) []string {
	found := NewSet()

	o := map[int]solve{}
	names := make([]string, len(fields))

	for i, field := range fields {
		o[field.Size()] = solve{set: field, idx: i}
	}

	for i := range fields {
		s := o[i+1]
		for _, element := range s.set.Elements() {
			if !found.Contains(element) {
				found.Add(element)

				names[s.idx] = element.(string)
				break
			}
		}

	}

	return names
}

func PartTwo(mine []int64, tickets [][]int64, rules []*Rule) int64 {
	validTickets := append(ValidTickets(tickets, rules), mine)

	var fields []*Set

	for i := 0; i < len(mine); i++ {
		fields = append(fields, FindField(validTickets, rules, i))
	}

	result := int64(1)
	for i, field := range Solve(fields) {
		if strings.HasPrefix(field, "departure") {
			result *= mine[i]
		}
	}

	return result
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Printf("opening file: %v", err)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(bufio.NewReader(f))
	var rules []*Rule

	scanRule := func(s string) {
		m := Capture(ruleregex, s)
		loa, hia, lob, hib := MustParse(m["loa"]), MustParse(m["hia"]), MustParse(m["lob"]), MustParse(m["hib"])

		rules = append(rules, &Rule{
			Field:     m["field"],
			Satisfied: Or(Between(loa, hia), Between(lob, hib)),
		})
	}

	scanRules := func(s *bufio.Scanner) {
		s.Scan()
		for s.Text() != "" {
			scanRule(s.Text())
			s.Scan()
		}
	}

	scanTicket := func(s string) []int64 {
		var ticket []int64

		for _, i := range strings.Split(s, ",") {
			ticket = append(ticket, MustParse(i))
		}

		return ticket
	}

	var mine []int64
	var nearby [][]int64

	scanMyTicket := func(s *bufio.Scanner) {
		s.Scan()
		if s.Text() != "your ticket:" {
			panic(fmt.Sprintf("expected your ticket!, got %v", s.Text()))
		}

		s.Scan()
		mine = scanTicket(s.Text())

		s.Scan()
		if s.Text() != "" {
			panic(fmt.Sprintf("expected empty, got %v", s.Text()))
		}
	}

	scanNearby := func(s *bufio.Scanner) {
		s.Scan()
		if s.Text() != "nearby tickets:" {
			panic(fmt.Sprintf("expected nearby ticket!, got %v", s.Text()))
		}

		for s.Scan() {
			nearby = append(nearby, scanTicket(s.Text()))
		}
	}

	scanRules(s)
	scanMyTicket(s)
	scanNearby(s)

	log.Printf("pt(1) %v", PartOne(nearby, rules))
	log.Printf("pt(2) %v", PartTwo(mine, nearby, rules))
}
