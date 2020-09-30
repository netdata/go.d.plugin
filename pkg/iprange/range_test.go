package iprange

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV4Range_String(t *testing.T) {
	tests := map[string]struct {
		input      string
		wantString string
	}{
		"IP":    {input: "192.0.2.0", wantString: "192.0.2.0-192.0.2.0"},
		"Range": {input: "192.0.2.0-192.0.2.10", wantString: "192.0.2.0-192.0.2.10"},
		"CIDR":  {input: "192.0.2.0/24", wantString: "192.0.2.1-192.0.2.254"},
		"Mask":  {input: "192.0.2.0/255.255.255.0", wantString: "192.0.2.1-192.0.2.254"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)
			require.NoError(t, err)

			assert.Equal(t, test.wantString, r.String())
		})
	}
}

func TestV4Range_Family(t *testing.T) {
	tests := map[string]struct {
		input string
	}{
		"IP":    {input: "192.0.2.0"},
		"Range": {input: "192.0.2.0-192.0.2.10"},
		"CIDR":  {input: "192.0.2.0/24"},
		"Mask":  {input: "192.0.2.0/255.255.255.0"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)
			require.NoError(t, err)

			assert.Equal(t, V4Family, r.Family())
		})
	}
}

func TestV4Range_Size(t *testing.T) {
	tests := map[string]struct {
		input    string
		wantSize *big.Int
	}{
		"IP":      {input: "192.0.2.0", wantSize: big.NewInt(1)},
		"Range":   {input: "192.0.2.0-192.0.2.10", wantSize: big.NewInt(11)},
		"CIDR":    {input: "192.0.2.0/24", wantSize: big.NewInt(254)},
		"CIDR 31": {input: "192.0.2.0/31", wantSize: big.NewInt(2)},
		"CIDR 32": {input: "192.0.2.0/32", wantSize: big.NewInt(1)},
		"Mask":    {input: "192.0.2.0/255.255.255.0", wantSize: big.NewInt(254)},
		"Mask 31": {input: "192.0.2.0/255.255.255.254", wantSize: big.NewInt(2)},
		"Mask 32": {input: "192.0.2.0/255.255.255.255", wantSize: big.NewInt(1)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)
			require.NoError(t, err)

			assert.Equal(t, test.wantSize, r.Size())
		})
	}
}

func TestV4Range_Contains(t *testing.T) {

}

func TestV6Range_String(t *testing.T) {
	tests := map[string]struct {
		input      string
		wantString string
	}{
		"IP":    {input: "2001:db8::", wantString: "2001:db8::-2001:db8::"},
		"Range": {input: "2001:db8::-2001:db8::10", wantString: "2001:db8::-2001:db8::10"},
		"CIDR":  {input: "2001:db8::/126", wantString: "2001:db8::1-2001:db8::2"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)
			require.NoError(t, err)

			assert.Equal(t, test.wantString, r.String())
		})
	}
}

func TestV6Range_Family(t *testing.T) {
	tests := map[string]struct {
		input string
	}{
		"IP":    {input: "2001:db8::"},
		"Range": {input: "2001:db8::-2001:db8::10"},
		"CIDR":  {input: "2001:db8::/126"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)
			require.NoError(t, err)

			assert.Equal(t, V6Family, r.Family())
		})
	}
}

func TestV6Range_Size(t *testing.T) {
	tests := map[string]struct {
		input    string
		wantSize *big.Int
	}{
		"IP":       {input: "2001:db8::", wantSize: big.NewInt(1)},
		"Range":    {input: "2001:db8::-2001:db8::10", wantSize: big.NewInt(17)},
		"CIDR":     {input: "2001:db8::/120", wantSize: big.NewInt(254)},
		"CIDR 127": {input: "2001:db8::/127", wantSize: big.NewInt(2)},
		"CIDR 128": {input: "2001:db8::/128", wantSize: big.NewInt(1)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := ParseRange(test.input)
			require.NoError(t, err)

			assert.Equal(t, test.wantSize, r.Size())
		})
	}
}

func TestV6Range_Contains(t *testing.T) {

}

