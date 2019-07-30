package match

import (
	"testing"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/stretchr/testify/assert"
)

func TestOrHostMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		lhs      HostMatcher
		rhs      HostMatcher
	}{
		{
			name:     "true, true",
			expected: true,
			lhs:      hostDcMatcher{matcher.TRUE()},
			rhs:      hostDcMatcher{matcher.TRUE()},
		},
		{
			name:     "true, false",
			expected: true,
			lhs:      hostDcMatcher{matcher.TRUE()},
			rhs:      hostDcMatcher{matcher.FALSE()},
		},
		{
			name:     "false, true",
			expected: true,
			lhs:      hostDcMatcher{matcher.FALSE()},
			rhs:      hostDcMatcher{matcher.TRUE()},
		},
		{
			name:     "false, false",
			expected: false,
			lhs:      hostDcMatcher{matcher.FALSE()},
			rhs:      hostDcMatcher{matcher.FALSE()},
		},
	}

	host := &rs.Host{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, orHostMatcher{lhs: test.lhs, rhs: test.rhs}.Match(host))
		})
	}
}

func TestAndHostMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		lhs      HostMatcher
		rhs      HostMatcher
	}{
		{
			name:     "true, true",
			expected: true,
			lhs:      hostDcMatcher{matcher.TRUE()},
			rhs:      hostDcMatcher{matcher.TRUE()},
		},
		{
			name:     "true, false",
			expected: false,
			lhs:      hostDcMatcher{matcher.TRUE()},
			rhs:      hostDcMatcher{matcher.FALSE()},
		},
		{
			name:     "false, true",
			expected: false,
			lhs:      hostDcMatcher{matcher.FALSE()},
			rhs:      hostDcMatcher{matcher.TRUE()},
		},
		{
			name:     "false, false",
			expected: false,
			lhs:      hostDcMatcher{matcher.FALSE()},
			rhs:      hostDcMatcher{matcher.FALSE()},
		},
	}

	host := &rs.Host{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, andHostMatcher{lhs: test.lhs, rhs: test.rhs}.Match(host))
		})
	}
}

func TestOrVMMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		lhs      VMMatcher
		rhs      VMMatcher
	}{
		{
			name:     "true, true",
			expected: true,
			lhs:      vmDcMatcher{matcher.TRUE()},
			rhs:      vmDcMatcher{matcher.TRUE()},
		},
		{
			name:     "true, false",
			expected: true,
			lhs:      vmDcMatcher{matcher.TRUE()},
			rhs:      vmDcMatcher{matcher.FALSE()},
		},
		{
			name:     "false, true",
			expected: true,
			lhs:      vmDcMatcher{matcher.FALSE()},
			rhs:      vmDcMatcher{matcher.TRUE()},
		},
		{
			name:     "false, false",
			expected: false,
			lhs:      vmDcMatcher{matcher.FALSE()},
			rhs:      vmDcMatcher{matcher.FALSE()},
		},
	}

	vm := &rs.VM{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, orVMMatcher{lhs: test.lhs, rhs: test.rhs}.Match(vm))
		})
	}
}

func TestAndVMMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		lhs      VMMatcher
		rhs      VMMatcher
	}{
		{
			name:     "true, true",
			expected: true,
			lhs:      vmDcMatcher{matcher.TRUE()},
			rhs:      vmDcMatcher{matcher.TRUE()},
		},
		{
			name:     "true, false",
			expected: false,
			lhs:      vmDcMatcher{matcher.TRUE()},
			rhs:      vmDcMatcher{matcher.FALSE()},
		},
		{
			name:     "false, true",
			expected: false,
			lhs:      vmDcMatcher{matcher.FALSE()},
			rhs:      vmDcMatcher{matcher.TRUE()},
		},
		{
			name:     "false, false",
			expected: false,
			lhs:      vmDcMatcher{matcher.FALSE()},
			rhs:      vmDcMatcher{matcher.FALSE()},
		},
	}

	vm := &rs.VM{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, andVMMatcher{lhs: test.lhs, rhs: test.rhs}.Match(vm))
		})
	}
}

