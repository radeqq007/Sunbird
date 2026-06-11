package ast

import (
	"bytes"
	"strings"

	"github.com/radeqq007/sunbird/internal/token"
)

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" ")
	out.WriteString(oe.Operator)
	out.WriteString(" ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

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

type CompoundAssignExpression struct {
	Token    token.Token // the token.ASSIGN token
	Name     Expression  // The identifier or property expression being assigned to
	Operator string
	Value    Expression
}

func (cae *CompoundAssignExpression) expressionNode()      {}
func (cae *CompoundAssignExpression) TokenLiteral() string { return cae.Token.Literal }
func (cae *CompoundAssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(cae.Name.String())
	out.WriteString(" ")
	out.WriteString(cae.Operator)
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString(cae.Value.String())
	out.WriteString(";")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type DeclarationExpression struct {
	Token   token.Token
	Name    Expression
	IsConst bool
	Value   Expression
}

func (ds *DeclarationExpression) expressionNode()      {}
func (ds *DeclarationExpression) TokenLiteral() string { return ds.Token.Literal }

func (ds *DeclarationExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ds.Name.String())

	if ds.IsConst {
		out.WriteString(" :: ")
	} else {
		out.WriteString(" := ")
	}

	if ds.Value != nil {
		out.WriteString(ds.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type RangeExpression struct {
	Token token.Token
	Start Expression
	End   Expression
	Step  Expression // 0..10:2 syntax, optional
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	var out bytes.Buffer

	out.WriteString(re.Start.String())
	out.WriteString("..")
	out.WriteString(re.End.String())
	if re.Step != nil {
		out.WriteString(":")
		out.WriteString(re.Step.String())
	}
	return out.String()
}
