package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math"
	"strconv"
	"strings"
	"sunbird/internal/ast"
	"unsafe"
)

// ApplyFunction is a hook to allow calling functions from modules
var ApplyFunction func(fn Value, args []Value) Value

type ValueKind uint8

const (
	IntKind ValueKind = iota
	FloatKind
	BoolKind
	NullKind
	StringKind
	ArrayKind
	HashKind
	FunctionKind
	BuiltinKind
	ReturnValueKind
	ErrorKind
	BreakKind
	ContinueKind
	RangeKind
	ModuleKind
)

func (vk ValueKind) String() string {
	switch vk {
	case IntKind:
		return "Integer"
	case FloatKind:
		return "Float"
	case BoolKind:
		return "Boolean"
	case NullKind:
		return "Null"
	case StringKind:
		return "String"
	case ArrayKind:
		return "Array"
	case HashKind:
		return "Hash"
	case FunctionKind:
		return "Function"
	case BuiltinKind:
		return "Builtin"
	case ReturnValueKind:
		return "ReturnValue"
	case ErrorKind:
		return "Error"
	case BreakKind:
		return "Break"
	case ContinueKind:
		return "Continue"
	case RangeKind:
		return "Range"
	case ModuleKind:
		return "Module"
	default:
		return "Unknown"
	}
}

// Primitives (int, float, bool, null) are stored inline
// Complex types use ptr to heap-allocated data
type Value struct {
	kind ValueKind
	bits uint64
	ptr  unsafe.Pointer
}

type String struct {
	Value string
}

type Array struct {
	Elements []Value
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Proto *Hash
}

type HashKey struct {
	Kind  ValueKind
	Value uint64
}

type HashPair struct {
	Key   Value
	Value Value
}

type Function struct {
	Parameters []*ast.Identifier
	ReturnType ast.TypeAnnotation
	Body       *ast.BlockStatement
	Env        *Environment
}

type BuiltinFunction func(args ...Value) Value

type Builtin struct {
	Fn BuiltinFunction
}

type ReturnValue struct {
	Value Value
}

type Error struct {
	Message     string
	Line        int
	Col         int
	Propagating bool
}

type Range struct {
	Start int64
	End   int64
	Step  int64
}

type Module struct {
	Name    string
	Exports map[string]Value
}

func (v Value) Kind() ValueKind {
	return v.kind
}

func (v Value) Inspect() string {
	switch v.kind {
	case IntKind:
		return strconv.FormatInt(v.AsInt(), 10)

	case FloatKind:
		return strconv.FormatFloat(v.AsFloat(), 'f', -1, 64)

	case BoolKind:
		return strconv.FormatBool(v.AsBool())

	case NullKind:
		return "null"

	case StringKind:
		return `"` + v.AsString().Value + `"`

	case ArrayKind:
		arr := v.AsArray()
		var out bytes.Buffer
		elements := make([]string, 0, len(arr.Elements))
		for _, e := range arr.Elements {
			elements = append(elements, e.Inspect())
		}
		out.WriteString("[")
		out.WriteString(strings.Join(elements, ", "))
		out.WriteString("]")
		return out.String()

	case HashKind:
		h := v.AsHash()
		var out bytes.Buffer
		pairs := []string{}
		for _, pair := range h.Pairs {
			pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
		}
		out.WriteString("{")
		out.WriteString(strings.Join(pairs, ", "))
		out.WriteString("}")
		return out.String()

	case FunctionKind:
		fn := v.AsFunction()
		var out bytes.Buffer
		params := []string{}
		for _, p := range fn.Parameters {
			params = append(params, p.String())
		}
		out.WriteString("func")
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(") {\n")
		out.WriteString(fn.Body.String())
		out.WriteString("\n}")
		return out.String()

	case BuiltinKind:
		return "builtin function"

	case ReturnValueKind:
		rv := v.AsReturnValue()
		return rv.Value.Inspect()

	case ErrorKind:
		err := v.AsError()
		if err.Line > 0 {
			return fmt.Sprintf("%s (at line %d, col %d)", err.Message, err.Line, err.Col)
		}
		return err.Message

	case BreakKind:
		return "break"

	case ContinueKind:
		return "continue"

	case RangeKind:
		r := v.AsRange()
		if r.Step != 1 {
			return fmt.Sprintf("%d..%d:%d", r.Start, r.End, r.Step)
		}
		return fmt.Sprintf("%d..%d", r.Start, r.End)

	case ModuleKind:
		m := v.AsModule()
		return "<module " + m.Name + ">"

	default:
		return "unknown"
	}
}

// Hashable interface implementation
func (v Value) HashKey() HashKey {
	switch v.kind {
	case IntKind:
		return HashKey{Kind: IntKind, Value: v.bits}

	case StringKind:
		h := fnv.New64a()
		_, _ = h.Write([]byte(v.AsString().Value))
		return HashKey{Kind: StringKind, Value: h.Sum64()}

	default:
		// Other types can't be used as hash keys
		panic(fmt.Sprintf("type %s is not hashable", v.kind))
	}
}

