package ast

import (
	"bytes"
	"sunbird/internal/token"
)

// Block
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// For
type ForStatement struct {
	Token     token.Token
	Init      Expression
	Condition Expression
	Update    Expression
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

// Return
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// Property Assign
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
