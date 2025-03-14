package lexer

import (
	"sunbird/token"
)

type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Current reading position
	ch           byte // Current char under examination
	// line				 int // TODO: add this in future
	// col 				 int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
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
	return l.input[position:l.position], token.FLOAT
		
	}

	return l.input[position:l.position], token.INT

}
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
		
	case '+':
		tok = newToken(token.PLUS, l.ch)
		
	case '-':
		tok = newToken(token.MINUS, l.ch)
		
	case ';':
			tok = newToken(token.SEMICOLON, l.ch)

	case '!':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
			} else {	
				tok = newToken(token.BANG, l.ch)
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
			tok = newToken(token.SLASH, l.ch)
	
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}

	case '>':	
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}

	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()

			tok = token.Token{ Type: token.PIPE, Literal: string(ch) + string(l.ch)}
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		}
	
	case '(':
			tok = newToken(token.LPAREN, l.ch)
		
	case ')':
			tok = newToken(token.RPAREN, l.ch)
	
	case ',':
			tok = newToken(token.COMMA, l.ch)

	case '{':
		tok = newToken(token.LBRACE, l.ch)

	case '}':
		tok = newToken(token.RBRACE, l.ch)

	case '[':
		tok = newToken(token.LBRACKET, l.ch)

	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	
	case '"', '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // Return earlier because readChar() is already being executed in LookupIdent()
		} else if isDigit(l.ch) {
			literal, tokenType := l.readNumber()


			tok.Literal = literal
			tok.Type = tokenType
			return tok // Return earlier because readChar() is already being executed in readNumber ()

		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}

	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
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
