package token

type TokenType uint8

type Token struct {
	Type    TokenType
	Literal string

	Pos Position
}

const (
	Illegal TokenType = iota
	EOF               // End of file

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
	Modulo
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

func (tt TokenType) String() string {
	switch tt {
	case Illegal:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case Ident:
		return "IDENT"
	case Float:
		return "FLOAT"
	case Int:
		return "INT"
	case String:
		return "STRING"
	case Assign:
		return "="
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Bang:
		return "!"
	case Asterisk:
		return "*"
	case Slash:
		return "/"
	case Modulo:
		return "%"
	case Pipe:
		return "|>"
	case Eq:
		return "=="
	case NotEq:
		return "!="
	case LT:
		return "<"
	case GT:
		return ">"
	case LE:
		return "<="
	case GE:
		return ">="
	case Or:
		return "||"
	case And:
		return "&&"
	case Comma:
		return ","
	case Semicolon:
		return ";"
	case LParen:
		return "("
	case RParen:
		return ")"
	case LBrace:
		return "{"
	case RBrace:
		return "}"
	case LBracket:
		return "["
	case RBracket:
		return "]"
	case Function:
		return "FUNCTION"
	case Var:
		return "VAR"
	case True:
		return "TRUE"
	case False:
		return "FALSE"
	case If:
		return "IF"
	case Else:
		return "ELSE"
	case Return:
		return "RETURN"
	case Null:
		return "NULL"
	case For:
		return "FOR"
	case While:
		return "WHILE"
	default:
		return "UNKNOWN"
	}
}

