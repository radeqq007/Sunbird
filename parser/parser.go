package parser

import (
	"sunbird/ast"
	"sunbird/lexer"
	"sunbird/token"
)

type Parser struct {
  l         *lexer.Lexer
  curToken  token.Token
  peekToken token.Token
  errors    []string

  prefixParseFns map[token.TokenType]prefixParseFn
  infixParseFns map[token.TokenType]infixParseFn
}

type (
  prefixParseFn func() ast.Expression
  infixParseFn  func(ast.Expression) ast.Expression
)

const  (
  _ int = iota
  LOWEST
  LOGICAL          // && or ||
  EQUALS           // ==
  LESSGREATER      // >, <, <= or >=
  SUM              // +
  PRODUCT          // *
  PREFIX           // -X or !X
  CALL             // foo()
  INDEX            // arr[x]
)

func New(l *lexer.Lexer) *Parser {
  p := &Parser{l: l, errors: []string{}}

  p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
  p.registerPrefix(token.IDENT, p.parseIdentifier)
  p.registerPrefix(token.INT, p.parseIntegerLiteral)
  p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
  p.registerPrefix(token.STRING, p.parseStringLiteral)
  p.registerPrefix(token.BANG, p.parsePrefixExpression)
  p.registerPrefix(token.MINUS, p.parsePrefixExpression)
  p.registerPrefix(token.TRUE, p.parseBoolean)
  p.registerPrefix(token.FALSE, p.parseBoolean)
  p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
  p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
  p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

  p.infixParseFns = make(map[token.TokenType]infixParseFn)
  p.registerInfix(token.PLUS, p.parseInfixExpression)	
  p.registerInfix(token.MINUS, p.parseInfixExpression)
  p.registerInfix(token.SLASH, p.parseInfixExpression)	
  p.registerInfix(token.ASTERISK, p.parseInfixExpression)
  p.registerInfix(token.EQ, p.parseInfixExpression)
  p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
  p.registerInfix(token.LT, p.parseInfixExpression)
  p.registerInfix(token.GT, p.parseInfixExpression)
  p.registerInfix(token.LE, p.parseInfixExpression)
  p.registerInfix(token.GE, p.parseInfixExpression)
  p.registerInfix(token.OR, p.parseInfixExpression)
  p.registerInfix(token.AND, p.parseInfixExpression)
  p.registerInfix(token.LPAREN, p.parseCallExpression)
  p.registerInfix(token.LBRACKET, p.parseIndexExpression)

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