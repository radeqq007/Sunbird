package json

import (
	"encoding/json"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("parse", parseJSON).
		AddFunction("stringify", stringify).
		Build()
}

func parseJSON(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	val := args[0].AsString().Value

	var data any
	errGo := json.Unmarshal([]byte(val), &data)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
	}

	return ToObject(data)
}

func ToObject(val any) object.Value {
	switch v := val.(type) {
	case string:
		return object.NewString(v)
	case float64:
		// JSON numbers are float64 by default in encoding/json
		if v == float64(int64(v)) {
			return object.NewInt(int64(v))
		}
		return object.NewFloat(v)
	case bool:
		return object.NewBool(v)
	case nil:
		return object.NewNull()
	case []any:
		elements := make([]object.Value, len(v))
		for i, el := range v {
			elements[i] = ToObject(el)
		}
		return object.NewArray(elements)
	case map[string]any:
		pairs := make(map[object.HashKey]object.HashPair)
		for k, val := range v {
			key := object.NewString(k)
			hashKey := key.HashKey()
			pairs[hashKey] = object.HashPair{
				Key:   key,
				Value: ToObject(val),
			}
		}
		return object.NewHash(pairs)
	default:
		return object.NewNull()
	}
}

func stringify(ctx object.CallContext, args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	data := FromObject(args[0])
	bytes, errGo := json.Marshal(data)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
	}

	return object.NewString(string(bytes))
}

func FromObject(obj object.Value) any {
	switch obj.Kind() {
	case object.StringKind:
		o := obj.AsString()
		return o.Value
	case object.IntKind:
		return obj.AsInt()
	case object.FloatKind:
		return obj.AsFloat()
	case object.BoolKind:
		return obj.AsBool()
	case object.NullKind:
		return object.NewNull()
	case object.ArrayKind:
		o := obj.AsArray()
		elements := make([]any, len(o.Elements))
		for i, el := range o.Elements {
			elements[i] = FromObject(el)
		}
		return elements
	case object.HashKind:
		o := obj.AsHash()
		m := make(map[string]any)
		for _, pair := range o.Pairs {
			var key string
			if pair.Key.IsString() {
				key = pair.Key.AsString().Value
			} else {
				key = pair.Key.Inspect()
			}
			m[key] = FromObject(pair.Value)
		}
		return m
	default:
		return nil
	}
}
