package logs

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegExpParser(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{name: "valid pattern", pattern: `(?P<A>\d+) (?P<B>\d+)`},
		{name: "no names subgroups in pattern", pattern: `(?:\d+) (?:\d+)`, wantErr: true},
		{name: "invalid pattern", pattern: `(((?P<A>\d+) (?P<B>\d+)`, wantErr: true},
		{name: "empty pattern", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewRegExpParser(RegExpConfig{Pattern: tt.pattern}, nil)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p)
			}
		})
	}
}

func TestRegExpParser_Info(t *testing.T) {
	pattern := `(?P<A>\d+) (?P<B>\d+)`
	p, err := NewRegExpParser(RegExpConfig{Pattern: pattern}, nil)
	require.NoError(t, err)
	expectedInfo := fmt.Sprintf("regexp: %s", pattern)
	assert.Equal(t, expectedInfo, p.Info())
}

func TestRegExpParser_ReadLine(t *testing.T) {
	tests := []struct {
		name         string
		row          string
		pattern      string
		wantErr      bool
		wantParseErr bool
	}{
		{name: "match no error", row: "1 2", pattern: `(?P<A>\d+) (?P<B>\d+)`},
		{name: "match error on assigning", row: "1 2", pattern: `(?P<AA>\d+) (?P<BB>\d+)`, wantErr: true, wantParseErr: true},
		{name: "no match", row: "A B", pattern: `(?P<A>\d+) (?P<B>\d+)`, wantErr: true, wantParseErr: true},
		{name: "error on reading EOF", row: "", pattern: `(?P<A>\d+) (?P<B>\d+)`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var line logLine
			r := strings.NewReader(tt.row)
			p, err := NewRegExpParser(RegExpConfig{Pattern: tt.pattern}, r)
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

func TestRegExpParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		row     string
		pattern string
		wantErr bool
	}{
		{name: "match no error", row: "1 2", pattern: `(?P<A>\d+) (?P<B>\d+)`},
		{name: "match error on assigning", row: "1 2", pattern: `(?P<AA>\d+) (?P<BB>\d+)`, wantErr: true},
		{name: "no match", row: "A B", pattern: `(?P<A>\d+) (?P<B>\d+)`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var line logLine
			p, err := NewRegExpParser(RegExpConfig{Pattern: tt.pattern}, nil)
			require.NoError(t, err)

			err = p.Parse([]byte(tt.row), &line)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, IsParseError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type logLine struct{}

func (l logLine) Assign(name, val string) error {
	switch name {
	case "A", "B":
	default:
		return errors.New("unknown var name")
	}
	return nil
}
