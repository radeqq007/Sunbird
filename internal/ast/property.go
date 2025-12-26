package ast

import (
	"bytes"
	"sunbird/internal/token"
)

type PropertyExpression struct {
	Token    token.Token
	Object   Expression
	Property *Identifier
}

func (pe *PropertyExpression) expressionNode()      {}
func (pe *PropertyExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PropertyExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Object.String())
	out.WriteString(".")
	out.WriteString(pe.Property.String())
	out.WriteString(")")
	return out.String()
}
