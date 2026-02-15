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

type request struct {
	r             *http.Request
	bodyCache     *string
	bodyJSONCache object.Value
}

func newRequest(r *http.Request) object.Value {
	req := &request{r: r}

	return modbuilder.NewHashBuilder().
		AddFunction("path_param", req.pathParam).
		AddFunction("query_param", req.queryParam).
		AddFunction("body", req.body).
		AddFunction("json", req.json).
		AddFunction("method", req.method).
		AddFunction("url", req.url).
		AddFunction("header", req.header).
		AddFunction("headers", req.headers).
		AddFunction("cookie", req.cookie).
		AddFunction("cookies", req.cookies).
		Build()
}

func (req *request) pathParam(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	val := req.r.PathValue(args[0].AsString().Value)
	return object.NewString(val)
}

func (req *request) queryParam(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	param := req.r.URL.Query().Get(args[0].AsString().Value)
	if param == "" {
		return object.NewNull()
	}

	return object.NewString(param)
}

func (req *request) body(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	if req.bodyCache != nil {
		return object.NewString(*req.bodyCache)
	}

	byteData, errGo := io.ReadAll(req.r.Body)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
	}
	defer req.r.Body.Close()

	bodyString := string(byteData)
	req.bodyCache = &bodyString
	return object.NewString(bodyString)
}

func (req *request) json(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	if !req.bodyJSONCache.IsNull() {
		return req.bodyJSONCache
	}

	if req.bodyCache == nil {
		byteData, errGo := io.ReadAll(req.r.Body)
		if errGo != nil {
			return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
		}
		defer req.r.Body.Close()
		bodyString := string(byteData)
		req.bodyCache = &bodyString
	}

	var data any
	errGo := gojson.Unmarshal([]byte(*req.bodyCache), &data)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
	}
	req.bodyJSONCache = json.ToObject(data)

	return req.bodyJSONCache
}

func (req *request) method(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	return object.NewString(req.r.Method)
}

func (req *request) url(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	return object.NewString(req.r.URL.String())
}

func (req *request) header(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	header := req.r.Header.Get(args[0].AsString().Value)
	if header == "" {
		return object.NewNull()
	}

	return object.NewString(header)
}

func (req *request) headers(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	pairs := make(map[object.HashKey]object.HashPair)
	for key, values := range req.r.Header {
		keyObj := object.NewString(key)
		hashKey := keyObj.HashKey()

		valueObj := object.NewString(strings.Join(values, ", "))
		pairs[hashKey] = object.HashPair{Key: keyObj, Value: valueObj}
	}

	return object.NewHash(pairs)
}

func (req *request) cookie(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}
	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	cookie, errGo := req.r.Cookie(args[0].AsString().Value)
	if errGo != nil {
		return object.NewNull()
	}

	return object.NewString(cookie.Value)
}

func (req *request) cookies(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err.IsError() {
		return err
	}

	pairs := make(map[object.HashKey]object.HashPair)
	for _, cookie := range req.r.Cookies() {
		keyObj := object.NewString(cookie.Name)
		hashKey := keyObj.HashKey()
		valueObj := object.NewString(cookie.Value)
		pairs[hashKey] = object.NewHashPair(keyObj, valueObj)
	}

	return object.NewHash(pairs)
}
