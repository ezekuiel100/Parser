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
		{Type: "eof", Value: "10"},
	}

	Statements := []Statement{}

	stmt := ParserProgram(tokens)

	if stmt != nil {
		Statements = append(Statements, *stmt)
		fmt.Printf("%+v\n", Statements)
	}

}

func (ls LetStatement) Node() {
	fmt.Println(ls)
}

type Parser struct {
	token  []Token
	errors []string

	curToken  Token
	peekToken Token
	position  int
}

func (p *Parser) advanceToken() {
	p.position++
}

func ParserProgram(tokens []Token) *LetStatement {
	p := &Parser{token: tokens, errors: []string{}, position: 0}
	p.curToken = tokens[p.position]
	p.peekToken = tokens[p.position+1]

	switch p.curToken.Type {
	case "let":
		return parseLetStatement(p)
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

func parseLetStatement(p *Parser) *LetStatement {
	p.advanceToken()

	if p.token[p.position].Type != "identifier" {
		return nil
	}

	identifier := p.token[p.position].Value

	p.advanceToken()

	if p.token[p.position].Type == "equal" {
		p.advanceToken()
	} else {
		p.peekError("equal")
		return nil
	}

	value := p.token[p.position].Value

	return &LetStatement{token: "let", name: identifier, value: value}
}
