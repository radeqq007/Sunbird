package lexer

import (
	"fmt"
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
	"func":     token.Function,
	"let":      token.Let,
	"true":     token.True,
	"false":    token.False,
	"if":       token.If,
	"else":     token.Else,
	"return":   token.Return,
	"null":     token.Null,
	"for":      token.For,
	"while":    token.While,
	"import":   token.Import,
	"as":       token.As,
	"break":    token.Break,
	"continue": token.Continue,
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
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

func (l *Lexer) readString() (string, error) {
	startingQuote := l.ch

	l.readChar() // skip the starting quote

	position := l.position

	for l.ch != startingQuote {
		if l.ch == 0 {
			return "", fmt.Errorf("unterminated string")
		}
		l.readChar()
	}

	return l.input[position:l.position], nil
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
	tok := token.Token{}

	l.skipWhitespace()

	startLine := l.line
	startCol := l.col

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.Eq, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		} else {
			tok = l.newToken(token.Assign, l.ch)
		}

	case '+':
		tok = l.newToken(token.Plus, l.ch)

	case '-':
		tok = l.newToken(token.Minus, l.ch)

	case ';':
		tok = l.newToken(token.Semicolon, l.ch)

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NotEq, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		} else {
			tok = l.newToken(token.Bang, l.ch)
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
		tok = l.newToken(token.Slash, l.ch)

	case '*':
		tok = l.newToken(token.Asterisk, l.ch)

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		} else {
			tok = l.newToken(token.LT, l.ch)
		}

	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		} else {
			tok = l.newToken(token.GT, l.ch)
		}

	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.Or, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()

			tok = token.Token{Type: token.Pipe, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.And, Literal: string(ch) + string(l.ch), Line: startLine, Col: startCol}
		}

	case '(':
		tok = l.newToken(token.LParen, l.ch)

	case ')':
		tok = l.newToken(token.RParen, l.ch)

	case ',':
		tok = l.newToken(token.Comma, l.ch)

	case '{':
		tok = l.newToken(token.LBrace, l.ch)

	case '}':
		tok = l.newToken(token.RBrace, l.ch)

	case '[':
		tok = l.newToken(token.LBracket, l.ch)

	case ']':
		tok = l.newToken(token.RBracket, l.ch)

	case ':':
		tok = l.newToken(token.Colon, l.ch)

	case '.':
		tok = l.newToken(token.Dot, l.ch)

	case '"', '\'':
		tok.Type = token.String
		lit, err := l.readString()
		if err != nil {
			tok.Type = token.Illegal
			tok.Literal = err.Error()
		} else {
			tok.Literal = lit
		}
		tok.Line = startLine
		tok.Col = startCol

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = startLine
		tok.Col = startCol

	default:
		switch {
		case isLetter(l.ch):
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = startLine
			tok.Col = startCol
			return tok // Return earlier because readChar() is already being executed in LookupIdent()

		case isDigit(l.ch):
			literal, tokenType := l.readNumber()

			tok.Literal = literal
			tok.Type = tokenType
			tok.Line = startLine
			tok.Col = startCol
			return tok // Return earlier because readChar() is already being executed in readNumber ()

		default:
			tok = l.newToken(token.Illegal, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.line, Col: l.col}
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
