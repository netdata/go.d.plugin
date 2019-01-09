package matcher

import "strings"

type StringContains struct{ SubStr string }

func (m StringContains) Match(line string) bool { return strings.Contains(line, m.SubStr) }

type StringPrefix struct{ Prefix string }

func (m StringPrefix) Match(line string) bool { return strings.HasPrefix(line, m.Prefix) }

type StringSuffix struct{ Suffix string }

func (m StringSuffix) Match(line string) bool { return strings.HasSuffix(line, m.Suffix) }
