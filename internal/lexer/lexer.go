package lexer

import (
	"errors"
	"fmt"
	"strings"
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
	"const":    token.Const,
	"true":     token.True,
	"false":    token.False,
	"if":       token.If,
	"else":     token.Else,
	"return":   token.Return,
	"null":     token.Null,
	"for":      token.For,
	"while":    token.While,
	"import":   token.Import,
	"export":   token.Export,
	"as":       token.As,
	"break":    token.Break,
	"continue": token.Continue,
	"try":      token.Try,
	"catch":    token.Catch,
	"finally":  token.Finally,
	"Int":      token.TypeInt,
	"Float":    token.TypeFloat,
	"String":   token.TypeString,
	"Bool":     token.TypeBool,
	"Void":     token.TypeVoid,
	"Array":    token.TypeArray,
	"Func":     token.TypeFunc,
	"Hash":     token.TypeHash,
	"Range":    token.TypeRange,
	"in":       token.In,
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
	var result strings.Builder

	for l.ch != startingQuote {
		if l.ch == 0 {
			return "", errors.New("unterminated string")
		}

		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case startingQuote:
				result.WriteByte(startingQuote)
			default:
				return "", fmt.Errorf("invalid escape sequence: %c", l.ch)
			}
			l.readChar()
			continue
		}

		result.WriteByte(l.ch)
		l.readChar()
	}

	return result.String(), nil
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	for isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	// this is kinda hacky, but it works
	if l.ch == '.' && l.peekChar() != '.' {
		l.readChar()

		for isDigit(l.ch) || l.ch == '_' {
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
			tok = l.newToken(token.Eq, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Assign, string(l.ch), startLine, startCol)
		}

	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.PlusEqual, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Plus, string(l.ch), startLine, startCol)
		}

	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.MinusEqual, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Minus, string(l.ch), startLine, startCol)
		}

	case ';':
		tok = l.newToken(token.Semicolon, string(l.ch), startLine, startCol)

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.NotEq, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Bang, string(l.ch), startLine, startCol)
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
				if l.ch == 0 {
					return l.newToken(token.Illegal, "unterminated comment", startLine, startCol)
				}
				l.readChar()
			}

			l.readChar() // skip the *
			l.readChar() // skip the /
			return l.NextToken()
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.SlashEqual, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Slash, string(l.ch), startLine, startCol)
		}

	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.AsteriskEqual, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Asterisk, string(l.ch), startLine, startCol)
		}

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.LE, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.LT, string(l.ch), startLine, startCol)
		}

	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.GE, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.GT, string(l.ch), startLine, startCol)
		}

	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.Or, string(ch)+string(l.ch), startLine, startCol)
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.And, string(ch)+string(l.ch), startLine, startCol)
		}

	case '(':
		tok = l.newToken(token.LParen, string(l.ch), startLine, startCol)

	case ')':
		tok = l.newToken(token.RParen, string(l.ch), startLine, startCol)

	case ',':
		tok = l.newToken(token.Comma, string(l.ch), startLine, startCol)

	case '{':
		tok = l.newToken(token.LBrace, string(l.ch), startLine, startCol)

	case '}':
		tok = l.newToken(token.RBrace, string(l.ch), startLine, startCol)

	case '[':
		tok = l.newToken(token.LBracket, string(l.ch), startLine, startCol)

	case ']':
		tok = l.newToken(token.RBracket, string(l.ch), startLine, startCol)

	case ':':
		tok = l.newToken(token.Colon, string(l.ch), startLine, startCol)

	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.DotDot, string(ch)+string(l.ch), startLine, startCol)
		} else {
			tok = l.newToken(token.Dot, string(l.ch), startLine, startCol)
		}

	case '?':
		tok = l.newToken(token.QuestionMark, string(l.ch), startLine, startCol)

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
			tok = l.newToken(token.Illegal, string(l.ch), startLine, startCol)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: line, Col: col}
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
