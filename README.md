# Sunbird

Sunbird is dynamically-typed (with optional type annotations), interpreted programming language that focuses on **ease of use** and **clarity**.

For detailed language reference, standard library docs, and guides, see the [`docs/`](./docs) directory.

## Overview

### Hello world in Sunbird
```ts
import "io"

io.println("Hello world")
```

### Defining variables and functions
```ts
let a = 1
const b: Int = 2

const add = func(a: Int, b: Int): Int {
    return a + b
}

add(a, b)
```

### Control flow
```ts
import "io"

let a = 1
let b = 2

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

try {
    let c = 1 / 0
} catch e {
    io.println(e)
} finally {
    io.println("finally")
}

```

### Type annotations
```ts
let a: Int = 1
let b: Float = 2.0
let c: String = "hello"
let d: Bool = true
let e: Void = null
let f: Array = [1, 2, 3]
let g: Func = func(a: Int, b: Int): Int {
  return a + b
}
let h: Hash = {1: 1, 2: 2, 3: 3}


// Nullable types
let i: Int? = null
let j: String? = "hello"
let d: Bool? = true
let e: Array? = [1, 2, 3]
let f: Func? = func(a: Int, b: Int): Int {
    return a + b
}
let g: Hash? = {1: 1, 2: 2, 3: 3}
```
