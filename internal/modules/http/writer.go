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
				err = errors.ExpectType(2, 0, args[2], object.HashKind)
				if err.IsError() {
					return err
				}

				options := args[2].AsHash()

				// MaxAge
				if pair, ok := options.Pairs[(object.NewString("max_age")).HashKey()]; ok && pair.Value.IsInt() {
					intVal := pair.Value.AsInt()
					cookie.MaxAge = int(intVal)
				}

				// Domain
				if pair, ok := options.Pairs[(object.NewString("domain")).HashKey()]; ok && pair.Value.IsString() {
					cookie.Domain = pair.Value.AsString().Value
				}

				// Path
				if pair, ok := options.Pairs[(object.NewString("path")).HashKey()]; ok && pair.Value.IsString() {
					cookie.Path = pair.Value.AsString().Value
				}

				// Secure
				if pair, ok := options.Pairs[(object.NewString("secure")).HashKey()]; ok && pair.Value.IsBool() {
					cookie.Secure = pair.Value.AsBool()
				}

				// HttpOnly
				if pair, ok := options.Pairs[(object.NewString("http_only")).HashKey()]; ok && pair.Value.IsBool() {
					cookie.HttpOnly = pair.Value.AsBool()
				}

				// SameSite
				if pair, ok := options.Pairs[(object.NewString("same_site")).HashKey()]; ok && pair.Value.IsString() {
					switch pair.Value.AsString().Value {
					case "strict":
						cookie.SameSite = http.SameSiteStrictMode
					case "lax":
						cookie.SameSite = http.SameSiteLaxMode
					case "none":
						cookie.SameSite = http.SameSiteNoneMode
					}
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
