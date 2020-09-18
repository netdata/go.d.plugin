# stm

This package helps you to convert a struct to `map[string]int64`.

## Tags

The encoding of each struct field can be customized by the format string stored under the `stm` key in the struct field's tag.
The format string gives the name of the field, possibly followed by a comma-separated list of options.

**Lack of tag means no conversion performed.**
If you don't want a field to be added to the `map[string]int64` just don't add a tag to it.
 
Tag syntax:

```
`stm:"name,multiplier,divisor"`
```

Both `multiplier` and `divisor` are optional, `name` is mandatory.

Examples of struct field tags and their meanings:

```
// Field appears in map as key "name".
Field int `stm:"name"`

// Field appears in map as key "name" and its value is multiplied by 10.
Field int `stm:"name,10"`

// Field appears in map as key "name" and its value is multiplied by 10 and divided by 5.
Field int `stm:"name,10,5"`
```

## Supported field value kinds

The list is:

-   `int`
-   `float`
-   `bool`
-   `map`
-   `array`
-   `slice`
-   `pointer`
-   `struct`
-   `interface { WriteTo(rv map[string]int64, key string, mul, div int) }`

It is ok to have nested structures.

## Usage

Use `ToMap` function. Keep in mind:

-   this function is variadic (can be called with any number of trailing arguments).
-   it doesn't allow to have duplicate in result map.
-   if there is a duplicate key it panics.
