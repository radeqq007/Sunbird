package fs

import (
	"os"
	"path/filepath"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() object.Value {
	return modbuilder.NewModuleBuilder().
		AddFunction("read", readFile).
		AddFunction("write", writeFile).
		AddFunction("append", appendFile).
		AddFunction("remove", removeFile).
		AddFunction("exists", exists).
		AddFunction("is_dir", isDir).
		AddFunction("list_dir", listDir).
		AddFunction("create_dir", createDir).
		AddFunction("rename", rename).
		AddFunction("copy", copyFile).
		Build()
}

func getFullPath(requestedPath string) string {
	var fullPath string
	if filepath.IsAbs(requestedPath) {
		fullPath = requestedPath
	} else {
		mainFileDir := ""
		if len(os.Args) > 1 {
			mainFileDir = filepath.Dir(os.Args[1])
		}
		fullPath = filepath.Join(mainFileDir, requestedPath)
	}
	return fullPath
}

func readFile(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	data, errGo := os.ReadFile(getFullPath(args[0].AsString().Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewString(string(data))
}

func writeFile(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 1, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	errGo := os.WriteFile(getFullPath(args[0].AsString().Value), []byte(args[1].AsString().Value), 0o644)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}
	
	return object.NewNull()
}

func appendFile(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 1, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	file, errGo := os.OpenFile(getFullPath(args[0].AsString().Value), os.O_APPEND|os.O_WRONLY, 0o644)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}
	defer file.Close()

	_, errGo = file.WriteString(args[1].AsString().Value)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewNull()
}

func removeFile(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	errGo := os.Remove(getFullPath(args[0].AsString().Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewNull()
}

func exists(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	_, errGo := os.Stat(getFullPath(args[0].AsString().Value))
	if errGo != nil {
		if os.IsNotExist(errGo) {
			return object.NewBool(false)
		}
		// Other errors (permissions, etc.)
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewBool(true)
}

func isDir(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	info, errGo := os.Stat(getFullPath(args[0].AsString().Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewBool(info.IsDir())
}

func listDir(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	entries, errGo := os.ReadDir(getFullPath(args[0].AsString().Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	var result []object.Value
	for _, entry := range entries {
		result = append(result, object.NewString(entry.Name()))
	}

	return object.NewArray(result)
}

func createDir(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	errGo := os.Mkdir(getFullPath(args[0].AsString().Value), 0o755)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}
	return object.NewNull()
}

func rename(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 1, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	errGo := os.Rename(getFullPath(args[0].AsString().Value), getFullPath(args[1].AsString().Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewNull()
}

func copyFile(args ...object.Value) object.Value {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringKind)
	if err.IsError() {
		return err
	}

	err = errors.ExpectType(0, 1, args[1], object.StringKind)
	if err.IsError() {
		return err
	}

	str := args[0].AsString().Value

	data, errGo := os.ReadFile(getFullPath(str))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	errGo = os.WriteFile(getFullPath(args[1].AsString().Value), data, 0o644)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, "%s", errGo.Error())
	}

	return object.NewNull()
}
