# io

io is a module that provides functions for handling input and output.

```ts
import "io"
```

## print

`print` is a function that prints a value to the console.

```ts
io.print("hello world")
io.print(123)
```

You can provide mulitple values and they will be printed separated by a space.

```ts
io.print("hello", "world")
io.print("hello", 123)
io.print("hello", true)
io.print("hello", 123, true)
io.print("hello", [1, 2, 3], { "a": 123 })
```

## println

println behaves the same as `print`, however it will print a newline at the end.
```ts
io.println("hello world")
io.println(123)
io.println("hello!", [1, 2, 3], { "a": 123 })
```

## read

`read` is a function that reads input from the console.

```ts
io.read()
```

You can provide a string argument to `read` to print it before reading the input.

```ts
io.read("Enter your name: ")
```

>[!NOTE]
> `read` will return the input including the newline character.
> To read input without the newline character, use `readln`.

## readln

`readln` is a function that reads a line from the console.

```ts
io.readln()
```

You can provide a string argument to `readln` to print it before reading the input.

```ts
io.readln("Enter your name: ")
```

# printf

`printf` is a function that prints a formatted string to the console.

The format string uses curly braces `{}` to mark where arguments should be inserted.

```ts
io.printf("hello {}", "world")
io.printf("hello {}", 123)
io.printf("hello {} {}", "world", 123)
```

# printfn

`printfn` is a function that prints a formatted string to the console with a newline at the end.

The format string uses curly braces `{}` to mark where arguments should be inserted.

```ts
io.printfn("hello {}", "world")
io.printfn("hello {} {}", "world", 123)
```

# sprintf

`sprintf` is a function that returns a formatted string.

The format string uses curly braces `{}` to mark where arguments should be inserted.

```ts
io.sprintf("hello {} {}", "world", 123)
```

The function returns a string.

## clear

`clear` is a function that clears the console.

```ts
io.clear()
```

## beep

`beep` is a function that emits a beep sound.

```ts
io.beep()
```

## args

`args` is an array of command line arguments.

```ts
let args_array = io.args
let first_arg = io.args[0]
```

The array is of strings.
