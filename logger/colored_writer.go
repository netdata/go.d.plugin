package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type coloredWriter struct{}

func (c *coloredWriter) Write(b []byte) (n int, err error) {
	msg := string(b)
	switch sevLevel {
	case DEBUG:
		switch {
		case strings.Contains(msg, DEBUG.String()):
			color.Magenta(msg)
		case strings.Contains(msg, INFO.String()):
			color.Green(msg)
		case strings.Contains(msg, WARNING.String()):
			color.HiYellow(msg)
		default:
			color.Red(msg)
		}
	default:
		return fmt.Fprint(os.Stderr, msg)
	}
	return len(b), nil
}
