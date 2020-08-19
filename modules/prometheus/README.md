<!--
title: "Prometheus endpoint monitoring with Netdata"
description: "Monitor 600+ services that support the Prometheus/OpenMetrics exposition format with Netdata's per-second frequency and zero configuration."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/prometheus/README.md
sidebar_label: "Prometheus endpoints"
-->

# Prometheus endpoint monitoring with Netdata

The generic Prometheus endpoint collector gathers metrics from [`Prometheus`](https://prometheus.io/) endpoints that use
the [OpenMetrics exposition format](https://prometheus.io/docs/instrumenting/exposition_formats/).

As of v1.24, Netdata can autodetect more than 600 Prometheus endpoints, including support for Windows 10 via
`windows_exporter`, and instantly generate new charts with the same high-granularity, per-second frequency as you expect
from other collectors. 

The full list of endpoints is available in the collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/prometheus.conf).

## Charts

Netdata will produce one or more charts for every metric collected via a Prometheus endpoint. The number of charts
depends entirely on the number of exposed metrics.

For example, scraping [`node_exporter`](https://github.com/prometheus/node_exporter) produces 3000+ metrics.

## Configuration

Edit the `go.d/prometheus.conf` configuration file using `edit-config` from the Agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

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

### Time Series Selector (filtering)

To filter unwanted time series (metrics) use `selector` configuration option.

Here is an example:
```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics
    selector:
      allow:
        - <SELECTOR_PATTERN>
        - <SELECTOR_PATTERN>
      deny:
        - <SELECTOR_PATTERN>
        - <SELECTOR_PATTERN>
```

To find `SELECTOR_PATTERN` syntax description and more examples see [selectors readme](https://github.com/netdata/go.d.plugin/pkg/prometheus/selector#time-series-selectors).

For all available options, see the Prometheus collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/prometheus.conf).


## Troubleshooting

To troubleshoot issues with the Prometheus collector, run the `go.d.plugin` orchestrator with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugins directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m prometheus
```
