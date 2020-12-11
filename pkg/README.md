# Helper Packages

- if you need IP ranges consider to
  use [`iprange`](https://github.com/netdata/go.d.plugin/tree/master/pkg/iprange#iprange).
- if you parse an application log files, then [`log`](https://github.com/netdata/go.d.plugin/tree/master/pkg/logs) is
  handy.
- if you need filtering
  check [`matcher`](https://github.com/netdata/go.d.plugin/tree/master/pkg/matcher#supported-format).
- if you collect metrics from an HTTP endpoint use [`web`](https://github.com/netdata/go.d.plugin/tree/master/pkg/web).
- if you collect metrics from a prometheus endpoint,
  then [`prometheus`](https://github.com/netdata/go.d.plugin/tree/master/pkg/prometheus)
  and [`web`](https://github.com/netdata/go.d.plugin/tree/master/pkg/web) is what you need.
- [`tlscfg`](https://github.com/netdata/go.d.plugin/tree/master/pkg/tlscfg) provides TLS support.
- [`stm`](https://github.com/netdata/go.d.plugin/tree/master/pkg/stm) helps you to convert any struct to
  a `map[string]int64`.