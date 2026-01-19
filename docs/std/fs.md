# fs

`fs` is a module for file system operations.

```ts
import "fs"
```

## read

`read` is a function that reads a file and returns its contents as a string.

```ts
let content = null
try {
  content = fs.read("file.txt")
} catch e {
  io.println(e)
}
```

## write

`write` is a function that writes a string to a file.

```ts
try {
  fs.write("file.txt", "Hello, world!")
} catch e {
  io.println(e)
}
```

## append

`append` is a function that appends a string to a file.

```ts
try {
  fs.append("file.txt", "Hello, world!")
} catch e {
  io.println(e)
}
```

## remove

`remove` is a function that removes a file or directory.

```ts
try {
  fs.remove("file.txt")
} catch e {
  io.println(e)
}
```

## exists

`exists` is a function that checks if a file or directory exists and returns a boolean.

```ts
let exists = false
try {
  exists = fs.exists("file.txt")
} catch e {
  io.println(e)
}
```

## is_dir

`is_dir` is a function that checks if a path is a directory and returns a boolean.

```ts
let is_dir = false
try {
  is_dir = fs.is_dir("file.txt")
} catch e {
  io.println(e)
}
```

## list_dir

`list_dir` is a function that lists the contents of a directory and returns an array of strings.

```ts
let files = []
try {
  files = fs.list_dir("./")
} catch e {
  io.println(e)
}
```

## create_dir

`create_dir` is a function that creates a directory.

```ts
try {
  fs.create_dir("./dir")
} catch e {
  io.println(e)
}
```

## rename

`rename` is a function that renames a file or directory.

```ts
try {
  fs.rename("file.txt", "file2.txt")
} catch e {
  io.println(e)
}
```

## copy

`copy` is a function that copies a file.

```ts
try {
  fs.copy("file.txt", "file2.txt")
} catch e {
  io.println(e)
}
```

