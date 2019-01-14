/*
Package matcher implements vary type of string matcher.

Supported match type:
  string
  glob
  regexp
  simple patterns

Short Syntax
  <line> 	  ::= <type> <space> <expr>
  <type>      ::= [ <not> ] <symbol>
  <not>       ::= '!'
                    positive expression
  <symbol>    ::= [ '=', '~', '*' ]
                    '=' means string match
                    '~' means regexp match
                    '*' means glob match
  <space>     ::= { ' ' | '\t' | '\f' | '\v' }
  <expr>      ::= any string

Long Syntax
  <line>      ::= [ <not> ] <name> <separator> <expr>
  <name>      ::= [ 'string' | 'glob' | 'regexp' | 'simple_patterns' ]
  <not>       ::= '!'
                    positive expression
  <separator> ::= ':'
  <expr>      ::= any string
*/
package matcher
