package unbound

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

type (
	stats []stat
	stat  struct {
		name  string
		value float64
	}
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

	for _, s := range stats {
		u.collectStat(s)
	}
	return nil
}

func (u *Unbound) collectStat(s stat) {
	switch {
	case strings.HasPrefix(s.name, "thread"):
		u.collectThread(s)
	case strings.HasPrefix(s.name, "total."):
		u.collectTotal(s)
	case strings.HasPrefix(s.name, "time."):
		u.collectTime(s)
	case strings.HasPrefix(s.name, "num.query.type."):
		u.collectQueryType(s)
	case strings.HasPrefix(s.name, "num.query.class."):
		u.collectQueryClass(s)
	case strings.HasPrefix(s.name, "num.query.opcode."):
		u.collectQueryOpCode(s)
	case strings.HasPrefix(s.name, "num.answer.rcode."):
		u.collectAnswerRCode(s)
	case strings.HasPrefix(s.name, "mem."):
		u.collectMem(s)
	}
}

func (u *Unbound) collectMem(s stat) {

}

func (u *Unbound) collectThread(s stat) {
	//thread0.*
	i := strings.IndexByte(s.name, '.')
	if i == -1 || i < 6 {
		return
	}
	threadID := s.name[6:i]
	c, ok := u.mx.Thread[threadID]
	if !ok {
		c = &common{}
		u.mx.Thread[threadID] = c
	}
	s.name = s.name[i+1:]
	collectCommon(c, s)
}

func (u *Unbound) collectTotal(s stat) {
	i := strings.IndexByte(s.name, '.')
	s.name = s.name[i+1:]
	collectCommon(&u.mx.common, s)
}

func (u *Unbound) collectTime(s stat) {
	switch s.name {
	case "time.up":
		u.mx.Uptime = s.value
	}
}

func (u *Unbound) collectQueryType(s stat) {
	i := len("num.query.type.")
	typ := s.name[i:]
	v, ok := u.mx.QueryType[typ]
	if !ok {
		//TODO:
	}
	u.mx.QueryType[typ] += v
}

func (u *Unbound) collectQueryClass(s stat) {
	i := len("num.query.class.")
	class := s.name[i:]
	v, ok := u.mx.QueryClass[class]
	if !ok {
		//TODO:
	}
	u.mx.QueryClass[class] += v
}

func (u *Unbound) collectQueryOpCode(s stat) {
	i := len("num.query.opcode.")
	opcode := s.name[i:]
	v, ok := u.mx.QueryOpCode[opcode]
	if !ok {
		//TODO:
	}
	u.mx.QueryOpCode[opcode] += v
}

func (u *Unbound) collectAnswerRCode(s stat) {
	i := len("num.answer.rcode.")
	rcode := s.name[i:]
	v, ok := u.mx.AnswerRCode[rcode]
	if !ok {
		//TODO:
	}
	u.mx.AnswerRCode[rcode] += v
}

func collectCommon(c *common, s stat) {
	switch s.name {
	case "num.queries":
		c.Queries = s.value
	case "num.queries_ip_ratelimited":
		c.QueriesIPRL = s.value
	case "num.cachehits":
		c.Cache.Hits = s.value
	case "num.cachemiss":
		c.Cache.Miss = s.value
	case "num.prefetch":
		c.Prefetch = s.value
	case "num.zero_ttl":
		c.ZeroTTL = s.value
	case "num.recursivereplies":
		c.RecursiveReplies = s.value
	case "num.dnscrypt.crypted":
		c.DNSCrypt.Crypted = s.value
	case "num.dnscrypt.cert":
		c.DNSCrypt.Cert = s.value
	case "num.dnscrypt.cleartext":
		c.DNSCrypt.ClearText = s.value
	case "num.dnscrypt.malformed":
		c.DNSCrypt.Malformed = s.value
	case "requestlist.avg":
		c.RequestList.Avg = s.value
	case "requestlist.max":
		c.RequestList.Max = s.value
	case "requestlist.overwritten":
		c.RequestList.Overwritten = s.value
	case "requestlist.exceeded":
		c.RequestList.Exceeded = s.value
	case "requestlist.current.all":
		c.RequestList.CurrentAll = s.value
	case "requestlist.current.user":
		c.RequestList.CurrentUser = s.value
	case "recursion.time.avg":
		c.RecursionTime.Avg = s.value
	case "recursion.time.median":
		c.RecursionTime.Median = s.value
	case "tcpusage":
		c.TCPUsage = s.value
	}
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
