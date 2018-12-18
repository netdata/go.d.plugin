package weblog

import (
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/modules"
)

func (w *WebLog) cleanup() {
	w.tail.stop()
}

func (w *WebLog) parseLoop() {
	lines := w.tail.lines()
LOOP:
	for {
		select {
		case <-w.stop:
			w.cleanup()
			break LOOP
		case <-w.pause:
			w.pause <- struct{}{}
		case line := <-lines:
			if w.filter.match(line.Text) {
				w.parseLine(line.Text)
			}
		}
	}
}

func (w *WebLog) parseLine(line string) {
	gm, ok := w.parser.parse(line)

	if !ok {
		w.metrics["unmatched"]++
		return
	}

	w.codeFam(gm)

	w.codeStatus(gm)

	if w.DoCodesDetailed {
		w.codeDetailed(gm)
	}

	if _, ok := gm.lookup(keyRequest); ok {
		w.request(gm)
	}

	if _, ok := gm.lookup(keyUserDefined); ok && len(w.userCats) > 0 {
		w.userCategory(gm)
	}

	if _, ok := gm.lookup(keyBytesSent); ok {
		w.bytesSent(gm)
	}

	if _, ok := gm.lookup(keyRespLength); ok {
		w.respLength(gm)
	}

	if v, ok := gm.lookup(keyRespTime); ok {
		i := w.timings.get(keyRespTime).set(v)
		if h, ok := w.histograms[keyRespTimeHistogram]; ok {
			h.set(i)
		}
	}

	if v, ok := gm.lookup(keyRespTimeUpstream); ok {
		i := w.timings.get(keyRespTimeUpstream).set(v)
		if h, ok := w.histograms[keyRespTimeUpstreamHistogram]; ok {
			h.set(i)
		}
	}

	if _, ok := gm.lookup(keyAddress); ok {
		w.ipProto(gm)
	}

	if w.DoPerURLCharts && w.matchedURL != "" {
		w.urlCategoryStats(gm)
	}

}

func (w *WebLog) codeFam(gm groupMap) {
	fam := gm.get(keyCode)[:1] + "xx"

	if _, ok := w.metrics[fam]; ok {
		w.metrics[fam]++
	} else {
		w.metrics["0xx"]++
	}
}

func (w *WebLog) codeDetailed(gm groupMap) {
	code := gm.get(keyCode)

	if _, ok := w.metrics[code]; ok {
		w.metrics[code]++
		return
	}

	var chart *Chart

	if w.DoCodesAggregate {
		chart = w.charts.Get(responseCodesDetailed.ID)
	} else {
		v := "other"
		if code[0] <= 53 {
			v = code[:1] + "xx"
		}
		chart = w.charts.Get(responseCodesDetailed.ID + "_" + v)
	}

	_ = chart.AddDim(&Dim{
		ID:   code,
		Algo: modules.Incremental,
	})
	chart.MarkNotCreated()

	w.metrics[code]++
}

func (w *WebLog) codeStatus(gm groupMap) {
	code, fam := gm.get(keyCode), gm.get(keyCode)[:1]

	switch {
	case fam == "2", code == "304", fam == "1":
		w.metrics["successful_requests"]++
	case fam == "3":
		w.metrics["redirects"]++
	case fam == "4":
		w.metrics["bad_requests"]++
	case fam == "5":
		w.metrics["server_errors"]++
	default:
		w.metrics["other_requests"]++
	}
}

func (w *WebLog) request(gm groupMap) {
	var ok bool

	if gm, ok = w.reqParser.parse(gm.get(keyRequest)); !ok {
		return
	}

	w.httpMethod(gm)
	w.urlCategory(gm)
	w.httpVersion(gm)
}

func (w *WebLog) httpMethod(gm groupMap) {
	method := gm.get(keyMethod)

	if _, ok := w.metrics[method]; !ok {
		chart := w.charts.Get(requestsPerHTTPMethod.ID)
		_ = chart.AddDim(&Dim{
			ID:   method,
			Algo: modules.Incremental,
		})
		chart.MarkNotCreated()
	}

	w.metrics[method]++
}

func (w *WebLog) urlCategory(gm groupMap) {
	url := gm.get(keyURL)

	for _, v := range w.urlCats {
		if v.match(url) {
			w.metrics[v.name]++
			w.matchedURL = v.name
			return
		}
	}

	w.matchedURL = ""
	w.metrics[keyURL+"_category_other"]++
}

func (w *WebLog) userCategory(gm groupMap) {
	userDefined := gm.get(keyUserDefined)

	for _, cat := range w.userCats {
		if cat.match(userDefined) {
			w.metrics[cat.name]++
			return
		}
	}

	w.metrics[keyUserDefined+"_category_other"]++
}

func (w *WebLog) httpVersion(gm groupMap) {
	version := gm.get(keyVersion)
	dimID := strings.Replace(version, ".", "_", 1)

	if _, ok := w.metrics[dimID]; !ok {
		chart := w.charts.Get(requestsPerHTTPVersion.ID)
		_ = chart.AddDim(&Dim{
			ID:   dimID,
			Name: version,
			Algo: modules.Incremental,
		})
		chart.MarkNotCreated()
	}

	w.metrics[dimID]++
}

func (w *WebLog) bytesSent(gm groupMap) {
	w.metrics["bytes_sent"] += toInt(gm.get(keyBytesSent))
}

func (w *WebLog) respLength(gm groupMap) {
	w.metrics["resp_length"] += toInt(gm.get(keyRespLength))
}

func (w *WebLog) ipProto(gm groupMap) {
	var (
		address = gm.get(keyAddress)
		proto   = "ipv4"
	)

	if strings.Contains(address, ":") {
		proto = "ipv6"
	}

	w.metrics["req_"+proto]++

	if _, ok := w.uniqIPs[address]; !ok {
		w.uniqIPs[address] = true
		w.metrics["unique_current_poll_"+proto]++
	}

	if !w.DoAllTimeIPs {
		return
	}

	if _, ok := w.uniqIPsAllTime[address]; !ok {
		w.uniqIPsAllTime[address] = true
		w.metrics["unique_all_time_"+proto]++
	}

}

func (w *WebLog) urlCategoryStats(gm groupMap) {
	code := gm.get(keyCode)
	id := w.matchedURL + "_" + code

	if _, ok := w.metrics[id]; !ok {
		chart := w.charts.Get(responseCodesDetailed.ID + "_" + w.matchedURL)
		_ = chart.AddDim(&Dim{
			ID:   id,
			Name: code,
			Algo: modules.Incremental,
		})
		chart.MarkNotCreated()
	}

	w.metrics[id]++

	if v, ok := gm.lookup(keyBytesSent); ok {
		w.metrics[w.matchedURL+"_bytes_sent"] += toInt(v)
	}

	if v, ok := gm.lookup(keyRespLength); ok {
		w.metrics[w.matchedURL+"_resp_length"] += toInt(v)
	}

	if v, ok := gm.lookup(keyRespTime); ok {
		w.timings.get(w.matchedURL).set(v)
	}
}

func toInt(s string) int64 {
	// TODO: 0.000
	if s == "-" {
		return 0
	}
	v, _ := strconv.Atoi(s)

	return int64(v)
}
