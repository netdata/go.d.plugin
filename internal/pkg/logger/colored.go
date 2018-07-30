package logger

import (
	"os"
	"strings"

	"fmt"
	"github.com/fatih/color"
)

type colored struct{}

func (c colored) Write(b []byte) (n int, err error) {
	msg := string(b)

	if sevLevel != DEBUG {
		return fmt.Fprint(os.Stderr, msg)
	}

	c.colorWrite(msg)
	return len(b), nil
}

func (c colored) colorWrite(m string) {
	switch {
	case strings.Contains(m, DEBUG.String()):
		color.Magenta(m)
	case strings.Contains(m, INFO.String()):
		color.Green(m)
	case strings.Contains(m, WARNING.String()):
		color.HiYellow(m)
	default:
		color.Red(m)
	}
}
