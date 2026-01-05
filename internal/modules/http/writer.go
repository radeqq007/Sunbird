package http

import (
	gojson "encoding/json"
	"net/http"
	"sunbird/internal/errors"
	"sunbird/internal/modules/json"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func newWriter(w http.ResponseWriter) object.Object {
	return modbuilder.NewHashBuilder().
		AddFunction("send", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			_, errGo := w.Write([]byte(args[0].(*object.String).Value))
			if errGo != nil {
				return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
			}

			return &object.Null{}
		}).
		AddFunction("json", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.HashObj)
			if err != nil {
				return err
			}

			data := json.FromObject(args[0])
			bytes, errGo := gojson.Marshal(data)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, errGo.Error())
			}

			w.Header().Set("Content-Type", "application/json")

			_, errGo = w.Write(bytes)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, errGo.Error())
			}

			return &object.Null{}
		}).
		AddValue("header", modbuilder.NewHashBuilder().
			AddFunction("set", func(args ...object.Object) object.Object {
				err := errors.ExpectNumberOfArguments(0, 0, 2, args)
				if err != nil {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringObj)
				if err != nil {
					return err
				}

				err = errors.ExpectType(1, 0, args[1], object.StringObj)
				if err != nil {
					return err
				}

				w.Header().Set(args[0].(*object.String).Value, args[1].(*object.String).Value)
				return &object.Null{}
			}).
			AddFunction("add", func(args ...object.Object) object.Object {
				err := errors.ExpectNumberOfArguments(0, 0, 2, args)
				if err != nil {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringObj)
				if err != nil {
					return err
				}

				err = errors.ExpectType(1, 0, args[1], object.StringObj)
				if err != nil {
					return err
				}

				w.Header().Add(args[0].(*object.String).Value, args[1].(*object.String).Value)

				return &object.Null{}
			}).
			AddFunction("del", func(args ...object.Object) object.Object {
				err := errors.ExpectNumberOfArguments(0, 0, 1, args)
				if err != nil {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringObj)
				if err != nil {
					return err
				}

				w.Header().Del(args[0].(*object.String).Value)

				return &object.Null{}
			}).
			AddFunction("get", func(args ...object.Object) object.Object {
				err := errors.ExpectNumberOfArguments(0, 0, 1, args)
				if err != nil {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringObj)
				if err != nil {
					return err
				}

				return &object.String{Value: w.Header().Get(args[0].(*object.String).Value)}
			}).
			Build(),
		).
		AddFunction("status", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.IntegerObj)
			if err != nil {
				return err
			}

			w.WriteHeader(int(args[0].(*object.Integer).Value))

			return &object.Null{}
		}).
		Build()
}
