package http

import (
	gojson "encoding/json"
	"io"
	"net/http"
	"strings"
	"sunbird/internal/errors"
	"sunbird/internal/modules/json"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func newRequest(r *http.Request) object.Object {
	var bodyCache *string
	var bodyJSONCache object.Object

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

			if bodyJSONCache != nil {
				return bodyJSONCache
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
			bodyJSONCache = json.ToObject(data)

			return bodyJSONCache
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
		AddFunction("header", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			header := r.Header.Get(args[0].(*object.String).Value)

			if header == "" {
				return &object.Null{}
			}

			return &object.String{Value: header}
		}).
		AddFunction("headers", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			pairs := make(map[object.HashKey]object.HashPair)
			for key, values := range r.Header {
				keyObj := &object.String{Value: key}
				hashKey := keyObj.HashKey()

				valueObj := &object.String{Value: strings.Join(values, ", ")}
				pairs[hashKey] = object.HashPair{Key: keyObj, Value: valueObj}
			}

			return &object.Hash{
				Pairs: pairs,
			}
		}).
		AddFunction("cookie", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}
			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			cookie, errGo := r.Cookie(args[0].(*object.String).Value)
			if errGo != nil {
				return &object.Null{}
			}

			return &object.String{Value: cookie.Value}
		}).
		AddFunction("cookies", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err != nil {
				return err
			}

			pairs := make(map[object.HashKey]object.HashPair)
			for _, cookie := range r.Cookies() {
				keyObj := &object.String{Value: cookie.Name}
				hashKey := keyObj.HashKey()
				valueObj := &object.String{Value: cookie.Value}
				pairs[hashKey] = object.HashPair{Key: keyObj, Value: valueObj}
			}

			return &object.Hash{Pairs: pairs}
		}).
		Build()
}
