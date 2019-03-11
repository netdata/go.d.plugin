package weblog

import (
	"encoding/csv"
	"io"
)

func NewLogParser(r io.Reader) *csv.Reader {
	parser := csv.NewReader(r)
	parser.Comma = ' '
	parser.ReuseRecord = true
	parser.TrimLeadingSpace = true
	parser.FieldsPerRecord = -1
	return parser
}
