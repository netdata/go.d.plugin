package matcher

import (
	"bytes"
	"strings"
)

type (
	// stringFullMatcher implements Matcher, it uses "==" to match.
	stringFullMatcher string

	// stringPartialMatcher implements Matcher, it uses strings.Contains to match.
	stringPartialMatcher string

	// stringPrefixMatcher implements Matcher, it uses strings.HasPrefix to match.
	stringPrefixMatcher string

	// stringSuffixMatcher implements Matcher, it uses strings.HasSuffix to match.
	stringSuffixMatcher string

	// stringLenMatcher implements Matcher, it uses len(s) == m to match.
	stringLenMatcher int
)

func (m stringFullMatcher) Match(b []byte) bool          { return string(m) == string(b) }
func (m stringFullMatcher) MatchString(line string) bool { return string(m) == line }

func (m stringPartialMatcher) Match(b []byte) bool          { return bytes.Contains(b, []byte(m)) }
func (m stringPartialMatcher) MatchString(line string) bool { return strings.Contains(line, string(m)) }

func (m stringPrefixMatcher) Match(b []byte) bool          { return bytes.HasPrefix(b, []byte(m)) }
func (m stringPrefixMatcher) MatchString(line string) bool { return strings.HasPrefix(line, string(m)) }

func (m stringSuffixMatcher) Match(b []byte) bool          { return bytes.HasSuffix(b, []byte(m)) }
func (m stringSuffixMatcher) MatchString(line string) bool { return strings.HasSuffix(line, string(m)) }

func (m stringLenMatcher) Match(b []byte) bool          { return len(b) == int(m) }
func (m stringLenMatcher) MatchString(line string) bool { return len(line) == int(m) }
