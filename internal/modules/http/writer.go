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
		AddValue("cookie", cookieHash(w)).
		Build()
}

func cookieHash(w http.ResponseWriter) object.Object {
	h := modbuilder.NewHashBuilder().
		AddFunction("set", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 2, args)
			if err != nil {
				err = errors.ExpectNumberOfArguments(0, 0, 3, args)
				if err != nil {
					return err
				}
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}
			err = errors.ExpectType(1, 0, args[1], object.StringObj)
			if err != nil {
				return err
			}

			cookie := &http.Cookie{
				Name:  args[0].(*object.String).Value,
				Value: args[1].(*object.String).Value,
				Path:  "/",
			}

			// Parse options if provided
			if len(args) == 3 {
				err = errors.ExpectType(2, 0, args[2], object.HashObj)
				if err != nil {
					return err
				}

				options := args[2].(*object.Hash)

				// MaxAge
				if pair, ok := options.Pairs[(&object.String{Value: "max_age"}).HashKey()]; ok {
					if intVal, ok := pair.Value.(*object.Integer); ok {
						cookie.MaxAge = int(intVal.Value)
					}
				}

				// Domain
				if pair, ok := options.Pairs[(&object.String{Value: "domain"}).HashKey()]; ok {
					if strVal, ok := pair.Value.(*object.String); ok {
						cookie.Domain = strVal.Value
					}
				}

				// Path
				if pair, ok := options.Pairs[(&object.String{Value: "path"}).HashKey()]; ok {
					if strVal, ok := pair.Value.(*object.String); ok {
						cookie.Path = strVal.Value
					}
				}

				// Secure
				if pair, ok := options.Pairs[(&object.String{Value: "secure"}).HashKey()]; ok {
					if boolVal, ok := pair.Value.(*object.Boolean); ok {
						cookie.Secure = boolVal.Value
					}
				}

				// HttpOnly
				if pair, ok := options.Pairs[(&object.String{Value: "http_only"}).HashKey()]; ok {
					if boolVal, ok := pair.Value.(*object.Boolean); ok {
						cookie.HttpOnly = boolVal.Value
					}
				}

				// SameSite
				if pair, ok := options.Pairs[(&object.String{Value: "same_site"}).HashKey()]; ok {
					if strVal, ok := pair.Value.(*object.String); ok {
						switch strVal.Value {
						case "strict":
							cookie.SameSite = http.SameSiteStrictMode
						case "lax":
							cookie.SameSite = http.SameSiteLaxMode
						case "none":
							cookie.SameSite = http.SameSiteNoneMode
						}
					}
				}
			}

			http.SetCookie(w, cookie)
			return &object.Null{}
		}).
		AddFunction("delete", func(args ...object.Object) object.Object {
			err := errors.ExpectNumberOfArguments(0, 0, 1, args)
			if err != nil {
				return err
			}

			err = errors.ExpectType(0, 0, args[0], object.StringObj)
			if err != nil {
				return err
			}

			cookie := &http.Cookie{
				Name:   args[0].(*object.String).Value,
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}

			http.SetCookie(w, cookie)

			return &object.Null{}
		}).
		Build()

	return h
}
