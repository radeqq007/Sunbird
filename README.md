<p align="center">
<img width="250" alt="Sunbird" src="https://github.com/user-attachments/assets/8c19c7b2-4d08-4d9f-a1da-0da1bb72ac5c" />
</p>

# The Sunbird programming language
[![Build](https://github.com/Sunbird-Lang/Sunbird/actions/workflows/build.yml/badge.svg)](https://github.com/Sunbird-Lang/Sunbird/actions/workflows/build.yml)
[![Tests](https://github.com/Sunbird-Lang/Sunbird/actions/workflows/tests.yml/badge.svg)](https://github.com/Sunbird-Lang/Sunbird/actions/workflows/tests.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
![Repo size](https://img.shields.io/github/repo-size/sunbird-lang/sunbird?cacheSeconds=0)

Sunbird is a simple, interpreted, dynamically typed language.

# Documentation
## Comments
Single line comments start with `//` and continue to the end of the line.
```go
// This is a single line comment
```

Block comments are enclosed within `/* ... */` and can span multiple lines.
```go
/*
This comment
spans multiple
lines
*/
```

## Declaring variables
In Sunbird you declare variables using the `var` keyword:
```go
var foo = "Hello, World!"
```

## Data types
Sunbird supports all the basic data types:
```go
var str = "hello!" // string
var int = 10 // integer
var float = 3.14 // float
var bool = true // booleans
var foo = null // null
```

<br />

## Arrays:

Arrays are an ordered list of elements of possibly different types identified by a number index. Each element in an array can be accessed individually by their index. Arrays are constructed as a comma separated list of elements, can contain any type of value, and are enclosed by square brackets:
```go
var arr = [ 1, "sunbird", 10, 2.2]
```

To get a value from an array you use the bracket notation:
```go
arr[0]
```

Negative indices can be used to access elements from the end of the array:
```go
arr[-1] // Returns the last element from the array
```

## Functions
Functions in Sunbird are defined using the func keyword:
```go
var add = func(a, b) {
  return a + b
}
```

You can also skip the `return`:
```go
var add = func(a, b) {
  a + b
}
```

Functions can be called like this:
```go
var result = add(10, 5) // 15
```

## Conditional statements
Sunbird supports `if`, `else if`, and `else` statements:

```go
var x = 5
var y = 10

if x > y {
    println("x is greater than y")
} else if x == y {
    println("x is equal to y")
} else {
    println("x is less than y")
}
```

## For loops
A `for` loop is used to execute a block of code multiple times.

```go
  for i = 0; i < 10; i = i + 1 {
    println("This code will run 10 times")
  }
```
## Pipe operator
Embedded function calls can get messsy and hard to follow. For example:
```go
var result = foo(bar(baz(another_func(data))))
```
in that case you can use the **pipe operator**:
```go
var result = data |> another_func |> baz |> bar |> foo
```

*Note: documentation is work in progress*
