package weblog

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/modules/weblog/parser"

	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go.d.plugin/pkg/simpletail"

	"github.com/netdata/go.d.plugin/pkg/logreader"

	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("web_log", creator)
}

func New() *WebLog {
	return &WebLog{
		Config: Config{
			Parser: parser.DefaultConfig,
		},
	}
}

type (
	Config struct {
		Parser                 parser.Config      `yaml:",inline"`
		Path                   string             `yaml:"path" validate:"required"`
		ExcludePath            string             `yaml:"exclude_path"`
		Filter                 matcher.SimpleExpr `yaml:"filter"`
		URLCategories          []RawCategory      `yaml:"categories"`
		UserCategories         []RawCategory      `yaml:"user_categories"`
		Histogram              []float64          `yaml:"histogram"`
		AggregateResponseCodes bool               `yaml:"aggregate_response_codes"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`
		charts *module.Charts

		file   *logreader.Reader
		parser parser.Parser

		metrics        *MetricsData
		filter         matcher.Matcher
		urlCategories  []*Category
		userCategories []*Category
	}
)

func (w *WebLog) Init() bool {
	if err := w.initFilter(); err != nil {
		w.Error(err)
		return false
	}

	if err := w.initCategories(); err != nil {
		w.Error(err)
		return false
	}

	w.metrics = NewMetricsData(w.Config)

	return true
}

func (w *WebLog) Check() bool {
	if err := w.initLogReader(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}
	lastLine, err := simpletail.ReadLastLine(w.file.CurrentFilename(), 0)
	if err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	w.parser, err = parser.NewParser(w.Config.Parser, w.file, lastLine)
	if err != nil {
		w.Warning("check failed: ", err)
		return false
	}
	log, err := w.parser.Parse(lastLine)
	if err != nil {
		w.Warning("check failed: ", err)
		return false
	}
	if err = log.Verify(); err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	return true
}

func (w *WebLog) Charts() *module.Charts {
	charts := make(module.Charts, 0, 10)
	_ = charts.Add(requests.Copy(), responseStatuses.Copy(), responseCodes.Copy())
	if w.AggregateResponseCodes {
		_ = charts.Add(responseCodesDetailedPerFamily()...)
	} else {
		_ = charts.Add(responseCodesDetailed.Copy())
	}
	_ = charts.Add(bandwidth.Copy())

	_ = charts.Add(requestsPerHTTPMethod.Copy())
	_ = charts.Add(requestsPerHTTPVersion.Copy())

	if len(w.urlCategories) > 0 {
		chart := requestsPerURL.Copy()
		for _, category := range w.urlCategories {
			_ = chart.AddDim(&Dim{
				ID:   category.name, //TODO: fix name
				Algo: module.Incremental,
			})
			for _, catChart := range perCategoryStats(category.name) {
				_ = charts.Add(catChart)
			}
		}
		_ = charts.Add(chart)
	}

	if len(w.userCategories) > 0 {
		chart := requestsPerUserDefined.Copy()
		for _, category := range w.userCategories {
			_ = chart.AddDim(&Dim{
				ID:   category.name,
				Algo: module.Incremental,
			})
		}
		_ = charts.Add(chart)
	}

	_ = charts.Add(responseTime.Copy())
	if len(w.Histogram) > 0 {
		chart := responseTimeHistogram.Copy()
		_ = charts.Add(chart)
		for _, v := range w.Histogram {
			name := fmt.Sprintf("%f", v)
			_ = chart.AddDim(&Dim{
				ID:   name, //FIXME
				Name: name,
				Algo: module.Incremental,
			})
		}
	}

	_ = charts.Add(responseTimeUpstream.Copy())
	if len(w.Histogram) > 0 {
		chart := responseTimeUpstreamHistogram.Copy()
		_ = charts.Add(chart)
		for _, v := range w.Histogram {
			name := fmt.Sprintf("%f", v)
			_ = chart.AddDim(&Dim{
				ID:   name, //FIXME
				Name: name,
				Algo: module.Incremental,
			})
		}
	}

	_ = charts.Add(requestsPerVhost.Copy())

	_ = charts.Add(requestsPerIPProto.Copy())
	_ = charts.Add(currentPollIPs.Copy())

	w.charts = &charts
	return w.charts
}

