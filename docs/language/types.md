# Types

Sunbird is a dynamically typed language, but you can declare the type of a variable using the `:` operator after the variable name.
If an assigned value is of a different type, you will get a runtime error.

There are several builtin data types, including simple types, data structures, and nullable types.

## Simple types
`Int` for 64 bit signed integers.

```ts
let age: Int = 18
```

`Float` for 64 bit floating point numbers.

```ts
let height: Float = 1.8
```

`String` for strings.

```ts
let name: String = "Bojack"
```

`Bool` for booleans.

```ts
let is_adult: Bool = true
```

`void` for null.

```ts
let x: void = null
```

## Data structures
`Array` for arrays.

```ts
let numbers: Array = [1, 2, 3]
```

`Hash` for hash maps.
```ts
let map: Hash = {"a": 1, "b": 2, "c": 3}
```

`Func` for functions.
```ts
let add: Func = func(a: Int, b: Int): Int { return a + b }
```

## Nullable types
To create a nullable type, you can use the `?` operator following the type name.
```ts
let age: Int? = null
let name: String? = null
```

## Type Conversion

To convert a value to a different type, you can use the the builtin functions `int()`, `float()`, `string()` and `bool()`.

### Integer conversion
```ts
let a: Int = int("5")   // Evaluates to 5
let b: Int = int(1.2)   // Will truncate the decimal part
let c: Int = int(true)  // Evaluates to 1
let d: Int = int(false) // Evaluates to 0
```

### Float conversion
```ts
let a: Float = float("5")   // evaluates to 5.0
let b: Float = float("5.2") // evaluates to 5.2
let c: Float = float(true)  // evaluates to 1.0
let d: Float = float(false) // evaluates to 0.0
```

### String conversion
```ts
let a: String = string(1)    // evaluates to "1"
let b: String = string(true) // evaluates to "true"
let c: String = string(false) // evaluates to "false"
let d: String = string(null) // evaluates to "null"
let e: String = string([1, 2, 3]) // evaluates to "[1, 2, 3]"
let f: String = string({ "a": 1, "b": "Hi!" }) // evaluates to "{\"a\": 1, \"b\": \"Hi!\"}"
```

### Boolean conversion
```ts
let a: Bool = bool(1)    // evaluates to true
let b: Bool = bool(0)    // evaluates to false
let c: Bool = bool("true") // evaluates to true
let d: Bool = bool("false") // also evaluates to true (strings are only evaluated to false if they are empty)
let e: Bool = bool("") // evaluates to false 
let f: Bool = bool(null) // evaluates to false
```