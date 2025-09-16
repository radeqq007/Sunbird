package lexer

import (
	"sunbird/internal/token"
)

type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Current reading position
	ch           byte // Current char under examination
	line         int  // Current line number
	col          int  // Current column number
}

var keywords = map[string]token.TokenType{
	"func":   token.Function,
	"var":    token.Var,
	"true":   token.True,
	"false":  token.False,
	"if":     token.If,
	"else":   token.Else,
	"return": token.Return,
	"null":   token.Null,
	"for":    token.For,
	"while":  token.While,
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 0, col: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1

	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	startingQuote := l.ch

	l.readChar() // skip the starting quote

	position := l.position

	for l.ch != startingQuote {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		l.readChar()

		for isDigit(l.ch) {
			l.readChar()
		}
		return l.input[position:l.position], token.Float
	}

	return l.input[position:l.position], token.Int
}

func (l *Lexer) NextToken() token.Token {
	pos := token.Position{
		Filename: "",
		Offset:   l.position,
		Line:     l.line,
		Col:      l.col - 1,
	}

	tok := token.Token{
		Pos: pos,
	}

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.Eq, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.Assign, l.ch, pos)
		}

	case '+':
		tok = newToken(token.Plus, l.ch, pos)

	case '-':
		tok = newToken(token.Minus, l.ch, pos)

	case ';':
		tok = newToken(token.Semicolon, l.ch, pos)

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NotEq, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.Bang, l.ch, pos)
		}

	case '/':
		if l.peekChar() == '/' { // Comment handling
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}

			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.readChar() // skip the /
			l.readChar() // skip the *

			for l.ch != '*' && l.peekChar() != '/' {
				l.readChar()
			}

			l.readChar() // skip the *
			l.readChar() // skip the /
			return l.NextToken()
		}
		tok = newToken(token.Slash, l.ch, pos)

	case '*':
		tok = newToken(token.Asterisk, l.ch, pos)

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch, pos)
		}

	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch, pos)
		}

	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.Or, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()

			tok = token.Token{Type: token.Pipe, Literal: string(ch) + string(l.ch)}
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.And, Literal: string(ch) + string(l.ch)}
		}

	case '(':
		tok = newToken(token.LParen, l.ch, pos)

	case ')':
		tok = newToken(token.RParen, l.ch, pos)

	case ',':
		tok = newToken(token.Comma, l.ch, pos)

	case '{':
		tok = newToken(token.LBrace, l.ch, pos)

	case '}':
		tok = newToken(token.RBrace, l.ch, pos)

	case '[':
		tok = newToken(token.LBracket, l.ch, pos)

	case ']':
		tok = newToken(token.RBracket, l.ch, pos)

	case '"', '\'':
		tok.Type = token.String
		tok.Literal = l.readString()

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		switch {
		case isLetter(l.ch):
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok // Return earlier because readChar() is already being executed in LookupIdent()

		case isDigit(l.ch):
			literal, tokenType := l.readNumber()

			tok.Literal = literal
			tok.Type = tokenType
			return tok // Return earlier because readChar() is already being executed in readNumber ()

		default:
			tok = newToken(token.Illegal, l.ch, pos)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, posistion token.Position) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Pos: posistion}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func LookupIdent(ident string) token.TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return token.Ident
}