func (w *WebLog) Collect() map[string]int64 {
	defer func() {
		if err := recover(); err != nil {
			w.Errorf("[ERROR] %s\n", err)
			for depth := 0; ; depth++ {
				_, file, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				w.Errorf("======> %d: %v:%d", depth, file, line)
			}
			panic(err)
		}
	}()
	w.metrics.Reset()

	for {
		line, err := w.parser.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			w.Logger.Errorf("collect error: %v", err)
			return nil
		}

		if !w.filter.MatchString(line.URI) {
			continue
		}

		w.metrics.Requests.Inc()

		if line.Status != parser.EmptyNumber {
			status := line.Status
			switch {
			case status >= 100 && status < 300, status == 304:
				w.metrics.RespSuccessful.Inc()
			case status >= 300 && status < 400:
				w.metrics.RespRedirect.Inc()
			case status >= 400 && status < 500:
				w.metrics.RespClientError.Inc()
			case status >= 500 && status < 600:
				w.metrics.RespServerError.Inc()
			}

			switch status / 100 {
			case 1:
				w.metrics.Resp1xx.Inc()
			case 2:
				w.metrics.Resp2xx.Inc()
			case 3:
				w.metrics.Resp3xx.Inc()
			case 4:
				w.metrics.Resp4xx.Inc()
			case 5:
				w.metrics.Resp5xx.Inc()
			}

			statusStr := strconv.Itoa(status)
			counter, ok := w.metrics.RespCode.GetP(statusStr)
			counter.Inc()
			if !ok {
				if w.AggregateResponseCodes {
					chartName := fmt.Sprintf(`%s_%dxx`, responseCodesDetailed.ID, status/100)
					w.charts.Get(chartName).AddDim(&module.Dim{
						ID:   "req_code_" + statusStr,
						Name: statusStr,
						Algo: module.Incremental,
					})
				} else {
					w.charts.Get(responseCodesDetailed.ID).AddDim(&module.Dim{
						ID:   "req_code_" + statusStr,
						Name: statusStr,
						Algo: module.Incremental,
					})
				}
			}
		}
		if line.Method != parser.EmptyString {
			counter, ok := w.metrics.ReqMethod.GetP(line.Method)
			counter.Inc()
			if !ok && line.Method != "GET" {
				w.charts.Get(requestsPerHTTPMethod.ID).AddDim(&module.Dim{
					ID:   "req_method_" + line.Method,
					Name: line.Method,
					Algo: module.Incremental,
				})
			}
		}

		if line.Version != parser.EmptyString {
			deDotVersion := strings.Replace(line.Version, ".", "_", 1)
			c, ok := w.metrics.ReqVersion.GetP(deDotVersion)
			c.Inc()
			if !ok {
				w.charts.Get(requestsPerHTTPVersion.ID).AddDim(&module.Dim{
					ID:   "req_version_" + deDotVersion,
					Name: line.Version,
					Algo: module.Incremental,
				})
			}
		}

		if line.RespSize != parser.EmptyNumber {
			w.metrics.BytesSent.Add(float64(line.RespSize))
		}
		if line.ReqSize != parser.EmptyNumber {
			w.metrics.BytesReceived.Add(float64(line.ReqSize))
		}

		if line.RespTime != parser.EmptyNumber {
			w.metrics.RespTime.Observe(line.RespTime)
			if w.metrics.RespTimeHist != nil {
				w.metrics.RespTimeHist.Observe(line.RespTime)
			}
		}
		if line.UpstreamRespTime != parser.EmptyNumber {
			w.metrics.RespTimeUpstream.Observe(line.UpstreamRespTime)
			if w.metrics.RespTimeUpstreamHist != nil {
				w.metrics.RespTimeUpstreamHist.Observe(line.UpstreamRespTime)
			}
		}

		if line.Client != parser.EmptyString {
			if strings.ContainsRune(line.Client, ':') {
				w.metrics.ReqIpv6.Inc()
				w.metrics.UniqueIPv6.Insert(line.Client)
			} else {
				w.metrics.ReqIpv4.Inc()
				w.metrics.UniqueIPv4.Insert(line.Client)
			}
		}
		if line.URI != parser.EmptyString {
			for _, cat := range w.urlCategories {
				if cat.Matcher.MatchString(line.URI) {
					// TODO add metrics
					break
				}
			}
		}

		if line.Custom != parser.EmptyString {
			for _, cat := range w.userCategories {
				if cat.Matcher.MatchString(line.Custom) {
					// TODO add metrics
					break
				}
			}
		}
	}

	result := stm.ToMap(w.metrics)
	return result
}

func (w *WebLog) Cleanup() {
	w.file.Close()
}

func (w *WebLog) initLogReader() error {
	file, err := logreader.Open(w.Path, w.ExcludePath, w.Logger)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *WebLog) initFilter() (err error) {
	if w.Filter.Empty() {
		w.filter = matcher.TRUE()
		return
	}
	m, err := w.Filter.Parse()
	if err != nil {
		return fmt.Errorf("error on creating filter %s: %v", w.Filter, err)
	}
	w.filter = m
	return
}

func (w *WebLog) initCategories() error {
	for _, raw := range w.URLCategories {
		cat, err := NewCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating Category %s: %v", raw, err)
		}
		w.urlCategories = append(w.urlCategories, cat)
	}

	for _, raw := range w.UserCategories {
		cat, err := NewCategory(raw)
		if err != nil {
			return fmt.Errorf("error on creating Category %s: %v", raw, err)
		}
		w.userCategories = append(w.userCategories, cat)
	}

	return nil
}
