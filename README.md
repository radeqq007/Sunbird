<p align="center">
<img width="250" alt="Sunbird" src="https://github.com/user-attachments/assets/8c19c7b2-4d08-4d9f-a1da-0da1bb72ac5c" />
</p>

# The Sunbird programming language
[![Tests](https://github.com/Sunbird-Lang/Sunbird/actions/workflows/tests.yml/badge.svg)](https://github.com/Sunbird-Lang/Sunbird/actions/workflows/tests.yml)

Sunbird is a simple, interpreted, dynamically typed language.

# Documentation
## Comments
Single line comments start with `//` and continue to the end of the line.
```go
// This is a single line comment
```

Block comments are enclosed within /* ... */ and can span as many lines as necessary.
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

Arrays:

Arrays are an ordered list of elements of possibly different types identified by a number index. Each element in an array can be accessed individually by their index. Arrays are constructed as a comma separated list of elements, can contain any type of value, and are enclosed by square brackets:
```go
var arr = [ 1, "sunbird", 10, 2.2]
```

To get a value from an array you use the bracket notation:
```go
arr[0]
```

You can also use negative index, to get values from the end of the array:
```go
arr[-1] // Returns the last element from the array
```

*Note: documentation is yet to be finished*



