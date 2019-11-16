package unbound

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// https://github.com/NLnetLabs/unbound/blob/master/smallapp/unbound-control.c
// https://github.com/NLnetLabs/unbound/blob/master/libunbound/unbound.h (ub_server_stats, ub_shm_stat_info)
// https://docs.menandmice.com/display/MM/Unbound+request-list+demystified
// https://docs.datadoghq.com/integrations/unbound/#metrics

func (u *Unbound) collect() (map[string]int64, error) {
	stats, err := u.scrapeUnboundStats()
	if err != nil {
		return nil, err
	}

	u.curCache.clear()
	mx := u.collectStats(stats)
	u.updateCharts()
	return mx, nil
}

func (u *Unbound) scrapeUnboundStats() ([]entry, error) {
	output, err := u.client.send("UBCT1 stats\n")
	if err != nil {
		return nil, err
	}
	switch len(output) {
	case 0:
		return nil, errors.New("empty response")
	case 1:
		// 	in case of error the first line of the response is: error <descriptive text possible> \n
		return nil, errors.New(output[0])
	}
	return parseStatsOutput(output)
}

func (u *Unbound) collectStats(stats []entry) map[string]int64 {
	mx := make(map[string]int64, len(stats))
	for _, e := range stats {
		if e.hasPrefix("histogram") {
			continue
		}
		// *.requestlist.avg, *.recursion.time.avg, *recursion.time.median
		if e.hasSuffix("avg") || e.hasSuffix("median") {
			e.value *= 1000
		}
		switch {
		case e.hasPrefix("thread"):
			v := extractThreadID(e.key)
			u.curCache.threads[v] = true
		case e.hasPrefix("num.query.type"):
			v := extractQueryType(e.key)
			u.curCache.queryType[v] = true
		case e.hasPrefix("num.query.class"):
			v := extractQueryClass(e.key)
			u.curCache.queryClass[v] = true
		case e.hasPrefix("num.query.opcode"):
			v := extractQueryOpCode(e.key)
			u.curCache.queryOpCode[v] = true
		case e.hasPrefix("num.query.flags"):
			v := extractQueryFlag(e.key)
			u.curCache.queryFlags[v] = true
		case e.hasPrefix("num.answer.rcode"):
			v := extractAnswerRCode(e.key)
			u.curCache.answerRCode[v] = true
		}
		mx[e.key] = int64(e.value)
	}
	return mx
}

func extractThreadID(key string) string    { idx := strings.IndexByte(key, '.'); return key[6:idx] }
func extractQueryType(key string) string   { i := len("num.query.type."); return key[i:] }
func extractQueryClass(key string) string  { i := len("num.query.class."); return key[i:] }
func extractQueryOpCode(key string) string { i := len("num.query.opcode."); return key[i:] }
func extractQueryFlag(key string) string   { i := len("num.query.flags."); return key[i:] }
func extractAnswerRCode(key string) string { i := len("num.answer.rcode."); return key[i:] }

type entry struct {
	key   string
	value float64
}

func (e entry) hasPrefix(prefix string) bool { return strings.HasPrefix(e.key, prefix) }
func (e entry) hasSuffix(suffix string) bool { return strings.HasSuffix(e.key, suffix) }

func parseStatsOutput(output []string) ([]entry, error) {
	var es []entry
	for _, v := range output {
		e, err := parseStatsLine(v)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, nil
}

func parseStatsLine(line string) (entry, error) {
	// 'stats' output is a list of [key]=[value] lines.
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		return entry{}, fmt.Errorf("bad line syntax: %s", line)
	}
	f, err := strconv.ParseFloat(parts[1], 10)
	return entry{key: parts[0], value: f}, err
}

func newCollectCache() collectCache {
	return collectCache{
		threads:     make(map[string]bool),
		queryType:   make(map[string]bool),
		queryClass:  make(map[string]bool),
		queryOpCode: make(map[string]bool),
		queryFlags:  make(map[string]bool),
		answerRCode: make(map[string]bool),
	}
}

type collectCache struct {
	threads     map[string]bool
	queryType   map[string]bool
	queryClass  map[string]bool
	queryOpCode map[string]bool
	queryFlags  map[string]bool
	answerRCode map[string]bool
}

func (c *collectCache) clear() {
	*c = newCollectCache()
}
