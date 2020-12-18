package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CGA1123/aoc"
)

type lexeme int

const (
	undefined lexeme = iota
	integer
	operation
	lparen
	rparen
)

func (l lexeme) String() string {
	return map[lexeme]string{
		undefined: "undefined",
		integer:   "integer",
		operation: "operation",
		lparen:    "lparen",
		rparen:    "rparen"}[l]
}

func TokenType(s string) lexeme {
	switch s {
	case "+":
		return operation
	case "*":
		return operation
	case "(":
		return lparen
	case ")":
		return rparen
	}

	return integer
}

type Token struct {
	Type  lexeme
	Value string
}

type Expression struct {
	token Token
	left  *Expression
	right *Expression
}

func (e *Expression) String() string {
	switch e.token.Type {
	case integer:
		return e.token.Value
	case operation:
		return fmt.Sprintf("(%v %v %v)", e.left.String(), e.token.Value, e.right.String())
	default:
		panic("nope")
	}
}

func (e *Expression) Eval() int64 {
	if e.token.Type == integer {
		return aoc.MustParse(e.token.Value)
	}

	switch e.token.Value {
	case "*":
		return e.left.Eval() * e.right.Eval()
	case "+":
		return e.left.Eval() + e.right.Eval()
	case "(":
		return e.left.Eval()
	}

	panic("Eval")
}

func Lex(exp []string) []Token {
	var tokens []Token
	for _, b := range exp {
		tokens = append(tokens, Token{Type: TokenType(b), Value: b})
	}
	return tokens
}

type Parser struct {
	t []Token
	i int
}

func (p *Parser) token() Token {
	return p.t[p.i]
}

func (p *Parser) lex(expected lexeme) Token {
	p.i += 1

	return p.t[p.i-1]
}

func (p *Parser) Parse() *Expression {
	return p.parseRight(p.parseSubExpRight())
}

func (p *Parser) parseSubExpRight() *Expression {
	token := p.lex(undefined)

	if token.Type == integer {
		return &Expression{token: token}
	}

	exp := p.parseRight(p.parseSubExpRight())
	p.lex(rparen)

	return exp
}

func (p *Parser) parseRight(left *Expression) *Expression {
	if p.i >= len(p.t) || p.token().Type == rparen {
		return left
	}

	exp := &Expression{left: left,
		token: p.lex(operation),
		right: p.parseSubExpRight()}

	return p.parseRight(exp)
}

func main() {
	var partOne int64
	var partTwo int64

	aoc.EachLine("input.txt", func(l string) {
		exp := strings.ReplaceAll(strings.ReplaceAll(l, "(", "( "), ")", " )")

		parserOne := &Parser{t: Lex(strings.Split(exp, " "))}
		one := parserOne.Parse()
		partOne += one.Eval()

		// lol
		two := strings.ReplaceAll(exp, "(", "( ( ( (")
		two = strings.ReplaceAll(two, ")", ") ) ) )")
		two = strings.ReplaceAll(two, "+", ") ) + ( (")
		two = strings.ReplaceAll(two, "*", ") ) ) * ( ( (")
		two = fmt.Sprintf("( ( ( ( %v ) ) ) )", two)
		parserTwo := &Parser{t: Lex(strings.Split(two, " "))}
		t := parserTwo.Parse()
		partTwo += t.Eval()
	})

	log.Printf("pt(1) %v", partOne)
	log.Printf("pt(2) %v", partTwo)
}
