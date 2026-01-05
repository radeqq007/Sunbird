# Control flow

## If expressions
If expressions are used to execute code conditionally.
```ts
if x > 0 {
  io.println("x is positive")
} else if x < 0 {
  io.println("x is negative")
} else {
  io.println("x is zero")
}
```

Else if and else are optional.
```ts
if x > 0 {
  io.println("x is positive")
}
```

If expressions can also be used to evaluate expressions.
```ts
let a = if b > 0 { 1 } else { 0 }
```

## For loops
For loops are used to execute code repeatedly.
```ts
for i in 0..10 {
  io.println(i)
}
```

Output:
```
0
1
2
3
4
5
6
7
8
9
```

To define a step you can use the `:` operator after the range.
```ts
for i in 0..10:2 {
  io.println(i)
}
```

Output:
```
0
2
4
6
8
```

You can also loop over arrays and strings.
```ts
let array = [1, 2, 3, 4, 5]
for i in array {
  io.println(i)
}
```

Output:
```
1
2
3
4
5
```

```ts
let string = "hello"
for char in string {
  io.println(char)
}
```

Output:
```
h
e
l
l
o
```

## While loops
While loops are used to execute code repeatedly while a condition is true.
```ts
while x > 0 {
  io.println(x)
  x = x - 1
}
```

## Break and continue
Break and continue are used to control the flow of loops.

Continue will skip the rest of the loop body and continue to the next iteration.
```ts
for i in 0..10 {
  if i % 2 == 0 {
    continue
  }
  io.println(i)
}
```

Break will exit the loop.
```ts
for i in 0..10 {
  if i == 5 {
    break
  }
  io.println(i)
}
```

## Error handling

You can handle potential errors using the `try` statement.
```ts
try {
  io.println("This will print")
} catch e {
  io.println("This will print if an error is thrown")
} finally {
  io.println("This will print no matter what")
}
```

If an error is thrown inside the `try` block, it will be caught by the `catch` block.

The `finally` block will always execute, regardless of whether an error was thrown or not.


To throw an error you can use the `errors` module.

```ts
import "errors"

errors.runtime_error("error message goes here")
```

To learn more about the errors module see the [errors](../std/errors.md) docs.