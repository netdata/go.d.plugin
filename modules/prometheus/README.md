<!--
title: "Prometheus endpoint monitoring with Netdata"
description: "Monitor the health and performance of 600+ services that support the Prometheus metrics with Netdata's per-second frequency and zero configuration."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/prometheus/README.md"
sidebar_label: "Prometheus endpoints"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Extended metric collection"
-->

# Prometheus endpoint monitoring with Netdata

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
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/prometheus.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

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

### Dimension algorithm

`incremental` algorithm (values displayed as rate) used when:

- the metric type is `Counter`, `Histogram` or `Summary`.
- the metrics suffix is `_total`, `_sum` or `_count`.

`absolute` algorithm (values displayed as is) is used in all other cases.

Use `force_absolute_algorithm` configuration option to overwrite the logic.

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics
    force_absolute_algorithm:
      - '*_sum'
      - '*_count'
```

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

This module groups time series into charts. It has built-in grouping logic (based on metric type). It is possible to
extend it via `group` configuration option.

#### Gauge and Counter

- A chart per every metric.
- Dimensions are labels sets.
- Dimensions per chart limit is `50`. If there is more dimensions the chart split into several charts.
- Values as is.

For instance, the following time series produce 1 chart.

```cmd
example_device_cur_state{name="0",type="Fan"} 0
example_device_cur_state{name="1",type="Fan"} 0
example_device_cur_state{name="10",type="Processor"} 0
example_device_cur_state{name="11",type="intel_powerclamp"} -1
example_device_cur_state{name="2",type="Fan"} 0
example_device_cur_state{name="3",type="Fan"} 0
example_device_cur_state{name="4",type="Fan"} 0
example_device_cur_state{name="5",type="Processor"} 0
example_device_cur_state{name="6",type="Processor"} 0
example_device_cur_state{name="7",type="Processor"} 0
example_device_cur_state{name="8",type="Processor"} 0
example_device_cur_state{name="9",type="Processor"} 0
```

#### Custom Grouping (Gauge and Counter only)

To group time series use `group` configuration option.

Here is an example:

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics
    group:
      - selector: <PATTERN>
        by_label: <a space separated list of labels names>
      - selector: <PATTERN>
        by_label: <a space separated list of labels names> 
```

To find `PATTERN` syntax description and more examples
see [selectors readme](https://github.com/netdata/go.d.plugin/tree/master/pkg/prometheus/selector#time-series-selector).

This example configuration groups all time series with metric names equal to `example_device_cur_state`
into multiple charts by `type` label. Number of charts is equal to number of `type` label values.

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics
    group:
      - selector: example_device_cur_state
        by_label: type 
```

#### Summary

- A chart per time series (label set).
- Dimensions are quantiles.
- Values as is.

For instance, the following time series produce 2 charts.

```cmd
example_duration_seconds{interval="15s",quantile="0"} 4.693e-06
example_duration_seconds{interval="15s",quantile="0.25"} 2.4383e-05
example_duration_seconds{interval="15s",quantile="0.5"} 0.00013458
example_duration_seconds{interval="15s",quantile="0.75"} 0.000195183
example_duration_seconds{interval="15s",quantile="1"} 0.005386229

example_duration_seconds{interval="30s",quantile="0"} 4.693e-06
example_duration_seconds{interval="30s",quantile="0.25"} 2.4383e-05
example_duration_seconds{interval="30s",quantile="0.5"} 0.00013458
example_duration_seconds{interval="30s",quantile="0.75"} 0.000195183
example_duration_seconds{interval="30s",quantile="1"} 0.005386229
```

#### Histogram

- A chart per time series (label set).
- Dimensions are `le` buckets.
- Values are not as is because histogram buckets are cumulative (`le="0.3"` contains `le="1.2"`). We calculate exact
  values for all buckets.

For instance, the following time series produce 2 charts.

```cmd
example_seconds_bucket{interval="15s",le="0.1"} 0
example_seconds_bucket{interval="15s",le="0.25"} 0
example_seconds_bucket{interval="15s",le="0.5"} 0
example_seconds_bucket{interval="15s",le="1"} 0
example_seconds_bucket{interval="15s",le="2.5"} 0
example_seconds_bucket{interval="15s",le="5"} 0
example_seconds_bucket{interval="15s",le="+Inf"} 0

example_seconds_bucket{interval="30s",le="0.1"} 0
example_seconds_bucket{interval="30s",le="0.25"} 0
example_seconds_bucket{interval="30s",le="0.5"} 0
example_seconds_bucket{interval="30s",le="1"} 0
example_seconds_bucket{interval="30s",le="2.5"} 0
example_seconds_bucket{interval="30s",le="5"} 0
example_seconds_bucket{interval="30s",le="+Inf"} 0
```

For all available options, see the Prometheus
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/prometheus.conf).

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
