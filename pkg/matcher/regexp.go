package matcher

import "regexp"

func NewRegExpMatcher(expr string) (Matcher, error) {
	switch expr {
	case "", "^", "$":
		return TRUE(), nil
	case "^$", "$^":
		return stringFullMatcher(""), nil
	}
	size := len(expr)
	chars := []rune(expr)
	var startWith, endWith bool
	startIdx := 0
	endIdx := size - 1
	if chars[startIdx] == '^' {
		startWith = true
		startIdx = 1
	}
	if chars[endIdx] == '$' {
		endWith = true
		endIdx--
	}

	unescapedExpr := make([]rune, 0, endIdx-startIdx+1)
	for i := startIdx; i <= endIdx; i++ {
		ch := chars[i]
		if ch == '\\' {
			if i == endIdx { // end with '\' => invalid format
				return regexp.Compile(expr)
			}
			nextCh := chars[i+1]
			if !isRegExpMeta(nextCh) { // '\' + mon-meta char => special meaning
				return regexp.Compile(expr)
			}
			unescapedExpr = append(unescapedExpr, nextCh)
			i++
		} else if isRegExpMeta(ch) {
			return regexp.Compile(expr)
		} else {
			unescapedExpr = append(unescapedExpr, ch)
		}
	}

	if startWith {
		if endWith {
			return stringFullMatcher(string(unescapedExpr)), nil
		}
		return stringPrefixMatcher(string(unescapedExpr)), nil
	}
	if endWith {
		return stringSuffixMatcher(string(unescapedExpr)), nil
	}
	return stringPrefixMatcher(string(unescapedExpr)), nil
}

// isRegExpMeta reports whether byte b needs to be escaped by QuoteMeta.
func isRegExpMeta(b rune) bool {
	switch b {
	case '\\', '.', '+', '*', '?', '(', ')', '|', '[', ']', '{', '}', '^', '$':
		return true
	default:
		return false
	}
}
