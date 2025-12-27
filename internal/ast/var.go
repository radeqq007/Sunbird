package ast

import (
	"bytes"
	"sunbird/internal/token"
)

type VarExpression struct {
	Token token.Token
	Name  Expression
	Value Expression
}

func (vs *VarExpression) expressionNode()      {}
func (vs *VarExpression) TokenLiteral() string { return vs.Token.Literal }

func (vs *VarExpression) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
