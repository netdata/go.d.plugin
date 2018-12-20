package csvparser

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const quoteLen = len(`"`)

var (
	ErrBareQuote  = errors.New("bare \" in non-quoted-field")
	ErrQuote      = errors.New("extraneous or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

type Parser struct {
	LazyQuotes       bool
	TrimLeadingSpace bool
	Comma            rune
	Comment          rune
	FieldsPerRecord  int

	recordBuffer []byte
	fieldIndexes []int
	lastRecord   []string
}

func (r *Parser) ParseString(line string) ([]string, error) {
	if r.Comment != 0 && nextRunesString(line) == r.Comment {
		return nil, nil // Skip Comment lines
	}
	if len(line) == lengthNLString(line) {
		return nil, nil // Skip empty lines
	}

	var err error
	commaLen := utf8.RuneLen(r.Comma)
	r.recordBuffer = r.recordBuffer[:0]
	r.fieldIndexes = r.fieldIndexes[:0]

parseField:
	for {
		if r.TrimLeadingSpace {
			line = strings.TrimLeftFunc(line, unicode.IsSpace)
		}
		if len(line) == 0 || line[0] != '"' {
			// Non-quoted string field
			i := strings.IndexRune(line, r.Comma)
			field := line
			if i >= 0 {
				field = field[:i]
			} else {
				field = field[:len(field)-lengthNLString(field)]
			}
			// Check to make sure a quote does not appear in field.
			if !r.LazyQuotes {
				if j := strings.IndexByte(field, '"'); j >= 0 {
					err = ErrBareQuote
					break parseField
				}
			}
			r.recordBuffer = append(r.recordBuffer, field...)
			r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
			if i >= 0 {
				line = line[i+commaLen:]
				continue parseField
			}
			break parseField
		} else {
			// Quoted string field
			line = line[quoteLen:]
			for {
				i := strings.IndexByte(line, '"')
				if i >= 0 {
					// Hit next quote.
					r.recordBuffer = append(r.recordBuffer, line[:i]...)
					line = line[i+quoteLen:]
					switch rn := nextRunesString(line); {
					case rn == '"':
						// `""` sequence (append quote).
						r.recordBuffer = append(r.recordBuffer, '"')
						line = line[quoteLen:]
					case rn == r.Comma:
						// `",` sequence (end of field).
						line = line[commaLen:]
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						continue parseField
					case lengthNLString(line) == len(line):
						// `"\n` sequence (end of line).
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						break parseField
					case r.LazyQuotes:
						// `"` sequence (bare quote).
						r.recordBuffer = append(r.recordBuffer, '"')
					default:
						// `"*` sequence (invalid non-escaped quote).
						err = ErrQuote
						break parseField
					}
				} else if len(line) > 0 {
					// Hit end of line (copy all data so far).
					r.recordBuffer = append(r.recordBuffer, line...)
					r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
					break parseField
				} else {
					// Abrupt end of file (EOF or error).
					if !r.LazyQuotes {
						err = ErrFieldCount
						break parseField
					}
					r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
					break parseField
				}
			}
		}
	}

	// Create a single string and create slices out of it.
	// This pins the memory of the fields together, but allocates once.
	str := string(r.recordBuffer) // Convert to string once to batch allocations
	r.lastRecord = r.lastRecord[:0]
	if cap(r.lastRecord) < len(r.fieldIndexes) {
		r.lastRecord = make([]string, len(r.fieldIndexes))
	}
	r.lastRecord = r.lastRecord[:len(r.fieldIndexes)]
	var preIdx int
	for i, idx := range r.fieldIndexes {
		r.lastRecord[i] = str[preIdx:idx]
		preIdx = idx
	}

	// Check or update the expected fields per record.
	if r.FieldsPerRecord > 0 {
		if len(r.lastRecord) != r.FieldsPerRecord && err == nil {
			err = ErrFieldCount
		}
	} else if r.FieldsPerRecord == 0 {
		r.FieldsPerRecord = len(r.lastRecord)
	}
	return r.lastRecord, err
}

func (r *Parser) Parse(line []byte) ([]string, error) {
	if r.Comment != 0 && nextRune(line) == r.Comment {
		return nil, nil // Skip Comment lines
	}
	if len(line) == lengthNL(line) {
		return nil, nil // Skip empty lines
	}

	var err error
	commaLen := utf8.RuneLen(r.Comma)
	r.recordBuffer = r.recordBuffer[:0]
	r.fieldIndexes = r.fieldIndexes[:0]

parseField:
	for {
		if r.TrimLeadingSpace {
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
		}
		if len(line) == 0 || line[0] != '"' {
			// Non-quoted string field
			i := bytes.IndexRune(line, r.Comma)
			field := line
			if i >= 0 {
				field = field[:i]
			} else {
				field = field[:len(field)-lengthNL(field)]
			}
			// Check to make sure a quote does not appear in field.
			if !r.LazyQuotes {
				if j := bytes.IndexByte(field, '"'); j >= 0 {
					err = ErrBareQuote
					break parseField
				}
			}
			r.recordBuffer = append(r.recordBuffer, field...)
			r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
			if i >= 0 {
				line = line[i+commaLen:]
				continue parseField
			}
			break parseField
		} else {
			// Quoted string field
			line = line[quoteLen:]
			for {
				i := bytes.IndexByte(line, '"')
				if i >= 0 {
					// Hit next quote.
					r.recordBuffer = append(r.recordBuffer, line[:i]...)
					line = line[i+quoteLen:]
					switch rn := nextRune(line); {
					case rn == '"':
						// `""` sequence (append quote).
						r.recordBuffer = append(r.recordBuffer, '"')
						line = line[quoteLen:]
					case rn == r.Comma:
						// `",` sequence (end of field).
						line = line[commaLen:]
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						continue parseField
					case lengthNL(line) == len(line):
						// `"\n` sequence (end of line).
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						break parseField
					case r.LazyQuotes:
						// `"` sequence (bare quote).
						r.recordBuffer = append(r.recordBuffer, '"')
					default:
						// `"*` sequence (invalid non-escaped quote).
						err = ErrQuote
						break parseField
					}
				} else if len(line) > 0 {
					// Hit end of line (copy all data so far).
					r.recordBuffer = append(r.recordBuffer, line...)
					r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
					break parseField
				} else {
					// Abrupt end of file (EOF or error).
					if !r.LazyQuotes {
						err = ErrFieldCount
						break parseField
					}
					r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
					break parseField
				}
			}
		}
	}

	// Create a single string and create slices out of it.
	// This pins the memory of the fields together, but allocates once.
	str := string(r.recordBuffer) // Convert to string once to batch allocations
	r.lastRecord = r.lastRecord[:0]
	if cap(r.lastRecord) < len(r.fieldIndexes) {
		r.lastRecord = make([]string, len(r.fieldIndexes))
	}
	r.lastRecord = r.lastRecord[:len(r.fieldIndexes)]
	var preIdx int
	for i, idx := range r.fieldIndexes {
		r.lastRecord[i] = str[preIdx:idx]
		preIdx = idx
	}

	// Check or update the expected fields per record.
	if r.FieldsPerRecord > 0 {
		if len(r.lastRecord) != r.FieldsPerRecord && err == nil {
			err = ErrFieldCount
		}
	} else if r.FieldsPerRecord == 0 {
		r.FieldsPerRecord = len(r.lastRecord)
	}
	return r.lastRecord, err
}

// lengthNLString reports the number of bytes for the trailing \n.
func lengthNLString(s string) int {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return 1
	}
	return 0
}

// nextRuneString returns the next rune in s or utf8.RuneError.
func nextRunesString(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

// lengthNL reports the number of bytes for the trailing \n.
func lengthNL(b []byte) int {
	if len(b) > 0 && b[len(b)-1] == '\n' {
		return 1
	}
	return 0
}

// nextRune returns the next rune in b or utf8.RuneError.
func nextRune(b []byte) rune {
	r, _ := utf8.DecodeRune(b)
	return r
}
