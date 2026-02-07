package json

import (
	"encoding/json"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() *object.Hash {
	return modbuilder.NewModuleBuilder().
		AddFunction("parse", parseJSON).
		AddFunction("stringify", stringify).
		Build()
}

func parseJSON(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	val := args[0].(*object.String).Value

	var data any
	errGo := json.Unmarshal([]byte(val), &data)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
	}

	return ToObject(data)
}

func ToObject(val any) object.Object {
	switch v := val.(type) {
	case string:
		return &object.String{Value: v}
	case float64:
		// JSON numbers are float64 by default in encoding/json
		if v == float64(int64(v)) {
			return &object.Integer{Value: int64(v)}
		}
		return &object.Float{Value: v}
	case bool:
		return &object.Boolean{Value: v}
	case nil:
		return &object.Null{}
	case []any:
		elements := make([]object.Object, len(v))
		for i, el := range v {
			elements[i] = ToObject(el)
		}
		return &object.Array{Elements: elements}
	case map[string]any:
		pairs := make(map[object.HashKey]object.HashPair)
		for k, val := range v {
			key := &object.String{Value: k}
			hashKey := key.HashKey()
			pairs[hashKey] = object.HashPair{
				Key:   key,
				Value: ToObject(val),
			}
		}
		return &object.Hash{Pairs: pairs}
	default:
		return &object.Null{}
	}
}

func stringify(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	data := FromObject(args[0])
	bytes, errGo := json.Marshal(data)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
	}

	return &object.String{Value: string(bytes)}
}

func FromObject(obj object.Object) any {
	switch o := obj.(type) {
	case *object.String:
		return o.Value
	case *object.Integer:
		return o.Value
	case *object.Float:
		return o.Value
	case *object.Boolean:
		return o.Value
	case *object.Null:
		return nil
	case *object.Array:
		elements := make([]any, len(o.Elements))
		for i, el := range o.Elements {
			elements[i] = FromObject(el)
		}
		return elements
	case *object.Hash:
		m := make(map[string]any)
		for _, pair := range o.Pairs {
			var key string
			if s, ok := pair.Key.(*object.String); ok {
				key = s.Value
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
