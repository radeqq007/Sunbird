package ast

import (
	"bytes"
	"sunbird/internal/token"
)

type ForStatement struct {
	Token     token.Token
	Init      Statement
	Condition Expression
	Update    Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }

func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fs.Init.String())
	out.WriteString(" ")
	out.WriteString(fs.Condition.String())
	out.WriteString("; ")
	out.WriteString(fs.Update.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}
