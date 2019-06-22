package weblog

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

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
		Config: Config{},
	}
}

type (
	Config struct {
		Path                   string             `yaml:"path" validate:"required"`
		ExcludePath            string             `yaml:"exclude_path"`
		Filter                 matcher.SimpleExpr `yaml:"filter"`
		URLCategories          []RawCategory      `yaml:"categories"`
		UserCategories         []RawCategory      `yaml:"user_categories"`
		LogFormat              string             `yaml:"log_format"`
		LogTimeScale           float64            `yaml:"log_time_scale"`
		Histogram              []float64          `yaml:"histogram"`
		AggregateResponseCodes bool               `yaml:"aggregate_response_codes"`
	}

	WebLog struct {
		module.Base
		Config `yaml:",inline"`
		charts *module.Charts

		file   *logreader.Reader
		parser *csv.Reader
		format *Format

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

	parser := NewCSVParser(ParserConfig{}, bytes.NewBuffer(lastLine))
	fields, err := parser.Read()
	if err != nil {
		w.Warning("check failed: ", err)
		return false
	}

	if w.LogFormat != "" {
		w.format = NewFormat(w.LogTimeScale, w.LogFormat)
		if w.format.Match(fields) != nil {
			w.Warning("check failed: ", err)
			return false
		}
	} else {
		w.format = GuessFormat(fields)
		if w.format == nil {
			w.Warning("check failed: cannot determine log format")
			return false
		}
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
	if w.format.BytesSent >= 0 || w.format.ReqLength >= 0 {
		_ = charts.Add(bandwidth.Copy())
	}

	if w.format.Request >= 0 {
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
	}

	if w.format.Custom >= 0 && len(w.userCategories) > 0 {
		chart := requestsPerUserDefined.Copy()
		for _, category := range w.userCategories {
			_ = chart.AddDim(&Dim{
				ID:   category.name,
				Algo: module.Incremental,
			})
		}
		_ = charts.Add(chart)
	}

	if w.format.ReqTime >= 0 {
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
	}

	if w.format.UpstreamRespTime >= 0 {
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
	}

	if w.format.Host >= 0 {
		_ = charts.Add(requestsPerVhost.Copy())
	}

	if w.format.RemoteAddr >= 0 {
		_ = charts.Add(requestsPerIPProto.Copy())
		_ = charts.Add(currentPollIPs.Copy())
	}

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
		fields, err := w.parser.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			w.Logger.Errorf("collect error: %v", err)
			return nil
		}
		line, err := w.format.Parse(fields)
		if err != nil {
			w.Logger.Errorf("parse error: %v", err)
			return nil
		}

		if !w.filter.MatchString(line.URI) {
			continue
		}

		w.metrics.Requests.Inc()

		if line.Status > 0 {
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
		if line.Method != "" {
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

		if line.Version != "" {
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

		if line.BytesSent > 0 {
			w.metrics.BytesSent.Add(float64(line.BytesSent))
		}
		if line.ReqLength > 0 {
			w.metrics.BytesReceived.Add(float64(line.ReqLength))
		}

		if line.ReqTime >= 0 {
			w.metrics.RespTime.Observe(line.ReqTime)
			if w.metrics.RespTimeHist != nil {
				w.metrics.RespTimeHist.Observe(line.ReqTime)
			}
		}
		if line.UpstreamRespTime != nil {
			for _, time := range line.UpstreamRespTime {
				w.metrics.RespTimeUpstream.Observe(time)
				if w.metrics.RespTimeUpstreamHist != nil {
					w.metrics.RespTimeUpstreamHist.Observe(line.ReqTime)
				}
			}
		}

		if line.RemoteAddr != "" {
			if strings.ContainsRune(line.RemoteAddr, ':') {
				w.metrics.ReqIpv6.Inc()
				w.metrics.UniqueIPv6.Insert(line.RemoteAddr)
			} else {
				w.metrics.ReqIpv4.Inc()
				w.metrics.UniqueIPv4.Insert(line.RemoteAddr)
			}
		}
		if line.URI != "" {
			for _, cat := range w.urlCategories {
				if cat.Matcher.MatchString(line.URI) {
					// TODO add metrics
					break
				}
			}
		}

		if w.format.Custom >= 0 {
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
	w.parser = NewCSVParser(ParserConfig{}, file)
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
