package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ = "FLOAT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
