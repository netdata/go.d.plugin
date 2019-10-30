package logs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCSVConfig = CSVConfig{
	Delimiter:        ' ',
	TrimLeadingSpace: false,
	Format:           "$A $B !C $D",
	CheckField:       checkCSVFormatField,
}

func TestNewCSVParser(t *testing.T) {

}

func TestCSVParser_ReadLine(t *testing.T) {
	tests := []struct {
		name         string
		row          string
		format       string
		wantErr      bool
		wantParseErr bool
	}{
		{name: "match and no error", row: "1 2 3", format: `$A $B $C`},
		{name: "match but error on assigning", row: "1 2 3", format: `$A $B $ERR`, wantErr: true, wantParseErr: true},
		{name: "not match", row: "1 2 3", format: `$A $B $C $d`, wantErr: true, wantParseErr: true},
		{name: "error on reading csv.Err", row: "1 2\"3", format: `$A $B $C`, wantErr: true, wantParseErr: true},
		{name: "error on reading EOF", row: "", format: `$A $B $C`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var line logLine
			r := strings.NewReader(tt.row)
			c := testCSVConfig
			c.Format = tt.format
			p, err := NewCSVParser(c, r)
			require.NoError(t, err)

			err = p.ReadLine(&line)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantParseErr {
					assert.True(t, IsParseError(err))
				} else {
					assert.False(t, IsParseError(err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCSVParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		row     string
		format  string
		wantErr bool
	}{
		{name: "match and no error", row: "1 2 3", format: `$A $B $C`},
		{name: "match but error on assigning", row: "1 2 3", format: `$A $B $ERR`, wantErr: true},
		{name: "not match", row: "1 2 3", format: `$A $B $C $d`, wantErr: true},
		{name: "error on reading csv.Err", row: "1 2\"3", format: `$A $B $C`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var line logLine
			r := strings.NewReader(tt.row)
			c := testCSVConfig
			c.Format = tt.format
			p, err := NewCSVParser(c, r)
			require.NoError(t, err)

			err = p.ReadLine(&line)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, IsParseError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestCSVParser_Info(t *testing.T) {
	p, err := NewCSVParser(testCSVConfig, nil)
	require.NoError(t, err)
	assert.NotZero(t, p.Info())
}
