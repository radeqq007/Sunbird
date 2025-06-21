package lexer_test

import (
	"sunbird/internal/lexer"
	"sunbird/internal/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
	var ten = 10;

var add = func(x, y) {
  x + y;
};

var result = add(five, ten);
!-/ *5;
5 <= 10 >= 5;
2.2 > 0;
if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;

"hello world"
'hi mom!'
"I'm using ' inside double quote string"
'this is a " inside single quote string'

[1, 2];
||
&&
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Var, "var"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Var, "var"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Var, "var"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "func"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Var, "var"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Int, "5"},
		{token.LE, "<="},
		{token.Int, "10"},
		{token.GE, ">="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Float, "2.2"},
		{token.GT, ">"},
		{token.Int, "0"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.LParen, "("},
		{token.Int, "5"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Else, "else"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.String, "hello world"},
		{token.String, "hi mom!"},
		{token.String, "I'm using ' inside double quote string"},
		{token.String, "this is a \" inside single quote string"},
		{token.LBracket, "["},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.RBracket, "]"},
		{token.Semicolon, ";"},
		{token.Or, "||"},
		{token.And, "&&"},
		{token.EOF, ""},
	}
	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i,
				tt.expectedType,
				tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expected=%q, got=%q",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}
	}
}
