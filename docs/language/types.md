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