package http

import (
	gojson "encoding/json"
	"io"
	"net/http"
	"sunbird/internal/errors"
	"sunbird/internal/modules/json"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func newRequest(r *http.Request) object.Object {
	var bodyCache *string
	var bodyJsonCache object.Object

	// defer r.Body.Close()
	return modbuilder.NewHashBuilder().
		AddFunction("path_param", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			val := r.PathValue(args[0].(*object.String).Value)
			return &object.String{Value: val}
		}).
		AddFunction("query_param", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			param := r.URL.Query().Get(args[0].(*object.String).Value)
			if param == "" {
				return &object.Null{}
			}

			return &object.String{Value: param}
		}).
		AddFunction("body", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			if bodyCache != nil {
				return &object.String{Value: *bodyCache}
			}

			byteData, errGo := io.ReadAll(r.Body)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, errGo.Error())
			}
			defer r.Body.Close()

			bodyString := string(byteData)
			bodyCache = &bodyString
			return &object.String{Value: bodyString}
		}).
		AddFunction("json", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			if bodyJsonCache != nil {
				return bodyJsonCache
			}

			if bodyCache == nil {
				byteData, errGo := io.ReadAll(r.Body)
				if errGo != nil {
					return errors.NewRuntimeError(0, 0, errGo.Error())
				}
				defer r.Body.Close()
				bodyString := string(byteData)
				bodyCache = &bodyString
			}

			var data any
			errgo := gojson.Unmarshal([]byte(*bodyCache), &data)
			if errgo != nil {
				return errors.NewRuntimeError(0, 0, errgo.Error())
			}
			bodyJsonCache = json.ToObject(data)

			return bodyJsonCache
		}).
		AddFunction("method", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			return &object.String{Value: r.Method}
		}).
		AddFunction("url", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			return &object.String{Value: r.URL.String()}
		}).
		Build()
}
