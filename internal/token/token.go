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
	Comma        TokenType = "COMMA"
	Semicolon    TokenType = "SEMICOLON"
	Colon        TokenType = "COLON"
	Dot          TokenType = "DOT"
	QuestionMark TokenType = "QUESTION_MARK"

	LParen   TokenType = "LPAREN"
	RParen   TokenType = "RPAREN"
	LBrace   TokenType = "LBRACE"
	RBrace   TokenType = "RBRACE"
	LBracket TokenType = "LBRACKET"
	RBracket TokenType = "RBRACKET"
	Pipe     TokenType = "PIPE"

	// Type keywords
	TypeInt    TokenType = "TYPE_INT"
	TypeFloat  TokenType = "TYPE_FLOAT"
	TypeString TokenType = "TYPE_STRING"
	TypeBool   TokenType = "TYPE_BOOL"
	TypeVoid   TokenType = "TYPE_VOID"
	TypeArray  TokenType = "TYPE_ARRAY"
	TypeFunc   TokenType = "TYPE_FUNC"
	TypeHash   TokenType = "TYPE_HASH"

	// Keywords
	Function TokenType = "FUNC"
	Let      TokenType = "LET"
	Const    TokenType = "CONST"
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
	Try      TokenType = "TRY"
	Catch    TokenType = "CATCH"
	Finally  TokenType = "FINALLY"
)

func (t Token) String() string {
	return fmt.Sprintf("%v %v", t.Type, t.Literal)
}
