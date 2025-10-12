package main

import "fmt"

type Program struct {
	Statements []string
}

type Token struct {
	Type  string
	Value string
}

type Statement struct {
	identifier  string
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



func ParserProgram(token Token) Statement {
	if token.Type == "let" {
		return parseLetStatement()
	}

	return Statement{}
}

func parseLetStatement() Statement {
	return Statement{identifier: "number" , value: "10"}
}