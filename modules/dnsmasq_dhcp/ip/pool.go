package ip

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Type is the IP Pool type.
type Type int

const (
	UnknownType Type = iota
	V4Type
	V6Type
)

// Pool represents IP Pool.
type Pool interface {
	Type() Type
	Contains(ip net.IP) bool
	Hosts() *big.Int
	fmt.Stringer
}

var replacer = strings.NewReplacer(" ", "", "-", ",")

func Parse(s string) Pool {
	if s = parseRawLine(s); s == "" {
		return nil
	}
	s = normalize(s)

	if strings.Contains(s, ",") {
		return ParseRange(s)
	}
	return ParseNet(s)
}

// TODO:
func parseRawLine(s string) string { return s }

func normalize(s string) string { return replacer.Replace(strings.TrimSpace(s)) }
