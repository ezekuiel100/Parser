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

	for r := range tokens {
		stmt := ParserProgram(tokens[r])
		Statements = append(Statements , stmt)
	}

	fmt.Println(Statements)
}


func (ls LetStatement) Node(){
	fmt.Println(ls)
}


func ParserProgram(token Token) LetStatement {
	if token.Type == "let" {
		return parseLetStatement()
	}

	return LetStatement{}
}

func parseLetStatement() LetStatement {
	return LetStatement{token: "let" , name: "number" , value: "10"}
}