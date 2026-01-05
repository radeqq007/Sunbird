# random

`random` is a module used for generating random numbers.

```ts
import "random"
```

## int

`int` is a function used for generating random integers between two provided integers.

```ts
random.int(min, max)
```

This returns a random integer between `min` (inclusive) and `max` (exclusive).

## float

`float` is a function used for generating random floating point numbers between two provided floating point numbers.

```ts
random.float(min, max)
```

This returns a random floating point number between `min` (inclusive) and `max` (exclusive).

## bool

`bool` is a function used for generating random boolean values.

```ts
random.bool()
```

This returns either `true` or `false`.

## seed

`seed` is a function used to set the seed for the random number generator.
It takes an integer as an argument.

```ts
random.seed(seed)
```

## shuffle

`shuffle` is a function used for shuffling the elements of an array.

```ts
random.shuffle(array)
```

This returns a new array with the elements of the provided in randomized order.

## choice

`choice` is a function used for selecting a random element from an array.

```ts
random.choice(array)
```

This returns a random element from the provided array.
