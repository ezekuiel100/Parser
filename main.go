package main

import (
	"fmt"
)

type Program struct {
	Statements []string
}

type Token struct {
	Type  string
	Value string
}

type Statement interface {
	Node()
}

type LetStatement struct {
	token string
	name  string
	value string
}

func main() {
	tokens := []Token{
		{Type: "let", Value: "let"},
		{Type: "identifier", Value: "number"},
		{Type: "equal", Value: "="},
		{Type: "int", Value: "10"},
		{Type: "return", Value: "9"},
		{Type: "eof", Value: ""},
	}

	Statements := []Statement{}

	p := &Parser{tokens: tokens, errors: []string{}, position: 0}
	p.curToken = tokens[p.position]
	p.position++
	p.peekToken = tokens[p.position]

	for p.curToken.Type != "eof" {
		stmt := p.ParserProgram()

		if stmt != nil {
			Statements = append(Statements, stmt)
		}

		p.advanceToken()
	}

	fmt.Printf("%+v\n", Statements)
}

func (ls LetStatement) Node() {
	fmt.Println(ls)
}

type Parser struct {
	tokens []Token
	errors []string

	curToken  Token
	peekToken Token
	position  int
}

func (p *Parser) advanceToken() {
	p.position++
	p.curToken = p.peekToken
	p.peekToken = p.tokens[p.position]
}

func (p *Parser) ParserProgram() Statement {
	switch p.curToken.Type {
	case "let":
		return p.parseLetStatement()
	}

	return nil
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t string) {
	msg := fmt.Sprintf("expexted next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseLetStatement() *LetStatement {
	p.advanceToken()

	if p.curToken.Type != "identifier" {
		p.peekError("identifier")
		return nil
	}

	identifier := p.curToken.Value

	p.advanceToken()

	if p.curToken.Type == "equal" {
		p.advanceToken()
	} else {
		p.peekError("equal")
		return nil
	}

	value := p.curToken.Value

	return &LetStatement{token: "let", name: identifier, value: value}
}
