package parser

type Parser interface {
	Parse(line string) (GroupMap, bool)
}
