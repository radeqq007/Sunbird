# Basic Concepts

In this page, you’ll learn the core concepts of Sunbird, including **variables**, **functions**, and **control flow**.

For more details, see the [Syntax](../language/syntax.md) page.

---

## Variables

Variables hold values. You can declare them with the `:=` operator or the `::` operator for constants.

```rs
age := 18
name :: "Bojack"

age = 19       // ✅ mutable variable can be updated
name = "Diane" // ❌ this will throw an error
```

## Functions

Functions are reusable blocks of code. You declare them with the `fn` keyword.

```rs
add :: fn(a, b) {
  return a + b
}

result := add(5, 3)
io.println(result) // 8
```

Functions can be passed around like variables:
```rs
double :: fn(x) { x * 2 }

applyFunc :: fn(func, value) { func(value) }

io.println(applyFunc(double, 10)) // 20
```

## Control Flow

### If Expression
Use `if` to execute code conditionally:

```rs
x := -5

if x > 0 {
    io.println("x is positive")
} else if x < 0 {
    io.println("x is negative")
} else {
    io.println("x is zero")
}

```

`if` can also return a value:

```rs
sign := if x > 0 { 1 } else { -1 }
io.println(sign) // -1
```

### Loops

For loops:
```rs
for i in 1..10 {
    io.println(i)
}
```

While loops:
```rs
i := 0
while i < 5 {
    io.println(i)
    i += 1
}
```

## Next steps
Next, you can check out the guides for:
- [Syntax](../language/syntax.md)
- [Control Flow](../language/control-flow.md)
- [Types](../language/types.md)
