package http

import (
	gojson "encoding/json"
	"net/http"
	"sunbird/internal/errors"
	"sunbird/internal/modules/json"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func newWriter(w http.ResponseWriter) object.Value {
	return modbuilder.NewHashBuilder().
		AddFunction("send", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}

			_, errGo := w.Write([]byte(args[0].AsString().Value))
			if errGo != nil {
				return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
			}

			return object.NewNull()
		}).
		AddFunction("json", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.HashKind)
			if err.IsError() {
				return err
			}

			data := json.FromObject(args[0])
			bytes, errGo := gojson.Marshal(data)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
			}

			w.Header().Set("Content-Type", "application/json")

			_, errGo = w.Write(bytes)
			if errGo != nil {
				return errors.NewRuntimeError(0, 0, "%s", errGo.Error())
			}

			return object.NewNull()
		}).
		AddValue("header", modbuilder.NewHashBuilder().
			AddFunction("set", func(args ...object.Value) object.Value {
				err := errors.ExpectNumberOfArguments(0, 0, 2, args)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringKind)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(1, 0, args[1], object.StringKind)
				if err.IsError() {
					return err
				}

				w.Header().Set(args[0].AsString().Value, args[1].AsString().Value)
				return object.NewNull()
			}).
			AddFunction("add", func(args ...object.Value) object.Value {
				err := errors.ExpectNumberOfArguments(0, 0, 2, args)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringKind)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(1, 0, args[1], object.StringKind)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(1, 0, args[1], object.StringKind)
				if err.IsError() {
					return err
				}

				w.Header().Add(args[0].AsString().Value, args[1].AsString().Value)

				return object.NewNull()
			}).
			AddFunction("del", func(args ...object.Value) object.Value {
				err := errors.ExpectNumberOfArguments(0, 0, 1, args)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringKind)
				if err.IsError() {
					return err
				}

				w.Header().Del(args[0].AsString().Value)

				return object.NewNull()
			}).
			AddFunction("get", func(args ...object.Value) object.Value {
				err := errors.ExpectNumberOfArguments(0, 0, 1, args)
				if err.IsError() {
					return err
				}

				err = errors.ExpectType(0, 0, args[0], object.StringKind)
				if err.IsError() {
					return err
				}

				return object.NewString(w.Header().Get(args[0].AsString().Value))
			}).
			Build(),
		).
		AddFunction("status", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.IntKind)
			if err.IsError() {
				return err
			}

			w.WriteHeader(int(args[0].AsInt()))

			return object.NewNull()
		}).
		AddValue("cookie", cookieHash(w)).
		Build()
}

func cookieHash(w http.ResponseWriter) object.Value {
	h := modbuilder.NewHashBuilder().
		AddFunction("set", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 2, args)
			if err.IsError() {
				err = errors.ExpectNumberOfArguments(0, 0, 3, args)
				if err.IsError() {
					return err
				}
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}
			err = errors.ExpectType(1, 0, args[1], object.StringKind)
			if err.IsError() {
				return err
			}

			cookie := &http.Cookie{
				Name:  args[0].AsString().Value,
				Value: args[1].AsString().Value,
				Path:  "/",
			}

			// Parse options if provided
			if len(args) == 3 {
				if err := applyCookieOptions(cookie, args[2]); err.IsError() {
					return err
				}
			}

			http.SetCookie(w, cookie)
			return object.NewNull()
		}).
		AddFunction("delete", func(args ...object.Value) object.Value {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err.IsError() {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringKind)
			if err.IsError() {
				return err
			}

			cookie := &http.Cookie{
				Name:   args[0].AsString().Value,
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}

			http.SetCookie(w, cookie)

			return object.NewNull()
		}).
		Build()

	return h
}

func applyCookieOptions(cookie *http.Cookie, optionsObj object.Value) object.Value {
	if err := errors.ExpectType(2, 0, optionsObj, object.HashKind); err.IsError() {
		return err
	}

	options := optionsObj.AsHash()

	getVal := func(key string) (object.Value, bool) {
		pair, ok := options.Pairs[object.NewString(key).HashKey()]
		if !ok {
			return object.NewNull(), false
		}
		return pair.Value, true
	}

	if val, ok := getVal("max_age"); ok && val.IsInt() {
		cookie.MaxAge = int(val.AsInt())
	}
	if val, ok := getVal("domain"); ok && val.IsString() {
		cookie.Domain = val.AsString().Value
	}
	if val, ok := getVal("path"); ok && val.IsString() {
		cookie.Path = val.AsString().Value
	}
	if val, ok := getVal("secure"); ok && val.IsBool() {
		cookie.Secure = val.AsBool()
	}
	if val, ok := getVal("http_only"); ok && val.IsBool() {
		cookie.HttpOnly = val.AsBool()
	}
	if val, ok := getVal("same_site"); ok && val.IsString() {
		switch val.AsString().Value {
		case "strict":
			cookie.SameSite = http.SameSiteStrictMode
		case "lax":
			cookie.SameSite = http.SameSiteLaxMode
		case "none":
			cookie.SameSite = http.SameSiteNoneMode
		}
	}

	return object.NewNull()
}
