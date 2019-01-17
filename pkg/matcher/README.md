# Supported Format

  * string
  * glob
  * regexp
  * simple patterns
  
### Syntax
```
Short Syntax

     <line>      ::= [ <not> ] <format> <space> <expr>
     <not>       ::= '!'
                       negative expression
     <format>    ::= [ '=', '~', '*' ]
                       '=' means string match
                       '~' means regexp match
                       '*' means glob match
     <space>     ::= { ' ' | '\t' | '\n' | '\n' | '\r' }
     <expr>      ::= any string

 Long Syntax

     <line>      ::= [ <not> ] <format> <separator> <expr>
     <format>    ::= [ 'string' | 'glob' | 'regexp' | 'simple_patterns' ]
     <not>       ::= '!'
                       negative expression
     <separator> ::= ':'
     <expr>      ::= any string
```

### String matcher
The string matcher reports whether the given value equals to the string ( use == ).

### Glob matcher
The glob matcher reports whether the given value matches the wildcard pattern.

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
      
### Regexp matcher
The regexp matcher reports whether the given value matches the RegExp pattern ( use regexp.Match ).

The RegExp syntax is described at https://golang.org/pkg/regexp/syntax/.

### Simple patterns matcher
The simple patterns matcher reports whether the given value matches the simple patterns.

The simple patterns is a custom format used in netdata, it's syntax is described at https://docs.netdata.cloud/libnetdata/simple_pattern/.



