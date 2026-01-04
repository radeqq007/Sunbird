# Basic Concepts

In this page, you’ll learn the core concepts of Sunbird, including **variables**, **functions**, **types**, and **control flow**.

For more details, see the [Types](../language/types.md) and [Syntax](../language/syntax.md) pages.

---

## Variables

Variables hold values. You can declare them with `let` (mutable) or `const` (immutable).

```ts
let age = 18
const name = "Bojack"

age = 19        // ✅ mutable variable can be updated

name = "Diane" // ❌ this will throw an error
```
> [!NOTE]
> You can also add a type declaration: `let age: Int = 18`

## Functions

Functions are reusable blocks of code. You declare them with the `func` keyword.

```ts
let add = func(a: Int, b: Int): Int {
  return a + b
}

let result = add(5, 3)
io.println(result) // 8
```

Functions can be passed around like variables:
```ts
let double = func(x: Int): Int { x * 2 }

let applyFunc = func(f: Func, value: Int): Int { f(value) }

io.println(applyFunc(double, 10)) // 20
```

## Control Flow

### If Expression
Use `if` to execute code conditionally:

```ts
let x = -5

if x > 0 {
    io.println("x is positive")
} else if x < 0 {
    io.println("x is negative")
} else {
    io.println("x is zero")
}

```
`if` can also return a value:
```ts
let sign = if x > 0 { 1 } else { -1 }
io.println(sign) // -1
```

### Loops

For loops:
```ts
for i in 1..5 {
    io.println(i)
}
```

While loops:
```ts
let i = 0
while i < 5 {
    io.println(i)
    i = i + 1
}
```

## Next steps
Next, you can check out the guides for:
- [Syntax](../language/syntax.md)
- [Control Flow](../language/control-flow.md)
- [Types](../language/types.md)