//func mustNewRange(start, end string) *Range {
//	r := &Range{}
//	if start != "" {
//		r.Start = net.ParseIP(start)
//	}
//	if end != "" {
//		r.End = net.ParseIP(end)
//	}
//	return r
//}
//
//func TestRange_String(t *testing.T) {
//	expected := "1.1.1.1-1.1.1.1"
//	r := mustNewRange("1.1.1.1", "1.1.1.1")
//
//	assert.Equal(t, expected, r.String())
//}
//
//func TestRange_Family(t *testing.T) {
//	v4 := "1.1.1.1"
//	v6 := "1234::1"
//	cases := []struct {
//		Expected Family
//		Range    *Range
//	}{
//		{
//			Expected: InvalidFamily,
//			Range:    mustNewRange("", ""),
//		},
//		{
//			Expected: InvalidFamily,
//			Range:    mustNewRange(v4, ""),
//		},
//		{
//			Expected: InvalidFamily,
//			Range:    mustNewRange("", v4),
//		},
//		{
//			Expected: InvalidFamily,
//			Range:    mustNewRange(v4, v6),
//		},
//		{
//			Expected: V4Family,
//			Range:    mustNewRange(v4, v4),
//		},
//		{
//			Expected: V6Family,
//			Range:    mustNewRange(v6, v6),
//		},
//	}
//
//	for _, c := range cases {
//		assert.Equal(t, c.Expected, c.Range.Family())
//	}
//}
//
//func TestRange_Contains(t *testing.T) {
//	v4 := mustNewRange("1.1.1.1", "1.1.2.255")
//	v6 := mustNewRange("1234::1", "1235::1")
//	cases := []struct {
//		Expected bool
//		Range    *Range
//		IP       net.IP
//	}{
//		{
//			Expected: false,
//			Range:    v4,
//			IP:       net.ParseIP("2.2.2.1"),
//		},
//		{
//			Expected: false,
//			Range:    v6,
//			IP:       net.ParseIP("1236::1"),
//		},
//		{
//			Expected: true,
//			Range:    v4,
//			IP:       net.ParseIP("1.1.2.200"),
//		},
//		{
//			Expected: true,
//			Range:    v6,
//			IP:       net.ParseIP("1234::4:d:1"),
//		},
//	}
//
//	for _, c := range cases {
//		assert.Equal(t, c.Expected, c.Range.Contains(c.IP))
//	}
//
//}
//
//func Test_isRangeValid(t *testing.T) {
//	v4Start, v4End := "1.1.1.1", "1.1.2.255"
//	v6Start, v6End := "1234::1", "1234::f"
//	cases := []struct {
//		Expected bool
//		Range    *Range
//	}{
//		{
//			Expected: false,
//			Range:    mustNewRange("", ""),
//		},
//		{
//			Expected: false,
//			Range:    mustNewRange(v4Start, ""),
//		},
//		{
//			Expected: false,
//			Range:    mustNewRange("", v4Start),
//		},
//		{
//			Expected: false,
//			Range:    mustNewRange(v4Start, v6End),
//		},
//		{
//			Expected: false,
//			Range:    mustNewRange(v4End, v4Start),
//		},
//		{
//			Expected: false,
//			Range:    mustNewRange(v6End, v6Start),
//		},
//		{
//			Expected: true,
//			Range:    mustNewRange(v4Start, v4End),
//		},
//		{
//			Expected: true,
//			Range:    mustNewRange(v6Start, v6End),
//		},
//	}
//
//	for _, c := range cases {
//		assert.Equal(t, c.Expected, isRangeValid(*c.Range))
//	}
//}
//
//func TestParseRange(t *testing.T) {
//	assert.Nil(t, ParseRange("not ip address"))
//	assert.NotNil(t, ParseRange("1234::1, 1234::2"))
//	assert.NotNil(t, ParseRange("1.1.1.1"))
//	assert.NotNil(t, ParseRange("1.1.1.1-1.1.1.2"))
//	assert.NotNil(t, ParseRange("1234::1-1234::2"))
//}
//
//func TestRange_Hosts(t *testing.T) {
//	cases := []struct {
//		Expected float64
//		Range    *Range
//	}{
//		{
//			Expected: 1,
//			Range:    mustNewRange("1.1.1.0", "1.1.1.0"),
//		},
//		{
//			Expected: math.Pow(2, 8),
//			Range:    mustNewRange("1.1.1.0", "1.1.1.255"),
//		},
//		{
//			Expected: math.Pow(2, 16),
//			Range:    mustNewRange("1.1.0.0", "1.1.255.255"),
//		},
//		{
//			Expected: math.Pow(2, 32),
//			Range:    mustNewRange("0.0.0.0", "255.255.255.255"),
//		},
//		{
//			Expected: math.Pow(2, 16),
//			Range:    mustNewRange("1234::ffff:ffff:ffff:0", "1234::ffff:ffff:ffff:ffff"),
//		},
//	}
//
//	for _, c := range cases {
//		assert.Equal(t, int64(c.Expected), c.Range.Hosts().Int64())
//	}
//
//	assert.Nil(t, Range{}.Hosts())
//}
//
//func Test_v4RangeSize(t *testing.T) {
//	v4 := mustNewRange("1.1.0.0", "1.1.255.255")
//	expected := int64(math.Pow(2, 16))
//
//	assert.Equal(t, expected, v4RangeSize(*v4).Int64())
//}
//
//func Test_v6RangeSize(t *testing.T) {
//	v6 := mustNewRange("1234::ffff:ffff:ffff:0", "1234::ffff:ffff:ffff:ffff")
//	expected := int64(math.Pow(2, 16))
//
//	assert.Equal(t, expected, v6RangeSize(*v6).Int64())
//}
//
//func Test_v4ToInt(t *testing.T) {
//	cases := []struct {
//		Expected int64
//		IP       net.IP
//	}{
//		{
//			Expected: 1,
//			IP:       net.ParseIP("0.0.0.1"),
//		},
//		{
//			Expected: 256,
//			IP:       net.ParseIP("0.0.1.0"),
//		},
//		{
//			Expected: 65536,
//			IP:       net.ParseIP("0.1.0.0"),
//		},
//		{
//			Expected: 16777216,
//			IP:       net.ParseIP("1.0.0.0"),
//		},
//	}
//
//	for _, c := range cases {
//		assert.Equal(t, c.Expected, v4ToInt(c.IP))
//	}
//}
