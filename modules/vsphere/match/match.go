package match

import (
	"strings"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type HostMatcher interface {
	Match(*rs.Host) bool
}

type VMMatcher interface {
	Match(*rs.VM) bool
}

type (
	hostDcMatcher      struct{ m matcher.Matcher }
	hostClusterMatcher struct{ m matcher.Matcher }
	hostHostMatcher    struct{ m matcher.Matcher }
	vmDcMatcher        struct{ m matcher.Matcher }
	vmClusterMatcher   struct{ m matcher.Matcher }
	vmHostMatcher      struct{ m matcher.Matcher }
	vmVMMatcher        struct{ m matcher.Matcher }
	orHostMatcher      struct{ lhs, rhs HostMatcher }
	orVMMatcher        struct{ lhs, rhs VMMatcher }
	andHostMatcher     struct{ lhs, rhs HostMatcher }
	andVMMatcher       struct{ lhs, rhs VMMatcher }
)

func (m hostDcMatcher) Match(host *rs.Host) bool { return m.m.MatchString(host.Hier.Dc.Name) }

func (m hostClusterMatcher) Match(host *rs.Host) bool { return m.m.MatchString(host.Hier.Cluster.Name) }

func (m hostHostMatcher) Match(host *rs.Host) bool { return m.m.MatchString(host.Name) }

func (m vmDcMatcher) Match(vm *rs.VM) bool { return m.m.MatchString(vm.Hier.Dc.Name) }

func (m vmClusterMatcher) Match(vm *rs.VM) bool { return m.m.MatchString(vm.Hier.Cluster.Name) }

func (m vmHostMatcher) Match(vm *rs.VM) bool { return m.m.MatchString(vm.Hier.Host.Name) }

func (m vmVMMatcher) Match(vm *rs.VM) bool { return m.m.MatchString(vm.Name) }

func (m orHostMatcher) Match(host *rs.Host) bool { return m.lhs.Match(host) || m.rhs.Match(host) }

func (m orVMMatcher) Match(vm *rs.VM) bool { return m.lhs.Match(vm) || m.rhs.Match(vm) }

func (m andHostMatcher) Match(host *rs.Host) bool { return m.lhs.Match(host) && m.rhs.Match(host) }

func (m andVMMatcher) Match(vm *rs.VM) bool { return m.lhs.Match(vm) && m.rhs.Match(vm) }

func newAndHostMatcher(lhs, rhs HostMatcher, others ...HostMatcher) andHostMatcher {
	m := andHostMatcher{lhs: lhs, rhs: rhs}
	if len(others) > 0 {
		return newAndHostMatcher(m, others[0], others[1:]...)
	}
	return m
}

func newAndVMMatcher(lhs, rhs VMMatcher, others ...VMMatcher) andVMMatcher {
	m := andVMMatcher{lhs: lhs, rhs: rhs}
	if len(others) > 0 {
		return newAndVMMatcher(m, others[0], others[1:]...)
	}
	return m
}

func newOrHostMatcher(lhs, rhs HostMatcher, others ...HostMatcher) orHostMatcher {
	m := orHostMatcher{lhs: lhs, rhs: rhs}
	if len(others) > 0 {
		return newOrHostMatcher(m, others[0], others[1:]...)
	}
	return m
}

func newOrVMMatcher(lhs, rhs VMMatcher, others ...VMMatcher) orVMMatcher {
	m := orVMMatcher{lhs: lhs, rhs: rhs}
	if len(others) > 0 {
		return newOrVMMatcher(m, others[0], others[1:]...)
	}
	return m
}

type (
	VMIncludes   []string
	HostIncludes []string
)

func (vi VMIncludes) Parse() (VMMatcher, error) {
	var ms []VMMatcher
	for _, v := range vi {
		m, err := parseVMIncludeString(v)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	switch len(ms) {
	case 0:
		return nil, nil
	case 1:
		return ms[0], nil
	default:
		return newOrVMMatcher(ms[0], ms[1], ms[2:]...), nil
	}
}

func (hi HostIncludes) Parse() (HostMatcher, error) {
	var ms []HostMatcher
	for _, v := range hi {
		m, err := parseHostIncludeString(v)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	switch len(ms) {
	case 0:
		return nil, nil
	case 1:
		return ms[0], nil
	default:
		return newOrHostMatcher(ms[0], ms[1], ms[2:]...), nil
	}
}

const (
	datacenter = iota
	cluster
	host
	vm
)

func parseHostIncludeString(s string) (HostMatcher, error) {
	s = strings.Trim(s, "/")
	// /dc/cluster/host
	parts := strings.Split(s, "/")
	var ms []HostMatcher
	for i, v := range parts {
		m, err := parseMatchSubString(v)
		if err != nil {
			return nil, err
		}
		switch i {
		case datacenter:
			ms = append(ms, hostDcMatcher{m})
		case cluster:
			ms = append(ms, hostClusterMatcher{m})
		case host:
			ms = append(ms, hostHostMatcher{m})
		}
	}
	switch len(ms) {
	case 0:
		return nil, nil
	case 1:
		return ms[0], nil
	default:
		return newAndHostMatcher(ms[0], ms[1], ms[2:]...), nil
	}
}

func parseVMIncludeString(s string) (VMMatcher, error) {
	s = strings.Trim(s, "/")
	// /dc/cluster/host/vm
	parts := strings.Split(s, "/")
	var ms []VMMatcher
	for i, v := range parts {
		m, err := parseMatchSubString(v)
		if err != nil {
			return nil, err
		}
		switch i {
		case datacenter:
			ms = append(ms, vmDcMatcher{m})
		case cluster:
			ms = append(ms, vmClusterMatcher{m})
		case host:
			ms = append(ms, vmHostMatcher{m})
		case vm:
			ms = append(ms, vmVMMatcher{m})
		}
	}
	switch len(ms) {
	case 0:
		return nil, nil
	case 1:
		return ms[0], nil
	default:
		return newAndVMMatcher(ms[0], ms[1], ms[2:]...), nil
	}
}

func parseMatchSubString(sub string) (matcher.Matcher, error) {
	sub = strings.TrimSpace(sub)
	if sub == "" || sub == "!*" {
		return matcher.FALSE(), nil
	}
	if sub == "*" {
		return matcher.TRUE(), nil
	}
	return matcher.NewSimplePatternsMatcher(sub)
}
