package ast

import (
	"bytes"
	"sunbird/internal/token"
)

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs []HashPair
}

type HashPair struct {
	Key   Expression
	Value Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, pair := range hl.Pairs {
		out.WriteString(pair.Key.String())
		out.WriteString(": ")
		out.WriteString(pair.Value.String())
		out.WriteString(", ")
	}
	out.WriteString("}")

	return out.String()
}
