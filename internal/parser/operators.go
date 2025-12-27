package parser

import (
	"fmt"
	"sunbird/internal/ast"
	"sunbird/internal/token"
)

// Prefix
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// Infix
var precedences = map[token.TokenType]int{
	token.Assign:   ASSIGN,
	token.Or:       LOGICAL,
	token.And:      LOGICAL,
	token.Eq:       EQUALS,
	token.NotEq:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.Plus:     SUM,
	token.Minus:    SUM,
	token.Slash:    PRODUCT,
	token.Asterisk: PRODUCT,
	token.Modulo:   PRODUCT,
	token.LParen:   CALL,
	token.Pipe:     PIPE,
	token.LBracket: INDEX,
	token.Dot:      PROPERTY,
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Assign
func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	exp := &ast.AssignExpression{
		Token: p.curToken,
		Name:  left,
	}

	if !p.validateAssignmentTarget(left) {
		p.errors = append(p.errors, fmt.Sprintf("invalid assignment target: %s", left.String()))
		return nil
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Value = p.parseExpression(precedence)

	return exp
}

func (p *Parser) validateAssignmentTarget(exp ast.Expression) bool {
	switch exp.(type) {
	case *ast.Identifier, *ast.PropertyExpression, *ast.IndexExpression:
		return true
	}
	return false
}
