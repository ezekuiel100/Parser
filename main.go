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
	token    []Token
	position int
}

func (p *Parser) advanceToken() {
	p.position++
}

func ParserProgram(tokens []Token) *LetStatement {
	p := &Parser{token: tokens, position: 0}
	curToken := tokens[p.position]

	switch curToken.Type {
	case "let":
		return parseLetStatement(p)
	}

	return nil
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
		fmt.Printf("Erro de sintaxe: esperado token 'equal' \n")
		return nil
	}

	value := p.token[p.position].Value

	return &LetStatement{token: "let", name: identifier, value: value}
}
