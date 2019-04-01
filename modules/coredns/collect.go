package coredns

import "github.com/netdata/go.d.plugin/pkg/stm"

func (cd *CoreDNS) collect() (map[string]int64, error) {
	raw, err := cd.prom.Scrape()

	if err != nil {
		return nil, err
	}

	mx := metrics{}
	mx.PanicCount.Set(raw.FindByName("coredns_panic_count_total").Max())

	return stm.ToMap(mx), nil
}
