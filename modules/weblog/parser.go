package weblog

type Parser interface {
	Parse(s string) (map[string]string, bool)
}
