package weblog

import (
	"bytes"
	"errors"
	"unicode"
	"unicode/utf8"
)

type csvReader struct {
	lazyQuotes       bool
	trimLeadingSpace bool
	comma            rune

	recordBuffer []byte
	fieldIndexes []int
	lastRecord   []string
}

func (r *csvReader) readRecord(line []byte) ([]string, error) {
	// Parse each field in the record.
	var err error
	const quoteLen = len(`"`)
	commaLen := utf8.RuneLen(r.comma)

	r.recordBuffer = r.recordBuffer[:0]
	r.fieldIndexes = r.fieldIndexes[:0]
parseField:
	for {
		if r.trimLeadingSpace {
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
		}
		if len(line) == 0 || line[0] != '"' {
			// Non-quoted string field
			i := bytes.IndexRune(line, r.comma)
			field := line
			if i >= 0 {
				field = field[:i]
			} else {
				field = field[:len(field)-lengthNL(field)]
			}
			// Check to make sure a quote does not appear in field.
			if !r.lazyQuotes {
				if j := bytes.IndexByte(field, '"'); j >= 0 {
					err = errors.New("bare \" in non-quoted-field")
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
					case rn == r.comma:
						// `",` sequence (end of field).
						line = line[commaLen:]
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						continue parseField
					case lengthNL(line) == len(line):
						// `"\n` sequence (end of line).
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						break parseField
					case r.lazyQuotes:
						// `"` sequence (bare quote).
						r.recordBuffer = append(r.recordBuffer, '"')
					default:
						// `"*` sequence (invalid non-escaped quote).
						err = errors.New("extraneous or missing \" in quoted-field")
						break parseField
					}
					// TODO: do we need this case?
					//} else if len(line) > 0 {
				} else {
					if !r.lazyQuotes {
						err = errors.New("extraneous or missing \" in quoted-field")
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

	return r.lastRecord, err
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
