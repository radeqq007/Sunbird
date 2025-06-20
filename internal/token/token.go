package token

type TokenType uint8

type Token struct {
	Type    TokenType
	Literal string

	Pos Position
}

const (
	ILLEGAL TokenType = iota
	EOF               // End of file

	// Identifiers + literals
	IDENT
	FLOAT
	INT
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	PIPE

	// Comparison operators
	EQ
	NOT_EQ

	LT
	GT

	LE
	GE

	// Logical operators
	OR
	AND

	// Delimeter
	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Keywords
	FUNCTION
	VAR
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	NULL
	FOR
	WHILE
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"null":   NULL,
	"for":    FOR,
	"while":  WHILE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
