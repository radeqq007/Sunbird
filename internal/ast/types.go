package ast

import (
	"bytes"
	"strings"
	"sunbird/internal/token"
)

type TypeAnnotation interface {
	Node
	typeAnnotation()
}

// SimpleType is for strings, ints, bools etc.
type SimpleType struct {
	Token token.Token
	Name  string
}

func (st *SimpleType) typeAnnotation()      {}
func (st *SimpleType) TokenLiteral() string { return st.Token.Literal }
func (st *SimpleType) String() string       { return st.Name }

type ArrayType struct {
	Token       token.Token
	ElementType TypeAnnotation
}

func (at *ArrayType) typeAnnotation()      {}
func (at *ArrayType) TokenLiteral() string { return at.Token.Literal }
func (at *ArrayType) String() string {
	// return "[" +  "]" + at.ElementType.String()
	return "Array"
}

type HashType struct {
	Token token.Token
}

func (ht *HashType) typeAnnotation()      {}
func (ht *HashType) TokenLiteral() string { return ht.Token.Literal }
func (ht *HashType) String() string       { return "Hash" }

type FunctionType struct {
	Token      token.Token
	Parameters []TypeAnnotation
	ReturnType TypeAnnotation
}

func (ft *FunctionType) typeAnnotation()      {}
func (ft *FunctionType) TokenLiteral() string { return ft.Token.Literal }
func (ft *FunctionType) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ft.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	if ft.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(ft.ReturnType.String())
	}

	return out.String()
}

// OptionalType is for optional types like int?
type OptionalType struct {
	Token    token.Token
	BaseType TypeAnnotation
}

func (ot *OptionalType) typeAnnotation()      {}
func (ot *OptionalType) TokenLiteral() string { return ot.Token.Literal }
func (ot *OptionalType) String() string {
	return ot.BaseType.String() + "?"
}
