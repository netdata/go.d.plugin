package match

import (
	"fmt"
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
	hostMatcher struct {
		dc      matcher.Matcher
		cluster matcher.Matcher
		host    matcher.Matcher
	}
	vmMatcher struct {
		dc      matcher.Matcher
		cluster matcher.Matcher
		host    matcher.Matcher
		vm      matcher.Matcher
	}
	orHostMatcher struct{ lhs, rhs HostMatcher }
	orVMMatcher   struct{ lhs, rhs VMMatcher }
)

func (m hostMatcher) Match(host *rs.Host) bool {
	if !m.dc.MatchString(host.Hier.Dc.Name) {
		return false
	}
	if !m.cluster.MatchString(host.Hier.Cluster.Name) {
		return false
	}
	return m.host.MatchString(host.Name)
}

func (m vmMatcher) Match(vm *rs.VM) bool {
	if !m.dc.MatchString(vm.Hier.Dc.Name) {
		return false
	}
	if !m.cluster.MatchString(vm.Hier.Cluster.Name) {
		return false
	}
	if !m.host.MatchString(vm.Hier.Host.Name) {
		return false
	}
	return m.vm.MatchString(vm.Name)
}

func (m orHostMatcher) Match(host *rs.Host) bool {
	return m.lhs.Match(host) || m.rhs.Match(host)
}

func (m orVMMatcher) Match(vm *rs.VM) bool {
	return m.lhs.Match(vm) || m.rhs.Match(vm)
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

func parseHostIncludeString(s string) (HostMatcher, error) {
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}
	// /dc/cluster/host
	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("bad format : %s", s)
	}
	var hm hostMatcher
	for i, v := range parts {
		m, err := parseMatchSubString(v)
		if err != nil {
			return nil, err
		}
		switch i {
		case 0:
			hm.dc = m
		case 1:
			hm.cluster = m
		case 2:
			hm.host = m
		}
	}
	return &hm, nil
}

func parseVMIncludeString(s string) (VMMatcher, error) {
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}
	// /dc/cluster/host/vm
	parts := strings.Split(s, "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("bad format : %s", s)
	}
	var vmm vmMatcher
	for i, v := range parts {
		m, err := parseMatchSubString(v)
		if err != nil {
			return nil, err
		}
		switch i {
		case 0:
			vmm.dc = m
		case 1:
			vmm.cluster = m
		case 2:
			vmm.host = m
		case 3:
			vmm.vm = m
		}
	}
	return &vmm, nil
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
