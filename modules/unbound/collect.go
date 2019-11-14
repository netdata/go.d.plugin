package unbound

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (u *Unbound) collect() (map[string]int64, error) {
	if err := u.collectStats(); err != nil {
		return nil, err
	}
	return stm.ToMap(u.mx), nil
}

func (u *Unbound) collectStats() error {
	resp, err := u.client.send("UBCT1 stats_noreset\n")
	if err != nil {
		return err
	}
	switch len(resp) {
	case 0:
		return errors.New("empty response")
	case 1:
		// 	In case of error the first line of the response is: error <descriptive text possible> \n
		//	For many commands the  response is 'ok\n', but it is not the case for 'stats'.
		return errors.New(resp[0])
	}

	stats, err := convertToUnboundStats(resp)
	if err != nil {
		return err
	}

	u.collectTotal(stats)
	return nil
}

func (u *Unbound) collectTotal(ss stats) {
	for _, s := range ss.find("total.") {
		v := s.value
		switch s.name {
		case "total.num.queries":
			u.mx.Total.Queries.Set(v)
		case "total.num.queries_ip_ratelimited":
			u.mx.Total.QueriesIPRL.Set(v)
		case "total.num.cachehits":
			u.mx.Total.Cache.Hits.Set(v)
		case "total.num.cachemiss":
			u.mx.Total.Cache.Miss.Set(v)
		case "total.num.prefetch":
			u.mx.Total.Prefetch.Set(v)
		case "total.num.zero_ttl":
			u.mx.Total.ZeroTTL.Set(v)
		case "total.num.recursivereplies":
			u.mx.Total.RecursiveReplies.Set(v)
		case "total.num.dnscrypt.crypted":
			u.mx.Total.DNSCrypt.Crypted.Set(v)
		case "total.num.dnscrypt.cert":
			u.mx.Total.DNSCrypt.Cert.Set(v)
		case "total.num.dnscrypt.cleartext":
			u.mx.Total.DNSCrypt.ClearText.Set(v)
		case "total.num.dnscrypt.malformed":
			u.mx.Total.DNSCrypt.Malformed.Set(v)
		case "total.requestlist.avg":
			u.mx.Total.RequestList.Avg.Set(v)
		case "total.requestlist.max":
			u.mx.Total.RequestList.Max.Set(v)
		case "total.requestlist.overwritten":
			u.mx.Total.RequestList.Overwritten.Set(v)
		case "total.requestlist.exceeded":
			u.mx.Total.RequestList.Exceeded.Set(v)
		case "total.requestlist.current.all":
			u.mx.Total.RequestList.CurrentAll.Set(v)
		case "total.requestlist.current.user":
			u.mx.Total.RequestList.CurrentUser.Set(v)
		case "total.recursion.time.avg":
			u.mx.Total.RecursionTime.Avg.Set(v)
		case "total.recursion.time.median":
			u.mx.Total.RecursionTime.Median.Set(v)
		case "total.tcpusage":
			u.mx.Total.TCPUsage.Set(v)
		}
	}
}

func (u *Unbound) collectTime(ss stats) {
	for _, s := range ss.find("time.") {
		switch s.name {
		case "time.up":
			u.mx.Uptime.Set(s.value)
		}
	}
}

type (
	stats []stat
	stat  struct {
		name  string
		value float64
	}
)

func (ss stats) empty() bool {
	return len(ss) == 0
}

func (ss stats) find(prefix string) stats {
	from := sort.Search(len(ss), func(i int) bool { return ss[i].name >= prefix })
	if from == len(ss) || !strings.HasPrefix(ss[from].name, prefix) {
		return nil
	}
	until := from + 1
	for until < len(ss) && strings.HasPrefix(ss[until].name, prefix) {
		until++
	}
	return ss[from:until]
}

func convertToUnboundStats(lines []string) (stats, error) {
	sort.Strings(lines)
	var ubs stats
	for _, line := range lines {
		// 'stats' output is a list of [name]=[value] lines.
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad stats syntax: %s", line)
		}

		name, value := parts[0], parts[1]
		v, err := strconv.ParseFloat(value, 10)
		if err != nil {
			return nil, err
		}

		ubs = append(ubs, stat{name, v})
	}
	return ubs, nil
}
