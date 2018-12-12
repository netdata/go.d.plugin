package matcher

import (
	"fmt"
	"strings"
)

const (
	methodString = "string"
	methodRegexp = "regexp"
)

type Matcher interface {
	Match(string) bool
}

// Valid options:
// 'string=GET'
// 'string=^GOT'
// 'regexp=G[QWERTY]T'
func New(rawExpr string) (Matcher, error) {
	v := strings.SplitN(rawExpr, "=", 2)

	if len(v) == 2 && v[1] == "" || len(v) != 2 {
		return nil, fmt.Errorf("unsupported Match syntax: %s", rawExpr)
	}

	method, expr := v[0], v[1]

	switch method {
	case methodString:
		return stringMatcherFactory(expr), nil
	case methodRegexp:
		return regexpMatcherFactory(expr)
	}

	return nil, fmt.Errorf("unsupported Match method: %s", method)
}
