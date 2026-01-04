package http

import (
	"net/http"
	"strconv"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

var Module = modbuilder.NewModuleBuilder().
	AddFunction("create_server", createServer).
	AddValue("status", statusCodes).
	AddValue("methods", methods).
	Build()

func createServer(args ...object.Object) object.Object {
	server := modbuilder.NewHashBuilder().
		AddFunction("get", getRoute).
		AddFunction("listen", listen).
		AddValue("status", statusCodes).
		Build()

	return server
}

func getRoute(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(1, 0, args[1], object.FunctionObj)
	if err != nil {
		return err
	}

	route := args[0].(*object.String).Value
	callback := args[1].(*object.Function)

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		args := []object.Object{
			newWriter(w),
			newRequest(r),
		}
		
		if object.ApplyFunction != nil {
			object.ApplyFunction(callback, args)
		}
	})
	
	return &object.Null{}
}

func listen(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.IntegerObj)
	if err != nil {
		return err
	}

	port := args[0].(*object.Integer).Value

	errGo := http.ListenAndServe(":" + strconv.FormatInt(port, 10), nil)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.Null{}
}

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

func newRequest(r *http.Request) object.Object {
	return &object.Null{}	
}

var statusCodes = modbuilder.NewHashBuilder().
		AddInteger("continue", http.StatusContinue).
		AddInteger("switching_protocols", http.StatusSwitchingProtocols).
		AddInteger("processing", http.StatusProcessing).
		AddInteger("early_hints", http.StatusEarlyHints).
		AddInteger("ok", http.StatusOK).
		AddInteger("created", http.StatusCreated).
		AddInteger("accepted", http.StatusAccepted).
		AddInteger("non_authoritative_info", http.StatusNonAuthoritativeInfo).
		AddInteger("no_content", http.StatusNoContent).
		AddInteger("reset_content", http.StatusResetContent).
		AddInteger("partial_content", http.StatusPartialContent).
		AddInteger("multi_status", http.StatusMultiStatus).
		AddInteger("already_reported", http.StatusAlreadyReported).
		AddInteger("im_used", http.StatusIMUsed).
		AddInteger("multiple_choices", http.StatusMultipleChoices).
		AddInteger("moved_permanently", http.StatusMovedPermanently).
		AddInteger("found", http.StatusFound).
		AddInteger("see_other", http.StatusSeeOther).
		AddInteger("not_modified", http.StatusNotModified).
		AddInteger("use_proxy", http.StatusUseProxy).
		AddInteger("temporary_redirect", http.StatusTemporaryRedirect).
		AddInteger("permanent_redirect", http.StatusPermanentRedirect).
		AddInteger("bad_request", http.StatusBadRequest).
		AddInteger("unauthorized", http.StatusUnauthorized).
		AddInteger("payment_required", http.StatusPaymentRequired).
		AddInteger("forbidden", http.StatusForbidden).
		AddInteger("not_found", http.StatusNotFound).
		AddInteger("method_not_allowed", http.StatusMethodNotAllowed).
		AddInteger("not_acceptable", http.StatusNotAcceptable).
		AddInteger("proxy_auth_required", http.StatusProxyAuthRequired).
		AddInteger("request_timeout", http.StatusRequestTimeout).
		AddInteger("conflict", http.StatusConflict).
		AddInteger("gone", http.StatusGone).
		AddInteger("length_required", http.StatusLengthRequired).
		AddInteger("precondition_failed", http.StatusPreconditionFailed).
		AddInteger("payload_too_large", http.StatusRequestEntityTooLarge).
		AddInteger("uri_too_long", http.StatusRequestURITooLong).
		AddInteger("unsupported_media_type", http.StatusUnsupportedMediaType).
		AddInteger("range_not_satisfiable", http.StatusRequestedRangeNotSatisfiable).
		AddInteger("expectation_failed", http.StatusExpectationFailed).
		AddInteger("im_a_teapot", http.StatusTeapot).
		AddInteger("misdirected_request", http.StatusMisdirectedRequest).
		AddInteger("unprocessable_entity", http.StatusUnprocessableEntity).
		AddInteger("locked", http.StatusLocked).
		AddInteger("failed_dependency", http.StatusFailedDependency).
		AddInteger("too_early", http.StatusTooEarly).
		AddInteger("upgrade_required", http.StatusUpgradeRequired).
		AddInteger("precondition_required", http.StatusPreconditionRequired).
		AddInteger("too_many_requests", http.StatusTooManyRequests).
		AddInteger("request_header_fields_too_large", http.StatusRequestHeaderFieldsTooLarge).
		AddInteger("unavailable_for_legal_reasons", http.StatusUnavailableForLegalReasons).
		AddInteger("internal_server_error", http.StatusInternalServerError).
		AddInteger("not_implemented", http.StatusNotImplemented).
		AddInteger("bad_gateway", http.StatusBadGateway).
		AddInteger("service_unavailable", http.StatusServiceUnavailable).
		AddInteger("gateway_timeout", http.StatusGatewayTimeout).
		AddInteger("http_version_not_supported", http.StatusHTTPVersionNotSupported).
		AddInteger("variant_also_negotiates", http.StatusVariantAlsoNegotiates).
		AddInteger("insufficient_storage", http.StatusInsufficientStorage).
		AddInteger("loop_detected", http.StatusLoopDetected).
		AddInteger("not_extended", http.StatusNotExtended).
		AddInteger("network_authentication_required", http.StatusNetworkAuthenticationRequired).
		Build()

var methods = modbuilder.NewHashBuilder().
		AddString("get", http.MethodGet).
		AddString("post", http.MethodPost).
		AddString("put", http.MethodPut).
		AddString("patch", http.MethodPatch).
		AddString("delete", http.MethodDelete).
		AddString("head", http.MethodHead).
		AddString("options", http.MethodOptions).
		AddString("connect", http.MethodConnect).
		AddString("trace", http.MethodTrace).
		Build()