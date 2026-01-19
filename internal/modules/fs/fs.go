package fs

import (
	"os"
	"path/filepath"
	"sunbird/internal/errors"
	"sunbird/internal/modules/modbuilder"
	"sunbird/internal/object"
)

func New() *object.Hash {
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
		AddFunction("copy", copy).
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


func readFile(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	data, errGo := os.ReadFile(getFullPath(args[0].(*object.String).Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.String{Value: string(data)}
}

func writeFile(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}	

	err = errors.ExpectType(0, 1, args[1], object.StringObj)
	if err != nil {
		return err
	}

	os.WriteFile(getFullPath(args[0].(*object.String).Value), []byte(args[1].(*object.String).Value), 0644)

	return &object.Null{}
}

func appendFile(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}	

	err = errors.ExpectType(0, 1, args[1], object.StringObj)
	if err != nil {
		return err
	}

	// os.WriteFile(getFullPath(args[0].(*object.String).Value), []byte(args[1].(*object.String).Value), 0644)
	file, errGo := os.OpenFile(getFullPath(args[0].(*object.String).Value), os.O_APPEND|os.O_WRONLY, 0644)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}
	defer file.Close()

	_, errGo = file.WriteString(args[1].(*object.String).Value)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.Null{}
}

func removeFile(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	errGo := os.Remove(getFullPath(args[0].(*object.String).Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.Null{}
}

func exists(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	_, errGo := os.Stat(getFullPath(args[0].(*object.String).Value))
	if errGo != nil {
		if os.IsNotExist(errGo) {
			return &object.Boolean{Value: false}
		}
		// Other errors (permissions, etc.)
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.Boolean{Value: true}
}

func isDir(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	info, errGo := os.Stat(getFullPath(args[0].(*object.String).Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.Boolean{Value: info.IsDir()}
}

func listDir(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	entries, errGo := os.ReadDir(getFullPath(args[0].(*object.String).Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	var result []object.Object
	for _, entry := range entries {
		result = append(result, &object.String{Value: entry.Name()})
	}

	return &object.Array{Elements: result}
}

func createDir(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 1, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	os.Mkdir(getFullPath(args[0].(*object.String).Value), 0755)

	return &object.Null{}
}

func rename(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 1, args[1], object.StringObj)
	if err != nil {
		return err
	}

	os.Rename(getFullPath(args[0].(*object.String).Value), getFullPath(args[1].(*object.String).Value))

	return &object.Null{}
}

func copy(args ...object.Object) object.Object {
	err := errors.ExpectNumberOfArguments(0, 0, 2, args)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 0, args[0], object.StringObj)
	if err != nil {
		return err
	}

	err = errors.ExpectType(0, 1, args[1], object.StringObj)
	if err != nil {
		return err
	}

	data, errGo := os.ReadFile(getFullPath(args[0].(*object.String).Value))
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	errGo = os.WriteFile(getFullPath(args[1].(*object.String).Value), data, 0644)
	if errGo != nil {
		return errors.New(errors.RuntimeError, 0, 0, errGo.Error())
	}

	return &object.Null{}
}
