<!--
title: "Logstash monitoring with Netdata"
description: "Monitor the health and performance of Logstash instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/logstash/README.md
sidebar_label: "Logstash"
-->

# Logstash monitoring with Netdata

[`Logstash`](https://www.elastic.co/products/logstash) is an open-source data processing pipeline that allows you to
collect, process, and load data into `Elasticsearch`.

This module will monitor one or more `Logstash` instances, depending on your configuration.

## Charts

It produces following charts:

- JVM Threads in `count`
- JVM Heap Memory Percentage in `percent`
- JVM Heap Memory in `KiB`
- JVM Pool Survivor Memory in `KiB`
- JVM Pool Old Memory in `KiB`
- JVM Pool Eden Memory in `KiB`
- Garbage Collection Count in `counts/s`
- Time Spent On Garbage Collection in `ms`
- Uptime in `time`

## Configuration

Edit the `go.d/logstash.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/logstash.conf
```

Here is a simple example for local and remote server:

```yaml
jobs:
  - name: local
    url: http://localhost:9600

  - name: remote
    url: http://203.0.113.10:9600
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/logstash.conf).

## Troubleshooting

To troubleshoot issues with the `logstash` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m logstash
```
