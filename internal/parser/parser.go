package parser

import (
	"fmt"
	"sunbird/internal/ast"
	"sunbird/internal/lexer"
	"sunbird/internal/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	LOGICAL     // && or ||
	EQUALS      // ==
	LESSGREATER // >, <, <= or >=
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // foo()
	INDEX       // arr[x]
	PROPERTY    // obj.prop
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Let, p.parseLetExpression)
	p.registerPrefix(token.Const, p.parseConstExpression)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.Float, p.parseFloatLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Function, p.parseFunctionLiteral)
	p.registerPrefix(token.LBracket, p.parseArrayLiteral)
	p.registerPrefix(token.LBrace, p.parseHashLiteral)
	p.registerPrefix(token.Null, p.parseNullLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Modulo, p.parseInfixExpression)
	p.registerInfix(token.Eq, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.PlusEqual, p.parseCompoundAssignExpression)
	p.registerInfix(token.MinusEqual, p.parseCompoundAssignExpression)
	p.registerInfix(token.AsteriskEqual, p.parseCompoundAssignExpression)
	p.registerInfix(token.SlashEqual, p.parseCompoundAssignExpression)
	p.registerInfix(token.ModuloEqual, p.parseCompoundAssignExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.Or, p.parseInfixExpression)
	p.registerInfix(token.And, p.parseInfixExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)
	p.registerInfix(token.LBracket, p.parseIndexExpression)
	p.registerInfix(token.Dot, p.parsePropertyExpression)
	p.registerInfix(token.Assign, p.parseAssignExpression)
	p.registerInfix(token.DotDot, p.parseRangeExpression)

	// Read 2 tokens so curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead (at line %d, col %d)",
		t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Col)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf(
		"no prefix parse function for %s found (at line %d, col %d)",
		t,
		p.curToken.Line,
		p.curToken.Col,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		program.Statements = append(program.Statements, stmt)

		p.nextToken()
	}

	return program
}

func (p *Parser) parseTypeAnnotation() ast.TypeAnnotation {
	switch p.curToken.Type {
	case token.TypeInt, token.TypeFloat, token.TypeString, token.TypeBool, token.TypeVoid:
		t := &ast.SimpleType{
			Token: p.curToken,
			Name:  p.curToken.Literal,
		}

		if p.peekTokenIs(token.QuestionMark) {
			p.nextToken()
			return &ast.OptionalType{
				Token:    p.curToken,
				BaseType: t,
			}
		}

		return t

	case token.TypeArray:
		t := &ast.ArrayType{Token: p.curToken}

		if p.peekTokenIs(token.QuestionMark) {
			p.nextToken()
			return &ast.OptionalType{
				Token:    p.curToken,
				BaseType: t,
			}
		}
		return t

	case token.TypeHash:
		t := &ast.HashType{Token: p.curToken}
		if p.peekTokenIs(token.QuestionMark) {
			p.nextToken()
			return &ast.OptionalType{
				Token:    p.curToken,
				BaseType: t,
			}
		}
		return t

	case token.TypeFunc:
		t := &ast.FunctionType{Token: p.curToken}
		if p.peekTokenIs(token.QuestionMark) {
			p.nextToken()
			return &ast.OptionalType{
				Token:    p.curToken,
				BaseType: t,
			}
		}
		return t

	default:
		return nil
	}
}

func (p *Parser) parseTypeList(end token.TokenType) []ast.TypeAnnotation {
	list := []ast.TypeAnnotation{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseTypeAnnotation())

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseTypeAnnotation())
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
