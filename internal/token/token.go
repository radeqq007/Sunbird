package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string

	Pos Position
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF" // End of file

	// Identifiers + literals
	IDENT  = "IDENT"
	FLOAT  = "FLOAT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PIPE     = "|>"

	// Comparison operators
	EQ     = "=="
	NOT_EQ = "!="

	LT = "<"
	GT = ">"

	LE = "<="
	GE = ">="

	// Logical operators
	OR  = "||"
	AND = "&&"

	// Delimeter
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	NULL     = "NULL"
	FOR      = "FOR"
	WHILE    = "WHILE"
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
