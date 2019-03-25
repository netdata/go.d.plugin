package k8s_kubeproxy

import (
	"reflect"

	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

func newMetrics() *metrics {
	var mx metrics
	value := reflect.Indirect(reflect.ValueOf(&mx))
	setMetrics(value, value.Type())
	return &mx
}

type metrics struct {
	SyncProxyRules struct {
		Count   mtx.Gauge `stm:"count"`
		Latency struct {
			LE1000     mtx.Gauge `stm:"1000"`
			LE2000     mtx.Gauge `stm:"2000"`
			LE4000     mtx.Gauge `stm:"4000"`
			LE8000     mtx.Gauge `stm:"8000"`
			LE16000    mtx.Gauge `stm:"16000"`
			LE32000    mtx.Gauge `stm:"32000"`
			LE64000    mtx.Gauge `stm:"64000"`
			LE128000   mtx.Gauge `stm:"128000"`
			LE256000   mtx.Gauge `stm:"256000"`
			LE512000   mtx.Gauge `stm:"512000"`
			LE1024000  mtx.Gauge `stm:"1024000"`
			LE2048000  mtx.Gauge `stm:"2048000"`
			LE4096000  mtx.Gauge `stm:"4096000"`
			LE8192000  mtx.Gauge `stm:"8192000"`
			LE16384000 mtx.Gauge `stm:"16384000"`
			Inf        mtx.Gauge `stm:"+Inf"`
		} `stm:"bucket"`
	} `stm:"sync_proxy_rules"`
	RESTClient struct {
		HTTPRequests struct {
			ByStatusCode map[string]mtx.Gauge `stm:""`
			ByMethod     map[string]mtx.Gauge `stm:""`
		} `stm:"http_requests"`
	} `stm:"rest_client"`
}

func setMetrics(v reflect.Value, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if ft.Type.Kind() == reflect.Struct {
			setMetrics(v.Field(i), ft.Type)
		}
		if ft.Type.Kind() != reflect.Map {
			continue
		}
		value := v.Field(i)
		if !value.IsNil() {
			continue
		}
		switch value.Interface().(type) {
		case map[string]mtx.Gauge:
			value.Set(reflect.ValueOf(map[string]mtx.Gauge{}))
		}
	}
}
