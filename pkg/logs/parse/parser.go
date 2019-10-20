package parse

import "errors"

type Error struct {
	msg string
	err error
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Unwrap() error {
	return e.err
}

func IsParseError(err error) bool {
	return errors.As(err, &Error{})
}

type (
	LogLine interface {
		Assign(field string, value string) error
	}

	Parser interface {
		ReadLine(LogLine) error
		Parse(line []byte, logLine LogLine) error
		Info() string
	}
)
