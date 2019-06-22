package ip

import (
	"math"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustNewRange(start, end string) *Range {
	r := &Range{}
	if start != "" {
		r.Start = net.ParseIP(start)
	}
	if end != "" {
		r.End = net.ParseIP(end)
	}
	return r
}

func TestRange_String(t *testing.T) {
	expected := "1.1.1.1-1.1.1.1"
	r := mustNewRange("1.1.1.1", "1.1.1.1")

	assert.Equal(t, expected, r.String())
}

func TestRange_Family(t *testing.T) {
	v4 := "1.1.1.1"
	v6 := "1234::1"
	cases := []struct {
		Expected Family
		Range    *Range
	}{
		{
			Expected: InvalidFamily,
			Range:    mustNewRange("", ""),
		},
		{
			Expected: InvalidFamily,
			Range:    mustNewRange(v4, ""),
		},
		{
			Expected: InvalidFamily,
			Range:    mustNewRange("", v4),
		},
		{
			Expected: InvalidFamily,
			Range:    mustNewRange(v4, v6),
		},
		{
			Expected: V4Family,
			Range:    mustNewRange(v4, v4),
		},
		{
			Expected: V6Family,
			Range:    mustNewRange(v6, v6),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, c.Range.Family())
	}
}

func TestRange_Contains(t *testing.T) {
	v4 := mustNewRange("1.1.1.1", "1.1.2.255")
	v6 := mustNewRange("1234::1", "1235::1")
	cases := []struct {
		Expected bool
		Range    *Range
		IP       net.IP
	}{
		{
			Expected: false,
			Range:    v4,
			IP:       net.ParseIP("2.2.2.1"),
		},
		{
			Expected: false,
			Range:    v6,
			IP:       net.ParseIP("1236::1"),
		},
		{
			Expected: true,
			Range:    v4,
			IP:       net.ParseIP("1.1.2.200"),
		},
		{
			Expected: true,
			Range:    v6,
			IP:       net.ParseIP("1234::4:d:1"),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, c.Range.Contains(c.IP))
	}

}

func Test_isRangeValid(t *testing.T) {
	v4Start, v4End := "1.1.1.1", "1.1.2.255"
	v6Start, v6End := "1234::1", "1234::f"
	cases := []struct {
		Expected bool
		Range    *Range
	}{
		{
			Expected: false,
			Range:    mustNewRange("", ""),
		},
		{
			Expected: false,
			Range:    mustNewRange(v4Start, ""),
		},
		{
			Expected: false,
			Range:    mustNewRange("", v4Start),
		},
		{
			Expected: false,
			Range:    mustNewRange(v4Start, v6End),
		},
		{
			Expected: false,
			Range:    mustNewRange(v4End, v4Start),
		},
		{
			Expected: false,
			Range:    mustNewRange(v6End, v6Start),
		},
		{
			Expected: true,
			Range:    mustNewRange(v4Start, v4End),
		},
		{
			Expected: true,
			Range:    mustNewRange(v6Start, v6End),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, isRangeValid(*c.Range))
	}
}

func TestParseRange(t *testing.T) {
	assert.Nil(t, ParseRange("not ip address"))
	assert.NotNil(t, ParseRange("1234::1, 1234::2"))
	assert.NotNil(t, ParseRange("1.1.1.1"))
	assert.NotNil(t, ParseRange("1.1.1.1-1.1.1.2"))
	assert.NotNil(t, ParseRange("1234::1-1234::2"))
}

func TestRange_Hosts(t *testing.T) {
	cases := []struct {
		Expected float64
		Range    *Range
	}{
		{
			Expected: 1,
			Range:    mustNewRange("1.1.1.0", "1.1.1.0"),
		},
		{
			Expected: math.Pow(2, 8),
			Range:    mustNewRange("1.1.1.0", "1.1.1.255"),
		},
		{
			Expected: math.Pow(2, 16),
			Range:    mustNewRange("1.1.0.0", "1.1.255.255"),
		},
		{
			Expected: math.Pow(2, 32),
			Range:    mustNewRange("0.0.0.0", "255.255.255.255"),
		},
		{
			Expected: math.Pow(2, 16),
			Range:    mustNewRange("1234::ffff:ffff:ffff:0", "1234::ffff:ffff:ffff:ffff"),
		},
	}

	for _, c := range cases {
		assert.Equal(t, int64(c.Expected), c.Range.Hosts().Int64())
	}

	assert.Nil(t, Range{}.Hosts())
}

func Test_v4RangeSize(t *testing.T) {
	v4 := mustNewRange("1.1.0.0", "1.1.255.255")
	expected := int64(math.Pow(2, 16))

	assert.Equal(t, expected, v4RangeSize(*v4).Int64())
}

func Test_v6RangeSize(t *testing.T) {
	v6 := mustNewRange("1234::ffff:ffff:ffff:0", "1234::ffff:ffff:ffff:ffff")
	expected := int64(math.Pow(2, 16))

	assert.Equal(t, expected, v6RangeSize(*v6).Int64())
}

func Test_v4ToInt(t *testing.T) {
	cases := []struct {
		Expected int64
		IP       net.IP
	}{
		{
			Expected: 1,
			IP:       net.ParseIP("0.0.0.1"),
		},
		{
			Expected: 256,
			IP:       net.ParseIP("0.0.1.0"),
		},
		{
			Expected: 65536,
			IP:       net.ParseIP("0.1.0.0"),
		},
		{
			Expected: 16777216,
			IP:       net.ParseIP("1.0.0.0"),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, v4ToInt(c.IP))
	}
}