func TestHostIncludes_Parse(t *testing.T) {
	tests := []struct {
		name     string
		valid    bool
		includes HostIncludes
		expected HostMatcher
	}{
		{
			name:     "",
			includes: []string{""},
			valid:    false,
		},
		{
			name:     "*/C1/H1",
			includes: []string{"*/C1/H1"},
			valid:    false,
		},
		{
			name:     "/",
			includes: []string{"/"},
			valid:    true,
			expected: hostDcMatcher{matcher.FALSE()},
		},
		{
			name:     "/*",
			includes: []string{"/*"},
			valid:    true,
			expected: hostDcMatcher{matcher.TRUE()},
		},
		{
			name:     "/!*",
			includes: []string{"/!*"},
			valid:    true,
			expected: hostDcMatcher{matcher.FALSE()},
		},
		{
			name:     "/!*/",
			includes: []string{"/!*/"},
			valid:    true,
			expected: hostDcMatcher{matcher.FALSE()},
		},
		{
			name:     "/!*/ ",
			includes: []string{"/!*/ "},
			valid:    true,
			expected: andHostMatcher{
				lhs: hostDcMatcher{matcher.FALSE()},
				rhs: hostClusterMatcher{matcher.FALSE()},
			},
		},
		{
			name:     "/DC1* DC2* !*/Cluster*",
			includes: []string{"/DC1* DC2* !*/Cluster*"},
			valid:    true,
			expected: andHostMatcher{
				lhs: hostDcMatcher{mustSimplePattern("DC1* DC2* !*")},
				rhs: hostClusterMatcher{mustSimplePattern("Cluster*")},
			},
		},
		{
			name:     "/*/*/HOST1*",
			includes: []string{"/*/*/HOST1*"},
			valid:    true,
			expected: andHostMatcher{
				lhs: andHostMatcher{
					lhs: hostDcMatcher{matcher.TRUE()},
					rhs: hostClusterMatcher{matcher.TRUE()},
				},
				rhs: hostHostMatcher{mustSimplePattern("HOST1*")},
			},
		},
		{
			name:     "/*/*/HOST1*/*/*",
			includes: []string{"/*/*/HOST1*/*/*"},
			valid:    true,
			expected: andHostMatcher{
				lhs: andHostMatcher{
					lhs: hostDcMatcher{matcher.TRUE()},
					rhs: hostClusterMatcher{matcher.TRUE()},
				},
				rhs: hostHostMatcher{mustSimplePattern("HOST1*")},
			},
		},
		{
			name:     "[/DC1*, /DC2*]",
			includes: []string{"/DC1*", "/DC2*"},
			valid:    true,
			expected: orHostMatcher{
				lhs: hostDcMatcher{mustSimplePattern("DC1*")},
				rhs: hostDcMatcher{mustSimplePattern("DC2*")},
			},
		},
		{
			name:     "[/DC1*, /DC2*, /DC3/Cluster1*/H*]",
			includes: []string{"/DC1*", "/DC2*", "/DC3*/Cluster1*/H*"},
			valid:    true,
			expected: orHostMatcher{
				lhs: orHostMatcher{
					lhs: hostDcMatcher{mustSimplePattern("DC1*")},
					rhs: hostDcMatcher{mustSimplePattern("DC2*")},
				},
				rhs: andHostMatcher{
					lhs: andHostMatcher{
						lhs: hostDcMatcher{mustSimplePattern("DC3*")},
						rhs: hostClusterMatcher{mustSimplePattern("Cluster1*")},
					},
					rhs: hostHostMatcher{mustSimplePattern("H*")},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := test.includes.Parse()
			if !test.valid {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, test.expected, m)
		})
	}
}

func TestVMIncludes_Parse(t *testing.T) {
	tests := []struct {
		name     string
		valid    bool
		includes VMIncludes
		expected VMMatcher
	}{
		{
			name:     "",
			includes: []string{""},
			valid:    false,
		},
		{
			name:     "*/C1/H1/V1",
			includes: []string{"*/C1/H1/V1"},
			valid:    false,
		},
		{
			name:     "/*",
			includes: []string{"/*"},
			valid:    true,
			expected: vmDcMatcher{matcher.TRUE()},
		},
		{
			name:     "/!*",
			includes: []string{"/!*"},
			valid:    true,
			expected: vmDcMatcher{matcher.FALSE()},
		},
		{
			name:     "/!*/",
			includes: []string{"/!*/"},
			valid:    true,
			expected: vmDcMatcher{matcher.FALSE()},
		},
		{
			name:     "/!*/ ",
			includes: []string{"/!*/ "},
			valid:    true,
			expected: andVMMatcher{
				lhs: vmDcMatcher{matcher.FALSE()},
				rhs: vmClusterMatcher{matcher.FALSE()},
			},
		},
		{
			name:     "/DC1* DC2* !*/Cluster*",
			includes: []string{"/DC1* DC2* !*/Cluster*"},
			valid:    true,
			expected: andVMMatcher{
				lhs: vmDcMatcher{mustSimplePattern("DC1* DC2* !*")},
				rhs: vmClusterMatcher{mustSimplePattern("Cluster*")},
			},
		},
		{
			name:     "/*/*/HOST1",
			includes: []string{"/*/*/HOST1"},
			valid:    true,
			expected: andVMMatcher{
				lhs: andVMMatcher{
					lhs: vmDcMatcher{matcher.TRUE()},
					rhs: vmClusterMatcher{matcher.TRUE()},
				},
				rhs: vmHostMatcher{mustSimplePattern("HOST1")},
			},
		},
		{
			name:     "/*/*/HOST1*/*/*",
			includes: []string{"/*/*/HOST1*/*/*"},
			valid:    true,
			expected: andVMMatcher{
				lhs: andVMMatcher{
					lhs: andVMMatcher{
						lhs: vmDcMatcher{matcher.TRUE()},
						rhs: vmClusterMatcher{matcher.TRUE()},
					},
					rhs: vmHostMatcher{mustSimplePattern("HOST1*")},
				},
				rhs: vmVMMatcher{matcher.TRUE()},
			},
		},
		{
			name:     "[/DC1*, /DC2*]",
			includes: []string{"/DC1*", "/DC2*"},
			valid:    true,
			expected: orVMMatcher{
				lhs: vmDcMatcher{mustSimplePattern("DC1*")},
				rhs: vmDcMatcher{mustSimplePattern("DC2*")},
			},
		},
		{
			name:     "[/DC1*, /DC2*, /DC3*/Cluster1*/H*/VM*]",
			includes: []string{"/DC1*", "/DC2*", "/DC3*/Cluster1*/H*/VM*"},
			valid:    true,
			expected: orVMMatcher{
				lhs: orVMMatcher{
					lhs: vmDcMatcher{mustSimplePattern("DC1*")},
					rhs: vmDcMatcher{mustSimplePattern("DC2*")},
				},
				rhs: andVMMatcher{
					lhs: andVMMatcher{
						lhs: andVMMatcher{
							lhs: vmDcMatcher{mustSimplePattern("DC3*")},
							rhs: vmClusterMatcher{mustSimplePattern("Cluster1*")},
						},
						rhs: vmHostMatcher{mustSimplePattern("H*")},
					},
					rhs: vmVMMatcher{mustSimplePattern("VM*")},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := test.includes.Parse()
			if !test.valid {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, test.expected, m)
		})
	}
}

func mustSimplePattern(expr string) matcher.Matcher {
	return matcher.Must(matcher.NewSimplePatternsMatcher(expr))
}
