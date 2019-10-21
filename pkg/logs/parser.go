package logs

import (
	"errors"
	"fmt"
	"io"
)

type ParseError struct {
	msg string
	err error
}

func (e ParseError) Error() string { return e.msg }

func (e ParseError) Unwrap() error { return e.err }

func IsParseError(err error) bool { return errors.As(err, &ParseError{}) }

type (
	LogLine interface {
		Assign(field string, value string) error
	}

	Parser interface {
		ReadLine(LogLine) error
		Parse(line []byte, logLine LogLine) error
		Info() string
	}

	Guesser interface {
		Guess(config ParserConfig, in io.Reader) (Parser, error)
	}

	GuessFunc func(config ParserConfig, in io.Reader) (Parser, error)
)

func (f GuessFunc) Guess(config ParserConfig, in io.Reader) (Parser, error) { return f(config, in) }

const (
	TypeAuto   = "auto"
	TypeCSV    = "csv"
	TypeLTSV   = "ltsv"
	TypeRegExp = "regexp"
)

type ParserConfig struct {
	LogType string       `yaml:"log_type"`
	CSV     CSVConfig    `yaml:"csv_config"`
	LTSV    LTSVConfig   `yaml:"ltsv_config"`
	RegExp  RegExpConfig `yaml:"regexp_config"`
}

func NewParser(config ParserConfig, in io.Reader, guess Guesser) (Parser, error) {
	switch config.LogType {
	case TypeAuto:
		if guess == nil {
			return nil, errors.New("guess is nil")
		}
		return guess.Guess(config, in)
	case TypeCSV:
		return NewCSVParser(config.CSV, in)
	case TypeLTSV:
		return NewLTSVParser(config.LTSV, in)
	case TypeRegExp:
		return NewRegExpParser(config.RegExp, in)
	default:
		return nil, fmt.Errorf("invalid type: %q", config.LogType)
	}
}
