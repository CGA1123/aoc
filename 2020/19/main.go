package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/CGA1123/aoc"
)

var (
	CharacterRuleExp = regexp.MustCompile(`^(?P<id>\d+): "(?P<char>a|b)"$`)
	OrRuleExp        = regexp.MustCompile(`^(?P<id>\d+): (?P<a>(\d|\s)+) \| (?P<b>(\d|\s)+)$`)
	SequenceRuleExp  = regexp.MustCompile(`^(?P<id>\d+): (?P<a>(\d|\s)+)$`)
)

type State int

const (
	Rules State = iota
	Data
)

type Rule interface {
	Compile(map[int64]Rule) string
}

type CharacterRule struct {
	c byte
}

func (cr *CharacterRule) Compile(rules map[int64]Rule) string {
	return string(cr.c)
}

type SequenceRule struct {
	rules []int64
}

func (seq *SequenceRule) Compile(rules map[int64]Rule) string {
	var str string

	for _, rule := range seq.rules {
		str += rules[rule].Compile(rules)
	}

	return str
}

type OrRule struct {
	a, b []int64
}

func (or *OrRule) Compile(rules map[int64]Rule) string {
	var a, b string

	for _, aRule := range or.a {
		a += rules[aRule].Compile(rules)
	}

	for _, bRule := range or.b {
		b += rules[bRule].Compile(rules)
	}

	return fmt.Sprintf("(%v|%v)", a, b)
}

type RepeatRule struct {
	rule int64
}

func (r *RepeatRule) Compile(rules map[int64]Rule) string {
	return fmt.Sprintf("(%v)+", rules[r.rule].Compile(rules))
}

type BoundedBalancingRule struct {
	left  int64
	right int64
	bound int
}

func (b *BoundedBalancingRule) Compile(rules map[int64]Rule) string {
	var parts []string

	for i := 1; i < b.bound; i++ {
		part := ""
		for j := 0; j < i; j++ {
			part += fmt.Sprintf("(%v)", rules[b.left].Compile(rules))
		}
		for j := 0; j < i; j++ {
			part += fmt.Sprintf("(%v)", rules[b.right].Compile(rules))
		}

		parts = append(parts, fmt.Sprintf("(%v)", part))
	}

	return fmt.Sprintf("(%v)", strings.Join(parts, "|"))
}

var RuleExp = map[*regexp.Regexp]func(string) (int64, Rule){
	CharacterRuleExp: ParseCharacterRule,
	OrRuleExp:        ParseOrRule,
	SequenceRuleExp:  ParseSequenceRule}

func ParseCharacterRule(line string) (int64, Rule) {
	capture := aoc.Capture(CharacterRuleExp, line)

	return aoc.MustParse(capture["id"]), &CharacterRule{c: []byte(capture["char"])[0]}
}

func ParseSequenceRule(line string) (int64, Rule) {
	capture := aoc.Capture(SequenceRuleExp, line)

	var rules []int64

	for _, id := range strings.Split(capture["a"], " ") {
		rules = append(rules, aoc.MustParse(id))
	}

	return aoc.MustParse(capture["id"]), &SequenceRule{rules: rules}
}

func ParseOrRule(line string) (int64, Rule) {
	capture := aoc.Capture(OrRuleExp, line)

	var aRules []int64
	var bRules []int64

	for _, id := range strings.Split(capture["a"], " ") {
		aRules = append(aRules, aoc.MustParse(id))
	}

	for _, id := range strings.Split(capture["b"], " ") {
		bRules = append(bRules, aoc.MustParse(id))
	}

	return aoc.MustParse(capture["id"]), &OrRule{a: aRules, b: bRules}
}

func ParseRule(line string, rules map[int64]Rule) map[int64]Rule {
	for exp, parser := range RuleExp {
		if !exp.MatchString(line) {
			continue
		}

		i, rule := parser(line)
		rules[i] = rule
		return rules
	}

	log.Printf("didn't match rule: %v", line)

	return rules
}

func Matches(rules map[int64]Rule, data []string) int {
	rule := rules[0].Compile(rules)
	exp := regexp.MustCompile(fmt.Sprintf("^%v$", rule))
	var matches int

	for _, line := range data {
		if !exp.MatchString(line) {
			continue
		}

		matches++
	}

	return matches
}

func main() {
	rules := map[int64]Rule{}
	var data []string

	var state State
	aoc.EachLine("input.txt", func(line string) {
		if line == "" {
			state = Data
			return
		}

		if state == Rules {
			rules = ParseRule(line, rules)
		} else {
			data = append(data, line)
		}
	})

	log.Printf("pt(1): %v", Matches(rules, data))

	rules[8] = &RepeatRule{rule: 42}
	rules[11] = &BoundedBalancingRule{bound: 10, left: 42, right: 31}

	log.Printf("pt(2): %v", Matches(rules, data))
}
