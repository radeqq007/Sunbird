# Syntax

## Variables
You can declare a variable using the `:=` operator or the `::` operator for constants.

```rs
x := 5
y :: 10
```

## Functions

Functions are expressions and can be assigned to variables or passed as arguments to other functions and are declared using the `fn` keyword.

```ts
add :: fn(a, b) {
  return a + b
}
```

To call a function you can use the function name followed by the arguments in parentheses.
```ts
add(1, 2)
```

> [!TIP]
> You can also omit the `return` keyword if the last expression is the return value.
>
> ```rs
> add :: fn(a, b) { a + b }
> ```
>

## Control flow
If expressions are used to execute code conditionally.
```ts
import "io"

if x > 0 {
  io.println("x is positive")
} else if x < 0 {
  io.println("x is negative")
} else {
  io.println("x is zero")
}
```

For loops are used to execute code repeatedly.
```rs
for i in 0..10 {
  io.println(i)
}
```

To define a step you can use the `:` operator after the range.
```rs
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

While loops are used to execute code repeatedly while a condition is true.
```ts
x := 10
while x > 0 {
  io.println(x)
  x = x - 1
}
```

To learn more about control flow see the [control flow](./control-flow.md) docs.

## Error handling
Try catch blocks are used to handle errors.
```ts
try {
  io.println("This will print")
} catch e {
  io.println("This will print if an error is thrown")
}
```




To learn more about types see the [types](./types.md) docs.

## Comments

Comments in Sunbird are denoted with `//` for single line comments and `/* */` for multi line comments.

```ts
// This is a single line comment

/*
This is a
multi line comment
*/
```

## Multiple Statements

You can write multiple expressions or statements on the same line by separating them with a semicolon.
```ts
x := 5; y := 10
```
