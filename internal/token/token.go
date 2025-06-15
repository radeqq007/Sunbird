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

func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case FLOAT:
		return "FLOAT"
	case INT:
		return "INT"
	case STRING:
		return "STRING"
	case ASSIGN:
		return "="
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case BANG:
		return "!"
	case ASTERISK:
		return "*"
	case SLASH:
		return "/"
	case PIPE:
		return "|>"
	case EQ:
		return "=="
	case NOT_EQ:
		return "!="
	case LT:
		return "<"
	case GT:
		return ">"
	case LE:
		return "<="
	case GE:
		return ">="
	case OR:
		return "||"
	case AND:
		return "&&"
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case FUNCTION:
		return "FUNCTION"
	case VAR:
		return "VAR"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case NULL:
		return "NULL"
	case FOR:
		return "FOR"
	case WHILE:
		return "WHILE"
	default:
		return "UNKNOWN"
	}
}

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
