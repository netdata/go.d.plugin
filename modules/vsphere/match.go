package vsphere

import (
	"fmt"
	"strings"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type hostMatcher struct {
	dc      matcher.Matcher
	cluster matcher.Matcher
	host    matcher.Matcher
}

func (m hostMatcher) Match(host *rs.Host) bool {
	if !m.dc.MatchString(host.Hier.Dc.Name) {
		return false
	}
	if !m.cluster.MatchString(host.Hier.Cluster.Name) {
		return false
	}
	return m.host.MatchString(host.Name)
}

type vmMatcher struct {
	dc      matcher.Matcher
	cluster matcher.Matcher
	host    matcher.Matcher
	vm      matcher.Matcher
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

type orHostMatcher []*hostMatcher

func (ms orHostMatcher) Match(host *rs.Host) bool {
	for _, m := range ms {
		if m.Match(host) {
			return true
		}
	}
	return false
}

type orVMMatcher []*vmMatcher

func (ms orVMMatcher) Match(vm *rs.VM) bool {
	for _, m := range ms {
		if m.Match(vm) {
			return true
		}
	}
	return false
}

type (
	hostInclude string
	vmInclude   string
)

func (hi hostInclude) parse() (*hostMatcher, error) {
	s := string(hi)
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}

	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("bad format : %s", string(hi))
	}

	dcm, err := parseMatchSubString(parts[0])
	if err != nil {
		return nil, err
	}
	cm, err := parseMatchSubString(parts[1])
	if err != nil {
		return nil, err
	}
	hm, err := parseMatchSubString(parts[2])
	if err != nil {
		return nil, err
	}
	return &hostMatcher{dc: dcm, cluster: cm, host: hm}, nil
}

func (vi vmInclude) parse() (*vmMatcher, error) {
	s := string(vi)
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}

	parts := strings.Split(s, "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("bad format : %s", string(vi))
	}

	dcm, err := parseMatchSubString(parts[0])
	if err != nil {
		return nil, err
	}
	cm, err := parseMatchSubString(parts[1])
	if err != nil {
		return nil, err
	}
	hm, err := parseMatchSubString(parts[2])
	if err != nil {
		return nil, err
	}
	vm, err := parseMatchSubString(parts[3])
	if err != nil {
		return nil, err
	}
	return &vmMatcher{dc: dcm, cluster: cm, host: hm, vm: vm}, nil
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

func parseHostIncludes(includes []hostInclude) (*orHostMatcher, error) {
	if len(includes) == 0 {
		return nil, nil
	}
	var ms orHostMatcher
	for _, v := range includes {
		m, err := v.parse()
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return &ms, nil
}

func parseVMIncludes(includes []vmInclude) (*orVMMatcher, error) {
	if len(includes) == 0 {
		return nil, nil
	}
	var ms orVMMatcher
	for _, v := range includes {
		m, err := v.parse()
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return &ms, nil
}
