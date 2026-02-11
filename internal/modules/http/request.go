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

func newRequest(r *http.Request) object.Value {
	var bodyCache *string
	var bodyJSONCache object.Value

	// defer r.Body.Close()
	return modbuilder.NewHashBuilder().
		AddFunction("path_param", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}

			val := r.PathValue(args[0].AsString().Value)
			return object.NewString(val)
		}).
		AddFunction("query_param", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}

			param := r.URL.Query().Get(args[0].AsString().Value)
			if param == "" {
				return object.NewNull()
			}

			return object.NewString(param)
		}).
		AddFunction("body", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err.IsError() {
				return err
			}

			if bodyCache != nil {
				return object.NewString(*bodyCache)
			}

			byteData, errGo := io.ReadAll(r.Body)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
			}
			defer r.Body.Close()

			bodyString := string(byteData)
			bodyCache = &bodyString
			return object.NewString(bodyString)
		}).
		AddFunction("json", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err.IsError() {
				return err
			}

			if !bodyJSONCache.IsNull() {
				return bodyJSONCache
			}

			if bodyCache == nil {
				byteData, errGo := io.ReadAll(r.Body)
				if errGo != nil {
					return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
				}
				defer r.Body.Close()
				bodyString := string(byteData)
				bodyCache = &bodyString
			}

			var data any
			errGo := gojson.Unmarshal([]byte(*bodyCache), &data)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
			}
			bodyJSONCache = json.ToObject(data)

			return bodyJSONCache
		}).
		AddFunction("method", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err.IsError() {
				return err
			}

			return object.NewString(r.Method)
		}).
		AddFunction("url", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err.IsError() {
				return err
			}

			return object.NewString(r.URL.String())
		}).
		AddFunction("header", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}

			header := r.Header.Get(args[0].AsString().Value)
			if header == "" {
				return object.NewNull()
			}

			return object.NewString(header)
		}).
		AddFunction("headers", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err.IsError() {
				return err
			}

			pairs := make(map[object.HashKey]object.HashPair)
			for key, values := range r.Header {
				keyObj := object.NewString(key)
				hashKey := keyObj.HashKey()

				valueObj := object.NewString(strings.Join(values, ", "))
				pairs[hashKey] = object.HashPair{Key: keyObj, Value: valueObj}
			}

			return object.NewHash(pairs)
		}).
		AddFunction("cookie", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}
			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}

			cookie, errGo := r.Cookie(args[0].AsString().Value)
			if errGo != nil {
				return object.NewNull()
			}

			return object.NewString(cookie.Value)
		}).
		AddFunction("cookies", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 0, args)
			if err.IsError() {
				return err
			}

			pairs := make(map[object.HashKey]object.HashPair)
			for _, cookie := range r.Cookies() {
				keyObj := object.NewString(cookie.Name)
				hashKey := keyObj.HashKey()
				valueObj := object.NewString(cookie.Value)
				pairs[hashKey] = object.NewHashPair(keyObj, valueObj)
			}

			return object.NewHash(pairs)
		}).
		Build()
}
