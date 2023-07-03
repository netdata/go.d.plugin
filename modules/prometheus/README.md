<!--
title: "Prometheus endpoint monitoring with Netdata"
description: "Monitor the health and performance of 600+ services that support the Prometheus metrics with Netdata's per-second frequency and zero configuration."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/prometheus/README.md"
sidebar_label: "Prometheus endpoints"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Anything"
-->

# Prometheus endpoint collector

The generic Prometheus endpoint collector gathers metrics from [`Prometheus`](https://prometheus.io/) endpoints that use
the [OpenMetrics exposition format](https://prometheus.io/docs/instrumenting/exposition_formats/).

- As of v1.24, Netdata can autodetect more than 600 Prometheus endpoints, including support for Windows 10 via
  `windows_exporter`, and instantly generate new charts with the same high-granularity, per-second frequency as you
  expect from other collectors.

- The full list of endpoints is available in the
  collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/prometheus.conf).

- Collecting metrics
  from [Prometheus endpoints in Kubernetes](https://github.com/netdata/helmchart#prometheus-endpoints).

## Charts

Netdata will produce one or more charts for every metric collected via a Prometheus endpoint. The number of charts
depends entirely on the number of exposed metrics.

For example, scraping [`node_exporter`](https://github.com/prometheus/node_exporter) produces 3000+ metrics.

## Configuration

Edit the `go.d/prometheus.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/prometheus.conf
```

Endpoints can be either local or remote as long as they provide their metrics at the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics

  - name: win10
    url: http://203.0.113.0:9182/metrics
```

Here is an example pulling from the [Prometheus demo site](https://demo.do.prometheus.io/) node
exporter [endpoint](https://node.demo.do.prometheus.io/metrics):

```yaml
jobs:
  - name: node_exporter_demo
    url: https://node.demo.do.prometheus.io/metrics
```

For all available options, see the Prometheus
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/prometheus.conf).

### Time Series Selector (filtering)

To filter unwanted time series (metrics) use `selector` configuration option.

Here is an example:

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics
    # (allow[0] || allow[1] || ...) && !(deny[0] || deny[1] || ...)
    selector:
      allow:
        - <PATTERN>
        - <PATTERN>
      deny:
        - <PATTERN>
        - <PATTERN>
```

To find `PATTERN` syntax description and more examples
see [selectors readme](https://github.com/netdata/go.d.plugin/tree/master/pkg/prometheus/selector#time-series-selector).

### Time Series Grouping

It has built-in grouping logic based on the [type of metrics](https://prometheus.io/docs/concepts/metric_types/).

| Metric                    | Chart                                     | Dimension(s)         | Algorithm   |
|---------------------------|-------------------------------------------|----------------------|-------------|
| Gauge                     | for each label set                        | one, the metric name | absolute    |
| Counter                   | for each label set                        | one, the metric name | incremental |
| Summary (quantiles)       | for each label set (excluding 'quantile') | for each quantile    | absolute    |
| Summary (sum and count)   | for each label set                        | the metric name      | incremental |
| Histogram (buckets)       | for each label set (excluding 'le')       | for each bucket      | incremental |
| Histogram (sum and count) | for each label set                        | the metric name      | incremental |

Untyped metrics (have no '# TYPE') processing:

- As Counter or Gauge depending on pattern match when 'fallback_type' is used.
- As Counter if it has suffix '_total'.
- As Summary if it has 'quantile' label.
- As Histogram if it has 'le' label.

The rest are ignored.

## Troubleshooting

To troubleshoot issues with the `prometheus` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m prometheus
  ```
