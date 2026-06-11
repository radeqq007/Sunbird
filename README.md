<div align="center">
  <img src=".github/assets/LOGO.svg" width="320" alt="Sunbird">
</div>

# Sunbird
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/radeqq007/Sunbird)
[![Build](https://github.com/radeqq007/Sunbird/actions/workflows/build.yml/badge.svg)](https://github.com/radeqq007/Sunbird/actions/workflows/build.yml)
[![Lint](https://github.com/radeqq007/Sunbird/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/radeqq007/Sunbird/actions/workflows/golangci-lint.yml)
[![Tests](https://github.com/radeqq007/Sunbird/actions/workflows/tests.yml/badge.svg)](https://github.com/radeqq007/Sunbird/actions/workflows/tests.yml)
[![Coverage](https://codecov.io/gh/radeqq007/sunbird/branch/main/graph/badge.svg)](https://codecov.io/gh/radeqq007/sunbird)
[![Go Report Card](https://goreportcard.com/badge/github.com/radeqq007/sunbird)](https://goreportcard.com/report/github.com/radeqq007/sunbird)
![Last commit](https://img.shields.io/github/last-commit/radeqq007/Sunbird)
![GitHub stars](https://img.shields.io/github/stars/radeqq007/sunbird?style=social)

Sunbird is dynamically-typed, interpreted programming language that focuses on **ease of use** and **simplicity**.

For detailed language reference, standard library docs, and guides, see the [`docs/`](./docs) directory.

## Overview

### Hello world in Sunbird
```ts
import "io"

io.println("Hello world")
```

### Defining variables and functions
```rs
// Use the := operator to declare a variable
a := 1
name := ""

// to define constants use the :: operator
x :: 20
pi :: 3.141

// Defining functions
add :: fn(a, b) {
    return a + b
}

add(a, b)
```

### Control flow
```rs
import "io"

a := 1
b := 2

if a > b {
    io.println("a is greater than b")
} else if a < b {
    io.println("a is less than b")
} else {
    io.println("a is equal to b")
}

for i in 0..10 {
    io.println(i)
}

while a <= b {
    io.println(a)
    a += 1

    if a == 1 {
      continue
    }

    if a == 2 {
      break
    }
}

loop {
    io.println("This will run forever!")
}

try {
    c := 1 / 0
} catch e {
    io.println(e)
} finally {
    io.println("finally")
}
```
