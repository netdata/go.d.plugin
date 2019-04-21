package tengine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

const (
	bytesIn               = "bytes_in"
	bytesOut              = "bytes_out"
	connTotal             = "conn_total"
	reqTotal              = "req_total"
	http2xx               = "http_2xx"
	http3xx               = "http_3xx"
	http4xx               = "http_4xx"
	http5xx               = "http_5xx"
	httpOtherStatus       = "http_other_status"
	rt                    = "rt"
	upsReq                = "ups_req"
	upsRT                 = "ups_rt"
	upsTries              = "ups_tries"
	http200               = "http_200"
	http206               = "http_206"
	http302               = "http_302"
	http304               = "http_304"
	http403               = "http_403"
	http404               = "http_404"
	http416               = "http_416"
	http499               = "http_499"
	http500               = "http_500"
	http502               = "http_502"
	http503               = "http_503"
	http504               = "http_504"
	http508               = "http_508"
	httpOtherDetailStatus = "http_other_detail_status"
	httpUps4xx            = "http_ups_4xx"
	httpUps5xx            = "http_ups_5xx"
)

var defaultLineFormat = []string{
	bytesIn,
	bytesOut,
	connTotal,
	reqTotal,
	http2xx,
	http3xx,
	http4xx,
	http5xx,
	httpOtherStatus,
	rt,
	upsReq,
	upsRT,
	upsTries,
	http200,
	http206,
	http302,
	http304,
	http403,
	http404,
	http416,
	http499,
	http500,
	http502,
	http503,
	http504,
	http508,
	httpOtherDetailStatus,
	httpUps4xx,
	httpUps5xx,
}

func (t *Tengine) collect() (map[string]int64, error) {
	raw, err := t.apiClient.getStatus()

	if err != nil {
		return nil, err
	}

	var ms metrics

	for _, line := range raw {
		var m metric

		err := parseLine(&m, line, defaultLineFormat)
		if err != nil {
			return nil, fmt.Errorf("error on parsing '%s' : %v", line, err)
		}

		ms = append(ms, m)
	}

	mx := make(map[string]int64)

	for _, m := range ms {
		for k, v := range stm.ToMap(m) {
			mx[k] += v
		}
	}

	return mx, nil
}

func parseLine(m *metric, line string, lineFormat []string) error {
	parts := strings.Split(line, ",")

	// NOTE: only default line format is supported
	// TODO: custom line format?
	// www.example.com,127.0.0.1:80,162,6242,1,1,1,0,0,0,0,10,1,10,1....
	if len(parts)-2 != len(lineFormat) {
		return fmt.Errorf("invalid response length, got %d, expected %d", len(parts), len(lineFormat))
	}

	m.Host = parts[0]
	m.ServerAddress = parts[1]

	for i, f := range lineFormat {
		// 1, 2: "$host,$server_addr:$server_port"
		value := parts[i+2]
		switch f {
		default:
			return fmt.Errorf("unknown line format value: %s", f)
		case bytesIn:
			m.BytesIn = mustParseInt(value)
		case bytesOut:
			m.BytesOut = mustParseInt(value)
		case connTotal:
			m.ConnTotal = mustParseInt(value)
		case reqTotal:
			m.ReqTotal = mustParseInt(value)
		case http2xx:
			m.HTTP2xx = mustParseInt(value)
		case http3xx:
			m.HTTP3xx = mustParseInt(value)
		case http4xx:
			m.HTTP4xx = mustParseInt(value)
		case http5xx:
			m.HTTP5xx = mustParseInt(value)
		case httpOtherStatus:
			m.HTTPOtherStatus = mustParseInt(value)
		case rt:
			m.RT = mustParseInt(value)
		case upsReq:
			m.UpsReq = mustParseInt(value)
		case upsRT:
			m.UpsRT = mustParseInt(value)
		case upsTries:
			m.UpsTries = mustParseInt(value)
		case http200:
			m.HTTP200 = mustParseInt(value)
		case http206:
			m.HTTP206 = mustParseInt(value)
		case http302:
			m.HTTP302 = mustParseInt(value)
		case http304:
			m.HTTP304 = mustParseInt(value)
		case http403:
			m.HTTP403 = mustParseInt(value)
		case http404:
			m.HTTP404 = mustParseInt(value)
		case http416:
			m.HTTP416 = mustParseInt(value)
		case http499:
			m.HTTP499 = mustParseInt(value)
		case http500:
			m.HTTP500 = mustParseInt(value)
		case http502:
			m.HTTP502 = mustParseInt(value)
		case http503:
			m.HTTP503 = mustParseInt(value)
		case http504:
			m.HTTP504 = mustParseInt(value)
		case http508:
			m.HTTP508 = mustParseInt(value)
		case httpOtherDetailStatus:
			m.HTTPOtherDetailStatus = mustParseInt(value)
		case httpUps4xx:
			m.HTTPUps4xx = mustParseInt(value)
		case httpUps5xx:
			m.HTTPUps5xx = mustParseInt(value)
		}
	}
	return nil
}

func mustParseInt(value string) *int64 {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}

	return &v
}
