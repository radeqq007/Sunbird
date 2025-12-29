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

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }

func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())

	return out.String()
}

type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.Token.Literal + ";" }

type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return cs.Token.Literal + ";" }

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

// Import
type ImportStatement struct {
	Token token.Token
	Path  *StringLiteral
	Alias *Identifier
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString("import ")
	out.WriteString(is.Path.String())
	if is.Alias != nil {
		out.WriteString(" as ")
		out.WriteString(is.Alias.String())
	}
	out.WriteString(";")
	return out.String()
}
