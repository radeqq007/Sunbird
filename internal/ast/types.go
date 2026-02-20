package ast

import (
	"sunbird/internal/token"
)

type TypeAnnotation interface {
	Node
	typeAnnotation()
}

// SimpleType is for strings, ints, bools, Range type etc.
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
	return "Func"
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
