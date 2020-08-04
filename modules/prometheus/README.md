<!--
title: "Prometheus endpoint monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/prometheus/README.md
sidebar_label: "Prometheus endpoints"
-->

# Prometheus endpoint monitoring with Netdata

This module collects metrics from one or more [`Prometheus`](https://prometheus.io/) endpoints.

## Charts

It produces one or more charts for every metric.

Number of charts depends on number of exposed metrics.
For example, scraping [`node_exporter`](https://github.com/prometheus/node_exporter) produces 3000+ metrics.


## Configuration

Edit the `go.d/prometheus.conf` configuration file using `edit-config` from the agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/prometheus.conf
```
 
Needs only `url`. Here is an example for 2 endpoints:

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9100/metrics

  - name: win10
    url: http://203.0.113.0:9182/metrics
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/prometheus.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m prometheus
