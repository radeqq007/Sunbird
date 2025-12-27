package parser

import (
	"fmt"
	"sunbird/internal/ast"
)

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
