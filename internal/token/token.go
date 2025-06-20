package token

type TokenType uint8

type Token struct {
	Type    TokenType
	Literal string

	Pos Position
}

const (
	Illegal TokenType = iota
	Eof               // End of file

	// Identifiers + literals
	Ident
	Float
	Int
	String

	// Operators
	Assign
	Plus
	Minus
	Bang
	Asterisk
	Slash
	Pipe

	// Comparison operators
	Eq
	NotEq

	LT
	GT

	LE
	GE

	// Logical operators
	Or
	And

	// Delimiter
	Comma
	Semicolon

	LParen
	RParen
	LBrace
	RBrace
	LBracket
	RBracket

	// Keywords
	Function
	Var
	True
	False
	If
	Else
	Return
	Null
	For
	While
)

var keywords = map[string]TokenType{
	"func":   Function,
	"var":    Var,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
	"null":   Null,
	"for":    For,
	"while":  While,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
