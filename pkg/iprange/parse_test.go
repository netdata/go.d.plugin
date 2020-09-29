package iprange

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRanges(t *testing.T) {

}

func TestParseRange(t *testing.T) {
	tests := map[string]struct {
		input     string
		wantRange Range
		wantErr   bool
	}{
		"v4 IP": {
			input:     "1.2.3.0",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.0"),
		},
		"v4 IP: invalid address": {
			input:   "1.2.3.",
			wantErr: true,
		},

		"v4 Range": {
			input:     "1.2.3.0-1.2.3.10",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.10"),
		},
		"v4 Range: start == end": {
			input:     "1.2.3.0-1.2.3.0",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.0"),
		},
		"v4 Range: start > end": {
			input:   "1.2.3.10-1.2.3.0",
			wantErr: true,
		},
		"v4 Range: invalid start": {
			input:   "1.2.3.-1.2.3.10",
			wantErr: true,
		},
		"v4 Range: invalid end": {
			input:   "1.2.3.0-1.2.3.",
			wantErr: true,
		},
		"v4 Range: v6 start": {
			input:   "::1-1.2.3.10",
			wantErr: true,
		},
		"v4 Range: v6 end": {
			input:   "1.2.3.10-::1",
			wantErr: true,
		},

		"v4 CIDR: /0": {
			input:     "1.2.3.0/0",
			wantRange: prepareV4Range("0.0.0.1", "255.255.255.254"),
		},
		"v4 CIDR: /24": {
			input:     "1.2.3.0/24",
			wantRange: prepareV4Range("1.2.3.1", "1.2.3.254"),
		},
		"v4 CIDR: /30": {
			input:     "1.2.3.0/30",
			wantRange: prepareV4Range("1.2.3.1", "1.2.3.2"),
		},
		"v4 CIDR: /31": {
			input:     "1.2.3.0/31",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.1"),
		},
		"v4 CIDR: /32": {
			input:     "1.2.3.0/32",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.0"),
		},
		"v4 CIDR: missing prefix length": {
			input:   "1.2.3.0/",
			wantErr: true,
		},
		"v4 CIDR: invalid prefix length": {
			input:   "1.2.3.0/99",
			wantErr: true,
		},

		"v4 Net: /0": {
			input:     "1.2.3.0/0.0.0.0",
			wantRange: prepareV4Range("0.0.0.1", "255.255.255.254"),
		},
		"v4 Net: /24": {
			input:     "1.2.3.0/255.255.255.0",
			wantRange: prepareV4Range("1.2.3.1", "1.2.3.254"),
		},
		"v4 Net: /30": {
			input:     "1.2.3.0/255.255.255.252",
			wantRange: prepareV4Range("1.2.3.1", "1.2.3.2"),
		},
		"v4 Net: /31": {
			input:     "1.2.3.0/255.255.255.254",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.1"),
		},
		"v4 Net: /32": {
			input:     "1.2.3.0/255.255.255.255",
			wantRange: prepareV4Range("1.2.3.0", "1.2.3.0"),
		},
		"v4 Net: missing prefix mask": {
			input:   "1.2.3.0/",
			wantErr: true,
		},
		"v4 Net: invalid mask": {
			input:   "1.2.3.0/mask",
			wantErr: true,
		},
		"v4 Net: not canonical form mask": {
			input:   "1.2.3.0/255.255.0.254",
			wantErr: true,
		},
		"v4 Net: v6 address": {
			input:   "::1/255.255.255.0",
			wantErr: true,
		},
	}

	for name, test := range tests {
		name = fmt.Sprintf("%s (%s)", name, test.input)
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)

			if test.wantErr {
				assert.Error(t, err)
				assert.Nilf(t, r, "want: nil, got: %s", r)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, test.wantRange, r, "want: %s, got: %s", test.wantRange, r)
			}
		})
	}
}

func prepareV4Range(start, end string) Range {
	return v4Range{
		start: net.ParseIP(start),
		end:   net.ParseIP(end),
	}
}

func prepareV6Range(start, end string) Range {
	return v6Range{
		start: net.ParseIP(start),
		end:   net.ParseIP(end),
	}
}
