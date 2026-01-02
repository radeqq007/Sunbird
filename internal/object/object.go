package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sunbird/internal/ast"
)

type ObjectType uint8

const (
	StringObj ObjectType = iota
	IntegerObj
	FloatObj
	BooleanObj
	NullObj
	FunctionObj
	ReturnValueObj
	ErrorObj
	BuiltinObj
	ArrayObj
	HashObj
	BreakObj
	ContinueObj
)

func (ot ObjectType) String() string {
	switch ot {
	case StringObj:
		return "STRING"
	case IntegerObj:
		return "INTEGER"
	case FloatObj:
		return "FLOAT"
	case BooleanObj:
		return "BOOLEAN"
	case NullObj:
		return "NULL"
	case FunctionObj:
		return "FUNCTION"
	case ReturnValueObj:
		return "RETURN_VALUE"
	case ErrorObj:
		return "ERROR"
	case BuiltinObj:
		return "BUILTIN"
	case ArrayObj:
		return "ARRAY"
	case BreakObj:
		return "BREAK"
	case ContinueObj:
		return "CONTINUE"
	default:
		return "UNKNOWN"
	}
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return StringObj }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) Type() ObjectType { return IntegerObj }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (f *Float) Type() ObjectType { return FloatObj }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() ObjectType { return BooleanObj }

type Null struct{}

func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Break struct{}

func (b *Break) Type() ObjectType { return BreakObj }
func (b *Break) Inspect() string  { return "break" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return ContinueObj }
func (c *Continue) Inspect() string  { return "continue" }

type Error struct {
	Message     string
	Line        int
	Col         int
	Propagating bool
}

func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s (at line %d, col %d)", e.Message, e.Line, e.Col)
	}
	return e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	ReturnType ast.TypeAnnotation
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FunctionObj }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BuiltinObj }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ArrayObj }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Proto *Hash // Prototype link for inheritance
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, pair := range h.Pairs {
		out.WriteString(pair.Key.Inspect())
		out.WriteString(": ")
		out.WriteString(pair.Value.Inspect())
		out.WriteString(", ")
	}
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
