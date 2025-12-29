package token

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

const (
	Illegal TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	Ident  TokenType = "IDENT"
	Float  TokenType = "FLOAT"
	Int    TokenType = "INT"
	String TokenType = "STRING"

	// Operators
	Assign   TokenType = "ASSIGN"
	Plus     TokenType = "PLUS"
	Minus    TokenType = "MINUS"
	Bang     TokenType = "BANG"
	Modulo   TokenType = "MODULO"
	Asterisk TokenType = "ASTERISK"
	Slash    TokenType = "SLASH"
	Pipe     TokenType = "PIPE"

	// Comparison operators
	Eq    TokenType = "EQ"
	NotEq TokenType = "NOTEQ"

	LT TokenType = "LT"
	GT TokenType = "GT"

	LE TokenType = "LE"
	GE TokenType = "GE"

	// Logical operators
	Or  TokenType = "OR"
	And TokenType = "AND"

	// Delimiter
	Comma     TokenType = "COMMA"
	Semicolon TokenType = "SEMICOLON"
	Colon     TokenType = "COLON"
	Dot       TokenType = "DOT"

	LParen   TokenType = "LPAREN"
	RParen   TokenType = "RPAREN"
	LBrace   TokenType = "LBRACE"
	RBrace   TokenType = "RBRACE"
	LBracket TokenType = "LBRACKET"
	RBracket TokenType = "RBRACKET"

	// Keywords
	Function TokenType = "FUNC"
	Let      TokenType = "LET"
	True     TokenType = "TRUE"
	False    TokenType = "FALSE"
	If       TokenType = "IF"
	Else     TokenType = "ELSE"
	Return   TokenType = "RETURN"
	Null     TokenType = "NULL"
	For      TokenType = "FOR"
	While    TokenType = "WHILE"
	Import   TokenType = "IMPORT"
	As       TokenType = "AS"
	Break    TokenType = "BREAK"
	Continue TokenType = "CONTINUE"
)

func (t Token) String() string {
	return fmt.Sprintf("%v %v", t.Type, t.Literal)
}
