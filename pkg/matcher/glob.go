package matcher

import (
	"regexp"
	"runtime"
	"unicode/utf8"

	"errors"
)

// globMatcher implements Matcher, it uses filepath.MatchString to match.
type globMatcher string

var (
	errBadGlobPattern = errors.New("bad glob pattern")
	erGlobPattern     = regexp.MustCompile(`(?s)^(?:[*?]|\[\^?([^\\-\]]|\\.|.-.)+\]|\\.|[^\*\?\\\[])*$`)
)

// NewGlobMatcher create a new matcher with glob format
func NewGlobMatcher(expr string) (Matcher, error) {
	switch expr {
	case "":
		return stringFullMatcher(""), nil
	case "*":
		return TRUE(), nil
	}

	// any strings pass this regexp check are valid pattern
	if !erGlobPattern.MatchString(expr) {
		return nil, errBadGlobPattern
	}

	size := len(expr)
	chars := []rune(expr)
	startWith := true
	endWith := true
	startIdx := 0
	endIdx := size - 1
	if chars[startIdx] == '*' {
		startWith = false
		startIdx = 1
	}
	if chars[endIdx] == '*' {
		endWith = false
		endIdx--
	}

	unescapedExpr := make([]rune, 0, endIdx-startIdx+1)
	for i := startIdx; i <= endIdx; i++ {
		ch := chars[i]
		if ch == '\\' {
			nextCh := chars[i+1]
			unescapedExpr = append(unescapedExpr, nextCh)
			i++
		} else if isGlobMeta(ch) {
			return globMatcher(expr), nil
		} else {
			unescapedExpr = append(unescapedExpr, ch)
		}
	}

	return NewStringMatcher(string(unescapedExpr), startWith, endWith)
}

func isGlobMeta(ch rune) bool {
	switch ch {
	case '*', '?', '[':
		return true
	default:
		return false
	}
}

// MatchString matches.
func (m globMatcher) Match(b []byte) bool {
	return m.MatchString(string(b))
}

// MatchString matches.
func (m globMatcher) MatchString(line string) bool {
	return globMatch(string(m), line)
}

const winOS = "windows"

func globMatch(pattern, name string) bool {
Pattern:
	for len(pattern) > 0 {
		var star bool
		var chunk string
		star, chunk, pattern = scanChunk(pattern)
		if star && chunk == "" {
			// Trailing * matches rest of string unless it has a /.
			// return !strings.Contains(name, string(Separator)), nil

			return true
		}
		// Look for match at current position.
		t, ok := matchChunk(chunk, name)
		// if we're the last chunk, make sure we've exhausted the name
		// otherwise we'll give a false result even if we could still match
		// using the star
		if ok && (len(t) == 0 || len(pattern) > 0) {
			name = t
			continue
		}
		if star {
			// Look for match skipping i+1 bytes.
			// Cannot skip /.
			for i := 0; i < len(name); i++ {
				//for i := 0; i < len(name) && name[i] != Separator; i++ {
				t, ok := matchChunk(chunk, name[i+1:])
				if ok {
					// if we're the last chunk, make sure we exhausted the name
					if len(pattern) == 0 && len(t) > 0 {
						continue
					}
					name = t
					continue Pattern
				}
			}
		}
		return false
	}
	return len(name) == 0
}

// scanChunk gets the next segment of pattern, which is a non-star string
// possibly preceded by a star.
func scanChunk(pattern string) (star bool, chunk, rest string) {
	for len(pattern) > 0 && pattern[0] == '*' {
		pattern = pattern[1:]
		star = true
	}
	inrange := false
	var i int
Scan:
	for i = 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			if runtime.GOOS != winOS {
				// error check handled in matchChunk: bad pattern.
				if i+1 < len(pattern) {
					i++
				}
			}
		case '[':
			inrange = true
		case ']':
			inrange = false
		case '*':
			if !inrange {
				break Scan
			}
		}
	}
	return star, pattern[0:i], pattern[i:]
}

// matchChunk checks whether chunk matches the beginning of s.
// If so, it returns the remainder of s (after the match).
// Chunk is all single-character operators: literals, char classes, and ?.
func matchChunk(chunk, s string) (rest string, ok bool) {
	for len(chunk) > 0 {
		switch chunk[0] {
		case '[':
			// character class
			r, n := utf8.DecodeRuneInString(s)
			s = s[n:]
			chunk = chunk[1:]
			// possibly negated
			negated := chunk[0] == '^'
			if negated {
				chunk = chunk[1:]
			}
			// parse all ranges
			match := false
			nrange := 0
			for {
				if len(chunk) > 0 && chunk[0] == ']' && nrange > 0 {
					chunk = chunk[1:]
					break
				}
				var lo, hi rune
				lo, chunk = getEsc(chunk)
				hi = lo
				if chunk[0] == '-' {
					hi, chunk = getEsc(chunk[1:])
				}
				if lo <= r && r <= hi {
					match = true
				}
				nrange++
			}
			if match == negated {
				return
			}

		case '?':
			_, n := utf8.DecodeRuneInString(s)
			s = s[n:]
			chunk = chunk[1:]

		case '\\':
			if runtime.GOOS != winOS {
				chunk = chunk[1:]
			}
			fallthrough

		default:
			if chunk[0] != s[0] {
				return
			}
			s = s[1:]
			chunk = chunk[1:]
		}
	}
	return s, true
}

// getEsc gets a possibly-escaped character from chunk, for a character class.
func getEsc(chunk string) (r rune, nchunk string) {
	if chunk[0] == '\\' && runtime.GOOS != winOS {
		chunk = chunk[1:]
	}
	r, n := utf8.DecodeRuneInString(chunk)
	nchunk = chunk[n:]
	return
}
