package dnsquery

import (
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func (d *DNSQuery) collect() (map[string]int64, error) {
	if d.dnsClient == nil {
		d.dnsClient = d.newDNSClient(d.Network, d.Timeout.Duration)
	}

	mx := make(map[string]int64)
	domain := randomDomain(d.Domains)
	d.Debugf("current domain : %s", domain)

	var wg sync.WaitGroup
	var mux sync.RWMutex
	for _, srv := range d.Servers {
		wg.Add(1)
		go func(srv string, wg *sync.WaitGroup) {
			defer wg.Done()

			msg := new(dns.Msg)
			msg.SetQuestion(dns.Fqdn(domain), d.rtype)
			address := net.JoinHostPort(srv, strconv.Itoa(d.Port))

			resp, rtt, err := d.dnsClient.Exchange(msg, address)
			if err != nil {
				d.Debugf("error on querying %s after %s query for %s : %s", srv, d.RecordType, domain, err)
				return
			}
			if resp != nil && resp.Rcode != dns.RcodeSuccess {
				d.Errorf("invalid answer from %s after %s query for %s", srv, d.RecordType, domain)
				return
			}

			mux.Lock()
			mx["server_"+srv+"_query_time"] = rtt.Nanoseconds()
			mux.Unlock()
		}(srv, &wg)
	}
	wg.Wait()

	return mx, nil
}

func randomDomain(domains []string) string {
	rand.Seed(time.Now().UnixNano())
	return domains[rand.Intn(len(domains))]
}
