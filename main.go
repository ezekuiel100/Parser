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
	Statements = append(Statements, stmt)

	fmt.Printf("%+v\n", Statements)
}

func (ls LetStatement) Node() {
	fmt.Println(ls)
}

func advanteToken(position *int) {
	*position = *position + 1
	fmt.Println(*position)
}

func ParserProgram(tokens []Token) LetStatement {
	position := 0
	curToken := tokens[position]

	advanteToken(&position)

	if curToken.Type == "let" {
		return parseLetStatement()
	}

	return LetStatement{}
}

func parseLetStatement() LetStatement {
	return LetStatement{token: "let", name: "number", value: "10"}
}
