package time

import (
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
	"time"
)

func New() *object.Hash {
	return modbuilder.NewModuleBuilder().
		AddFunction("now", now).
		AddFunction("now_ms", nowMs).
		AddFunction("now_ns", nowNs).
		AddFunction("sleep", sleep).
		AddFunction("unix", unix).
		AddFunction("unix_ms", unixMs).
		AddFunction("unix_ns", unixNs).
		AddFunction("format", formatTime).
		AddFunction("parse", parseTime).
		AddFloat("millisecond", 1.0/1000).
		AddInteger("second", 1).
		AddInteger("minute", 60).
		AddInteger("hour", 60*60).
		AddInteger("day", 60*60*24).
		AddInteger("week", 60*60*24*7).
		Build()
}

// now returns current Unix timestamp as integer
func now(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 0, args)
	if err != nil {
		return err
	}

	return &object.Integer{Value: time.Now().Unix()}
}

// nowMs returns current Unix timestamp in milliseconds
func nowMs(args ...object.Object) object.Object {
	if err := errors.ExpectNumberOfArguments(0, 0, 0, args); err != nil {
		return err
	}
	return &object.Integer{Value: time.Now().UnixMilli()}
}

// nowNs returns current Unix timestamp in nanoseconds
func nowNs(args ...object.Object) object.Object {
	if err := errors.ExpectNumberOfArguments(0, 0, 0, args); err != nil {
		return err
	}
	return &object.Integer{Value: time.Now().UnixNano()}
}

func sleep(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectOneOfTypes(0, 0, args[0], object.IntegerObj, object.FloatObj)
	if err != nil {
		return err
	}

	var d time.Duration
	switch v := args[0].(type) {
	case *object.Integer:
		d = time.Duration(v.Value) * time.Second
	case *object.Float:
		d = time.Duration(v.Value * float64(time.Second))
	}

	time.Sleep(d)
	return &object.Null{}
}

// unix converts a Unix timestamp to a time object
func unix(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.IntegerObj)
	if err != nil {
		return err
	}

	timestamp := args[0].(*object.Integer).Value
	t := time.Unix(timestamp, 0)

	return timeToHash(t)
}

func unixMs(args ...object.Object) object.Object {
	if err := errors.ExpectNumberOfArguments(0, 0, 1, args); err != nil {
		return err
	}

	if err := errors.ExpectType(0, 0, args[0], object.IntegerObj); err != nil {
		return err
	}

	ms := args[0].(*object.Integer).Value
	t := time.UnixMilli(ms)
	return timeToHash(t)
}

func unixNs(args ...object.Object) object.Object {
	if err := errors.ExpectNumberOfArguments(0, 0, 1, args); err != nil {
		return err
	}

	if err := errors.ExpectType(0, 0, args[0], object.IntegerObj); err != nil {
		return err
	}

	ns := args[0].(*object.Integer).Value
	t := time.Unix(0, ns)
	return timeToHash(t)
}

func formatTime(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	formatStr := args[1].(*object.String).Value
	var t time.Time

	switch v := args[0].(type) {
	case *object.Integer:
		t = time.Unix(v.Value, 0)
	case *object.Hash:
		var errObj *object.Error
		t, errObj = hashToTime(v)
		if errObj != nil {
			return errObj
		}
	default:
		return errors.NewTypeError(0, 0, "expected Integer or Hash, got %s", args[0].Type().String())
	}

	// Convert Go format to common format patterns
	goFormat := convertFormat(formatStr)
	return &object.String{Value: t.Format(goFormat)}
}

func parseTime(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[1], object.StringObj)
	if err != nil {
		return err
	}

	timeStr := args[0].(*object.String).Value
	formatStr := args[1].(*object.String).Value

	goFormat := convertFormat(formatStr)
	t, errGo := time.Parse(goFormat, timeStr)
	if errGo != nil {
		return errors.NewRuntimeError(0, 0, "failed to parse time: %s", errGo.Error())
	}

	return timeToHash(t)
}

// convertFormat converts common format patterns to Go's layout
func convertFormat(format string) string {
	conversions := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"mm":   "04",
		"ss":   "05",
		"SSS":  "000",
	}

	result := format
	for k, v := range conversions {
		result = replaceAll(result, k, v)
	}
	return result
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

func timeToHash(t time.Time) *object.Hash {
	return modbuilder.NewHashBuilder().
		AddInteger("unix", t.Unix()).
		AddInteger("unix_ms", t.UnixMilli()).
		AddInteger("unix_ns", t.UnixNano()).
		AddInteger("year", int64(t.Year())).
		AddInteger("month", int64(t.Month())).
		AddInteger("day", int64(t.Day())).
		AddInteger("hour", int64(t.Hour())).
		AddInteger("minute", int64(t.Minute())).
		AddInteger("second", int64(t.Second())).
		AddInteger("millisecond", int64(t.Nanosecond()/1e6)).
		AddInteger("nanosecond", int64(t.Nanosecond())).
		AddInteger("weekday", int64(t.Weekday())).
		Build()
}

func hashToTime(h *object.Hash) (time.Time, *object.Error) {
	// Check for unix_ns first for maximum precision
	if pair, ok := h.Pairs[(&object.String{Value: "unix_ns"}).HashKey()]; ok {
		if ns, ok := pair.Value.(*object.Integer); ok {
			return time.Unix(0, ns.Value), nil
		}
	}

	// Fallback to standard unix seconds
	unixKey := &object.String{Value: "unix"}
	if pair, ok := h.Pairs[unixKey.HashKey()]; ok {
		if unixInt, ok := pair.Value.(*object.Integer); ok {
			return time.Unix(unixInt.Value, 0), nil
		}
	}

	return time.Time{}, errors.NewRuntimeError(0, 0, "time hash missing 'unix' or 'unix_ns' field")
}
