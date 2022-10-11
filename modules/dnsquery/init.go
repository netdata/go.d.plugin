// SPDX-License-Identifier: GPL-3.0-or-later

package dnsquery

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/miekg/dns"
)

func (d *DNSQuery) verifyConfig() error {
	if len(d.Domains) == 0 {
		return errors.New("no domains specified")
	}

	if len(d.Servers) == 0 {
		return errors.New("no servers specified")
	}

	if !(d.Network == "" || d.Network == "udp" || d.Network == "tcp" || d.Network == "tcp-tls") {
		return fmt.Errorf("wrong network transport : %s", d.Network)
	}

	return nil
}

func (d *DNSQuery) initCharts() (*module.Charts, error) {
	var charts module.Charts

	for _, srv := range d.Servers {
		cs := newDNSServerCharts(srv, d.Network, d.RecordType)
		if err := charts.Add(*cs...); err != nil {
			return nil, err
		}
	}

	return &charts, nil
}

func (d *DNSQuery) initRecordType() (uint16, error) {
	return parseRecordType(d.RecordType)
}

func parseRecordType(recordType string) (uint16, error) {
	var rtype uint16

	switch recordType {
	case "A":
		rtype = dns.TypeA
	case "AAAA":
		rtype = dns.TypeAAAA
	case "ANY":
		rtype = dns.TypeANY
	case "CNAME":
		rtype = dns.TypeCNAME
	case "MX":
		rtype = dns.TypeMX
	case "NS":
		rtype = dns.TypeNS
	case "PTR":
		rtype = dns.TypePTR
	case "SOA":
		rtype = dns.TypeSOA
	case "SPF":
		rtype = dns.TypeSPF
	case "SRV":
		rtype = dns.TypeSRV
	case "TXT":
		rtype = dns.TypeTXT
	default:
		return 0, fmt.Errorf("unknown record type : %s", recordType)
	}

	return rtype, nil
}
