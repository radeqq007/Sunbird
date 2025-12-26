package ast

import (
	"bytes"
	"sunbird/internal/token"
)

type PropertyAssignStatement struct {
	Statement
	Token    token.Token // identifier
	Object   Expression
	Property *Identifier
	Value    Expression
}

func (pas *PropertyAssignStatement) statementNode()       {}
func (pas *PropertyAssignStatement) TokenLiteral() string { return pas.Token.Literal }
func (pas *PropertyAssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(pas.Object.String())
	out.WriteString(".")
	out.WriteString(pas.Property.String())
	out.WriteString(" = ")
	if pas.Value != nil {
		out.WriteString(pas.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
