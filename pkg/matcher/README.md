# Supported Format

  * string
  * glob
  * regexp
  * simple patterns

Depending on the symbol at the start of the string, the `matcher` will use one of the supported formats. 

e.g `'* *'`: Because of the first `*`, `matcher` understands that it's in `glob` format, and according to the `glob` format the second `*` means that it will match *any* string. 
  
### Syntax

**Tip**: Read `::=` as `is defined as`.
```
Short Syntax
     [ <not> ] <format> <space> <expr>
     
     <not>       ::= '!'
                       negative expression
     <format>    ::= [ '=', '~', '*' ]
                       '=' means string match
                       '~' means regexp match
                       '*' means glob match
     <space>     ::= { ' ' | '\t' | '\n' | '\n' | '\r' }
     <expr>      ::= any string

 Long Syntax

     [ <not> ] <format> <separator> <expr>
     <format>    ::= [ 'string' | 'glob' | 'regexp' | 'simple_patterns' ]
     <not>       ::= '!'
                       negative expression
     <separator> ::= ':'
     <expr>      ::= any string
```

**Tip**: In short syntax, you can enable the glob format by starting the string with a `*`, while in the long syntax you need to define it more explicitly. The following examples are identical. `simple_patterns` can be used **only** with the long syntax.

Short Syntax: `'* * '`
Long Syntax: `'glob':'*'`

### String matcher
The string matcher reports whether the given value equals to the string ( use == ).

e.g <br>
`'== foo'` will match only if the string is `'foo'`. <br>
`'[!] == bar'` will match any string that is not `'bar'`.


### Glob matcher

The glob matcher reports whether the given value matches the wildcard pattern. It uses the standard `golang` library `path`. You can read more about the library in the [golang documentation](https://golang.org/pkg/path/#Match), where you can also practice with the library in order to learn the syntax and use it in your netdata configuration.

The pattern syntax is:
```
    pattern:
        { term }
    term:
        '*'         matches any sequence of characters
        '?'         matches any single character
        '[' [ '^' ] { character-range } ']'
        character class (must be non-empty)
        c           matches character c (c != '*', '?', '\\', '[')
        '\\' c      matches character c

    character-range:
        c           matches character c (c != '\\', '-', ']')
        '\\' c      matches character c
        lo '-' hi   matches character c for lo <= c <= hi
```
e.g 
 - `* ?` will match any string that is a single character. 
 - `'?a'` will match any 2 character string that starts with any character and the second character is `a`, like `ba` but not `bb` or `bba`. 
 - `'[^abc]'` will match any character that is NOT a,b,c. `'[abc]'` will match only a, b, c.
 - `'*[a-d]'` will match any string (`*`) that ends with a character that is between `a` and `d` (i.e `a,b,c,d`).


      
### Regexp matcher
The regexp matcher reports whether the given value matches the RegExp pattern ( use regexp.Match ).

The RegExp syntax is described at https://golang.org/pkg/regexp/syntax/.

We know that Regular Expressions are hard, but they are powerfull. If you want to learn more about them, you can visit [RegexOne](https://regexone.com/).

### Simple patterns matcher
The simple patterns matcher reports whether the given value matches the simple patterns.

Simple patterns are a space separated list of words, that can have `*` as a wildcard. Each world may use any number of `*`. Simple patterns allow negative matches by prefixing a word with `!`.

So, pattern = `!*bad* *` will match anything, except all those that contain the word bad.

Simple patterns are quite powerful: pattern = `*foobar* !foo* !*bar *` matches everything containing foobar, except strings that start with foo or end with bar.




