package matcher

import (
	"errors"
	"regexp"
	"strings"
)

const (
	methodString        = "string"
	methodRegexp        = "regexp"
	methodSimplePattern = "simplepattern"
)

type Mather interface {
	Match(string) bool
}

func CreateMatcher(line string) (Mather, error) {
	parts := strings.SplitN(line, "=", 2)

	if len(parts) == 2 && parts[1] == "" || len(parts) != 2 {
		return nil, errors.New("unsupported match syntax")
	}

	method, expr := parts[0], parts[1]

	switch method {
	case methodSimplePattern:
		return CreateSimplePatterns(expr)
	case methodRegexp:
		return createRegexpMatcher(expr)
	case methodString:
		return createStringMatcher(expr), nil
	default:
		return nil, errors.New("unsupported match method")
	}
}

func createStringMatcher(expr string) Mather {
	if strings.HasPrefix(expr, "^") {
		return &StringPrefix{expr[1:]}
	}
	if strings.HasSuffix(expr, "$") {
		return &StringSuffix{expr[:len(expr)-1]}
	}
	return &StringContains{expr}
}

func createRegexpMatcher(expr string) (Mather, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &RegexpMatch{re}, nil
}