func (v Value) IsInt() bool      { return v.kind == IntKind }
func (v Value) IsFloat() bool    { return v.kind == FloatKind }
func (v Value) IsBool() bool     { return v.kind == BoolKind }
func (v Value) IsNull() bool     { return v.kind == NullKind }
func (v Value) IsString() bool   { return v.kind == StringKind }
func (v Value) IsArray() bool    { return v.kind == ArrayKind }
func (v Value) IsHash() bool     { return v.kind == HashKind }
func (v Value) IsFunction() bool { return v.kind == FunctionKind }
func (v Value) IsBuiltin() bool  { return v.kind == BuiltinKind }
func (v Value) IsError() bool    { return v.kind == ErrorKind }
func (v Value) IsRange() bool    { return v.kind == RangeKind }
func (v Value) IsModule() bool   { return v.kind == ModuleKind }

// Getters
func (v Value) AsInt() int64 {
	return int64(v.bits)
}

func (v Value) AsFloat() float64 {
	return math.Float64frombits(v.bits)
}

func (v Value) AsBool() bool {
	return v.bits != 0
}

func (v Value) AsString() *String {
	return (*String)(v.ptr)
}

func (v Value) AsArray() *Array {
	return (*Array)(v.ptr)
}

func (v Value) AsHash() *Hash {
	return (*Hash)(v.ptr)
}

func (v Value) AsFunction() *Function {
	return (*Function)(v.ptr)
}

func (v Value) AsBuiltin() *Builtin {
	return (*Builtin)(v.ptr)
}

func (v Value) AsReturnValue() *ReturnValue {
	return (*ReturnValue)(v.ptr)
}

func (v Value) AsError() *Error {
	return (*Error)(v.ptr)
}

func (v Value) AsRange() *Range {
	return (*Range)(v.ptr)
}

func (v Value) AsModule() *Module {
	return (*Module)(v.ptr)
}

func NewInt(val int64) Value {
	return Value{kind: IntKind, bits: uint64(val)}
}

func NewFloat(val float64) Value {
	return Value{kind: FloatKind, bits: math.Float64bits(val)}
}

func NewBool(val bool) Value {
	var bits uint64 = 0
	if val {
		bits = 1
	}
	return Value{kind: BoolKind, bits: bits}
}

func NewNull() Value {
	return Value{kind: NullKind}
}

func NewString(val string) Value {
	s := &String{Value: val}
	return Value{
		kind: StringKind,
		ptr:  unsafe.Pointer(s),
	}
}

func NewArray(elements []Value) Value {
	arr := &Array{Elements: elements}
	return Value{
		kind: ArrayKind,
		ptr:  unsafe.Pointer(arr),
	}
}

func NewHash(pairs map[HashKey]HashPair) Value {
	h := &Hash{Pairs: pairs}
	return Value{
		kind: HashKind,
		ptr:  unsafe.Pointer(h),
	}
}

func NewHashPair(key, value Value) HashPair {
	switch key.Kind() {
	case IntKind, StringKind:
		return HashPair{
			Key:   key,
			Value: value,
		}
	default:
		panic(fmt.Sprintf("type %s is not hashable", key.Kind()))
	}
}

func NewFunction(
	parameters []*ast.Identifier,
	returnType ast.TypeAnnotation,
	body *ast.BlockStatement,
	env *Environment,
) Value {
	fn := &Function{
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
		Env:        env,
	}

	return Value{
		kind: FunctionKind,
		ptr:  unsafe.Pointer(fn),
	}
}

func NewBuiltin(fn BuiltinFunction) Value {
	b := &Builtin{Fn: fn}
	return Value{
		kind: BuiltinKind,
		ptr:  unsafe.Pointer(b),
	}
}

func NewReturnValue(val Value) Value {
	rv := &ReturnValue{Value: val}
	return Value{
		kind: ReturnValueKind,
		ptr:  unsafe.Pointer(rv),
	}
}

func NewError(message string, line, col int, propagating bool) Value {
	err := &Error{
		Message:     message,
		Line:        line,
		Col:         col,
		Propagating: propagating,
	}
	return Value{
		kind: ErrorKind,
		ptr:  unsafe.Pointer(err),
	}
}

func NewBreak() Value {
	return Value{kind: BreakKind}
}

func NewContinue() Value {
	return Value{kind: ContinueKind}
}

func NewRange(start, end, step int64) Value {
	r := &Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	return Value{
		kind: RangeKind,
		ptr:  unsafe.Pointer(r),
	}
}

func NewModule(name string, exports map[string]Value) Value {
	m := &Module{
		Name:    name,
		Exports: exports,
	}
	return Value{
		kind: ModuleKind,
		ptr:  unsafe.Pointer(m),
	}
}
