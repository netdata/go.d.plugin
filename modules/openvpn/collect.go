package openvpn

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

var (
	reLoadStats = regexp.MustCompile(`SUCCESS: nclients=([0-9]+),bytesin=([0-9]+),bytesout=([0-9]+)`)
	reVersion   = regexp.MustCompile(`OpenVPN Version: OpenVPN ([0-9]+)\.([0-9]+)\.([0-9]+) .+?Management Version: ([0-9])`)
)

type loadStats struct {
	Clients  int `stm:"clients"`
	BytesIn  int `stm:"bytes_in"`
	BytesOut int `stm:"bytes_out"`
}

type version struct {
	major      int
	minor      int
	patch      int
	management int
}

func (o *OpenVPN) collect() (map[string]int64, error) {
	if !o.apiClient.isConnected() {
		err := o.apiClient.reconnect()
		if err != nil {
			return nil, err
		}
	}

	ls, err := o.collectLoadStats()
	if err != nil {
		_ = o.apiClient.disconnect()
		return nil, err
	}

	return stm.ToMap(ls), nil
}

func (o *OpenVPN) collectVersion() (*version, error) {
	if err := o.apiClient.send(commandVersion); err != nil {
		return nil, err
	}

	// one line is enough
	resp, err := o.apiClient.read(func(s string) bool { return true })
	if err != nil {
		return nil, err
	}
	return parseVersion(resp)
}

func (o *OpenVPN) collectLoadStats() (*loadStats, error) {
	err := o.apiClient.send(commandLoadStats)
	if err != nil {
		return nil, err
	}
	// one line is enough
	resp, err := o.apiClient.read(func(s string) bool { return true })
	if err != nil {
		return nil, err
	}
	return parseLoadStats(resp)
}

func parseVersion(raw []string) (*version, error) {
	m := reVersion.FindStringSubmatch(strings.Join(raw, " "))
	if len(m) == 0 {
		return nil, fmt.Errorf("regex parse filed, invalid format : %v", raw)
	}
	ver := version{
		major:      mustAtoi(m[1]),
		minor:      mustAtoi(m[2]),
		patch:      mustAtoi(m[3]),
		management: mustAtoi(m[4]),
	}
	return &ver, nil
}

func parseLoadStats(raw []string) (*loadStats, error) {
	m := reLoadStats.FindStringSubmatch(strings.Join(raw, " "))
	if len(m) == 0 {
		return nil, fmt.Errorf("regex parse filed, invalid format : %v", raw)
	}
	ls := loadStats{
		Clients:  mustAtoi(m[1]),
		BytesIn:  mustAtoi(m[2]),
		BytesOut: mustAtoi(m[3]),
	}
	return &ls, nil
}

func mustAtoi(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return v
}
