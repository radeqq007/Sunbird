package ast

import (
	"bytes"
	"sunbird/internal/token"
)

type AssignExpression struct {
	Token token.Token // the token.ASSIGN token
	Name  Expression  // The identifier or property expression being assigned to
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Name.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())
	out.WriteString(";")

	return out.String()
}
