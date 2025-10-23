package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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

type Identifier struct {
	Token Token
	Value string
}

type Expression interface {
	ExpressionNode()
}

const (
	_ int = iota
	LOWEST
	EQUALS      //==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     //*
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

func main() {
	tokens := []Token{
		{Type: "let", Value: "let"},
		{Type: "identifier", Value: "number"},
		{Type: "equal", Value: "="},
		{Type: "int", Value: "10"},
		{Type: "eol", Value: ""},
		{Type: "return", Value: "9"},
		{Type: "eof", Value: ""},
	}

	Statements := []Statement{}

	p := &Parser{tokens: tokens, errors: []string{}, position: 0}
	p.curToken = tokens[p.position]
	p.position++
	p.peekToken = tokens[p.position]

	p.prefixParseFns = make(map[string]func() Expression)
	p.resgisterPrefix("identifier", p.parseIdentifier)
	p.resgisterPrefix("int", p.parseIntegerLiteral)
	p.resgisterPrefix("bang", p.parsePrefixExpression)
	p.resgisterPrefix("minus", p.parsePrefixExpression)

	p.infixParseFns = make(map[string]func(Expression) Expression)
	p.registerInfix("plus", p.parseInfixExpression)
	p.registerInfix("minus", p.parseInfixExpression)
	p.registerInfix("slash", p.parseInfixExpression)
	p.registerInfix("asterisk", p.parseInfixExpression)
	p.registerInfix("equal", p.parseInfixExpression)
	p.registerInfix("no_equal", p.parseInfixExpression)
	p.registerInfix("less_than", p.parseInfixExpression)
	p.registerInfix("greater_than", p.parseInfixExpression)

	for p.curToken.Type != "eof" {
		stmt := p.ParserProgram()

		if stmt != nil {
			Statements = append(Statements, stmt)
		}

		p.advanceToken()
	}

	b, _ := json.MarshalIndent(Statements, "", " ")
	fmt.Println(string(b))
}

type Parser struct {
	tokens []Token
	errors []string

	curToken  Token
	peekToken Token
	position  int

	prefixParseFns map[string]func() Expression
	infixParseFns  map[string]func(Expression) Expression
}

var precedences = map[string]int{
	"equal":        EQUALS,
	"not_equal":    EQUALS,
	"less_than":    LESSGREATER,
	"greater_than": LESSGREATER,
	"plus":         SUM,
	"minus":        SUM,
	"slash":        PRODUCT,
	"asterisk":     PRODUCT,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) advanceToken() {
	p.position++
	p.curToken = p.peekToken

	if p.position >= len(p.tokens) {
		p.peekToken = p.curToken
	} else {
		p.peekToken = p.tokens[p.position]
	}
}

func (p *Parser) ParserProgram() Statement {
	switch p.curToken.Type {
	case "let":
		return p.parseLetStatement()
	case "return":
		return p.parseReturnStatement()
	default:
		return p.parserExpressionStatement()
	}

}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t string) {
	msg := fmt.Sprintf("expexted next token to be %s, got %s instead", t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t string) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (Identifier) ExpressionNode() {}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Value}
}

type LetStatement struct {
	Token Token
	Name  string
	Value string
}

func (LetStatement) Node() {}

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

	return &LetStatement{Token: p.curToken, Name: identifier, Value: value}
}

type ReturnStatement struct {
	Token       Token
	ReturnValue string
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.advanceToken()

	//handle expression

	return &ReturnStatement{Token: p.curToken, ReturnValue: p.curToken.Value}
}

func (ReturnStatement) Node() {}

type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (ExpressionStatement) Node() {}

func (p *Parser) parserExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == "eol" {
		p.advanceToken()
	}

	return stmt
}

func (p *Parser) resgisterPrefix(tokenType string, fn func() Expression) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType string, fn func(Expression) Expression) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	return leftExp
}

type IntergerLiteral struct {
	Token Token
	Value int64
}

func (IntergerLiteral) ExpressionNode() {}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntergerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Value, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (PrefixExpression) ExpressionNode() {}

func (p *Parser) parsePrefixExpression() Expression {
	expression := PrefixExpression{Token: p.curToken, Operator: p.curToken.Value}

	p.advanceToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression

}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (InfixExpression) ExpressionNode() {}

func (p *Parser) parseInfixOperator() Expression {
	infix := InfixExpression{Token: p.curToken, Operator: p.curToken.Value}

	return infix
}
