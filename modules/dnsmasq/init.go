package dnsmasq

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (d Dnsmasq) validateConfig() error {
	if d.Address == "" {
		return errors.New("'address' parameter not set")
	}
	if !isNetworkValid(d.Network) {
		return fmt.Errorf("'network' (%s) is not valid, expected one of %v", d.Network, validNetworks)
	}
	return nil
}

func (d Dnsmasq) initDNSClient() (dnsClient, error) {
	return d.newDNSClient(d.Network, d.Timeout.Duration), nil
}

func (d Dnsmasq) initCharts() (*module.Charts, error) {
	return cacheCharts.Copy(), nil
}

func isNetworkValid(network string) bool {
	for _, v := range validNetworks {
		if network == v {
			return true
		}
	}
	return false
}

var validNetworks = []string{
	"udp",
	"tcp",
	"tcp-tls",
}
