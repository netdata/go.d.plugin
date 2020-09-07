package logs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testJSONConfig = JSONConfig{
	Mapping: map[string]string{"from_field_1": "to_field_1"},
}

func TestNewJSONParser(t *testing.T) {
	tests := []struct {
		name    string
		config  JSONConfig
		wantErr bool
	}{
		{
			name:    "empty config",
			config:  JSONConfig{},
			wantErr: false,
		},
		{
			name:    "basic config",
			config:  testJSONConfig,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewJSONParser(tt.config, nil)

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

func TestJSONParser_ReadLine(t *testing.T) {
	tests := []struct {
		name    string
		config  JSONConfig
		data    string
		wantErr bool
	}{
		{
			name:    "no error",
			config:  JSONConfig{},
			wantErr: false,
			data:    `{ "host": "example.com" }`,
		},
		{
			name:    "splits on newline",
			config:  JSONConfig{},
			wantErr: false,
			data:    "{\"host\": \"example.com\"}\n{\"host\": \"acme.org\"}",
		},
		{
			name:    "error on malformed JSON",
			config:  JSONConfig{},
			wantErr: true,
			data:    `{ "host"": unquoted_string}`,
		},
		{
			name:    "error on no data",
			config:  JSONConfig{},
			wantErr: true,
			data:    ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var line logLine
			in := strings.NewReader(tt.data)
			p, err := NewJSONParser(tt.config, in)
			require.NoError(t, err)
			require.NotNil(t, p)

			err = p.ReadLine(&line)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJSONParser_Parse(t *testing.T) {
	tests := []struct {
		name         string
		row          string
		fieldMap     map[string]string
		wantParseErr bool
	}{
		{
			name:         "malformed JSON",
			row:          `{`,
			wantParseErr: true,
		},
		{
			name:         "malformed JSON #2",
			row:          `{ host: "example.com" }`,
			wantParseErr: true,
		},
		{
			name:         "empty string",
			row:          "",
			wantParseErr: true,
		},
		{
			name:         "no field mapping",
			row:          `{ "host": "example.com", "remote_addr": "127.0.0.1", "request_time": 0.05 }`,
			wantParseErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var line logLine
			p, err := NewJSONParser(JSONConfig{Mapping: tt.fieldMap}, nil)
			assert.NoError(t, err)

			parseErr := p.Parse([]byte(tt.row), line)

			if tt.wantParseErr {
				assert.Error(t, parseErr)
			} else {
				assert.NoError(t, parseErr)
			}

		})
	}
}
