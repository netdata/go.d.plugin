package traefik

import (
	"errors"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricEntrypointRequestsTotal               = "traefik_entrypoint_requests_total"
	metricEntrypointRequestDurationSecondsSum   = "traefik_entrypoint_request_duration_seconds_sum"
	metricEntrypointRequestDurationSecondsCount = "traefik_entrypoint_request_duration_seconds_count"
	metricEntrypointOpenConnections             = "traefik_entrypoint_open_connections"
)

const (
	prefixEntrypointRequests  = "entrypoint_requests_"
	prefixEntrypointReqDurAvg = "entrypoint_request_duration_average_"
	prefixEntrypointOpenConn  = "entrypoint_open_connections_"
)

func isTraefikMetrics(pms prometheus.Metrics) bool {
	for _, pm := range pms {
		if strings.HasPrefix(pm.Name(), "traefik_") {
			return true
		}
	}
	return false
}

func (t *Traefik) collect() (map[string]int64, error) {
	pms, err := t.prom.Scrape()
	if err != nil {
		return nil, err
	}

	if t.checkMetrics && !isTraefikMetrics(pms) {
		return nil, errors.New("unexpected metrics (not Traefik)")
	}
	t.checkMetrics = false

	mx := make(map[string]int64)

	t.collectEntrypointRequestsTotal(mx, pms)
	t.collectEntrypointRequestDuration(mx, pms)
	t.collectEntrypointOpenConnections(mx, pms)
	t.updateCodeClassMetrics(mx)

	return mx, nil
}

func (t *Traefik) collectEntrypointRequestsTotal(mx map[string]int64, pms prometheus.Metrics) {
	if pms = pms.FindByName(metricEntrypointRequestsTotal); pms.Len() == 0 {
		return
	}

	for _, pm := range pms {
		code := pm.Labels.Get("code")
		ep := pm.Labels.Get("entrypoint")
		proto := pm.Labels.Get("protocol")
		codeClass := getCodeClass(code)
		if code == "" || ep == "" || proto == "" || codeClass == "" {
			continue
		}

		key := prefixEntrypointRequests + ep + "_" + proto + "_" + codeClass
		mx[key] += int64(pm.Value)

		id := ep + "_" + proto
		ce := t.cacheGetOrPutEntrypoint(id)
		if ce.requests == nil {
			chart := newChartEntrypointRequests(ep, proto)
			ce.requests = chart
			if err := t.Charts().Add(chart); err != nil {
				t.Warning(err)
			}
		}
	}
}

func (t *Traefik) collectEntrypointRequestDuration(mx map[string]int64, pms prometheus.Metrics) {
	if pms = pms.FindByNames(
		metricEntrypointRequestDurationSecondsCount,
		metricEntrypointRequestDurationSecondsSum,
	); pms.Len() == 0 {
		return
	}

	reqCount, durSum := make(map[string]float64), make(map[string]float64)
	for _, pm := range pms {
		code := pm.Labels.Get("code")
		ep := pm.Labels.Get("entrypoint")
		proto := pm.Labels.Get("protocol")
		codeClass := getCodeClass(code)
		if code == "" || ep == "" || proto == "" || codeClass == "" {
			continue
		}

		key := ep + "_" + proto + "_" + codeClass
		if pm.Name() == metricEntrypointRequestDurationSecondsCount {
			reqCount[key] += pm.Value
		} else {
			durSum[key] += pm.Value
		}

		id := ep + "_" + proto
		ce := t.cacheGetOrPutEntrypoint(id)
		if ce.reqDur == nil {
			chart := newChartEntrypointRequestDuration(ep, proto)
			ce.reqDur = chart
			if err := t.Charts().Add(chart); err != nil {
				t.Warning(err)
			}
		}
	}
	for k, count := range reqCount {
		if sum, ok := durSum[k]; ok && count > 0 {
			mx[prefixEntrypointReqDurAvg+k] = int64(sum * 1000 / count)
		} else {
			mx[prefixEntrypointReqDurAvg+k] = 0
		}
	}
}

func (t *Traefik) collectEntrypointOpenConnections(mx map[string]int64, pms prometheus.Metrics) {
	if pms = pms.FindByName(metricEntrypointOpenConnections); pms.Len() == 0 {
		return
	}

	for _, pm := range pms {
		method := pm.Labels.Get("method")
		ep := pm.Labels.Get("entrypoint")
		proto := pm.Labels.Get("protocol")
		if method == "" || ep == "" || proto == "" {
			continue
		}

		key := prefixEntrypointOpenConn + ep + "_" + proto + "_" + method
		mx[key] += int64(pm.Value)

		id := ep + "_" + proto
		ce := t.cacheGetOrPutEntrypoint(id)
		if ce.openConn == nil {
			chart := newChartEntrypointOpenConnections(ep, proto)
			ce.openConn = chart
			if err := t.Charts().Add(chart); err != nil {
				t.Warning(err)
			}
		}

		if !ce.openConnMethods[method] {
			ce.openConnMethods[method] = true
			dim := &module.Dim{ID: key, Name: method}
			if err := ce.openConn.AddDim(dim); err != nil {
				t.Warning(err)
			}
		}
	}
}

var httpRespCodeClasses = []string{"1xx", "2xx", "3xx", "4xx", "5xx"}

func (t Traefik) updateCodeClassMetrics(mx map[string]int64) {
	for k, v := range t.cache.entrypoints {
		if v.requests != nil {
			for _, c := range httpRespCodeClasses {
				key := prefixEntrypointRequests + k + "_" + c
				mx[key] += 0
			}
		}
		if v.reqDur != nil {
			for _, c := range httpRespCodeClasses {
				key := prefixEntrypointReqDurAvg + k + "_" + c
				mx[key] += 0
			}
		}
	}
}

func getCodeClass(code string) string {
	if len(code) != 3 {
		return ""
	}
	return string(code[0]) + "xx"
}

func (t *Traefik) cacheGetOrPutEntrypoint(id string) *cacheEntrypoint {
	if _, ok := t.cache.entrypoints[id]; !ok {
		t.cache.entrypoints[id] = &cacheEntrypoint{
			openConnMethods: make(map[string]bool),
		}
	}
	return t.cache.entrypoints[id]
}
