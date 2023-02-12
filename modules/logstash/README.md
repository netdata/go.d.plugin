<!--
title: "Logstash monitoring with Netdata"
description: "Monitor the health and performance of Logstash instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/logstash/README.md"
sidebar_label: "Logstash"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Apm"
-->

# Logstash monitoring with Netdata

[Logstash](https://www.elastic.co/products/logstash) is an open-source data processing pipeline that allows you to
collect, process, and load data into Elasticsearch.

This module will monitor one or more Logstash instances, depending on your configuration.

## Metrics

All metrics have "logstash." prefix.

Labels per scope:

- global: no labels.
- pipeline: pipeline.

| Metric                 |  Scope   |    Dimensions     |   Units    |
|------------------------|:--------:|:-----------------:|:----------:|
| jvm_threads            |  global  |      threads      |   count    |
| jvm_mem_heap_used      |  global  |      in_use       | percentage |
| jvm_mem_heap           |  global  |  committed, used  |    KiB     |
| jvm_mem_pools_eden     |  global  |  committed, used  |    KiB     |
| jvm_mem_pools_survivor |  global  |  committed, used  |    KiB     |
| jvm_mem_pools_old      |  global  |  committed, used  |    KiB     |
| jvm_gc_collector_count |  global  |     eden, old     |  counts/s  |
| jvm_gc_collector_time  |  global  |     eden, old     |     ms     |
| open_file_descriptors  |  global  |       open        |     fd     |
| event                  |  global  | in, filtered, out |  events/s  |
| event_duration         |  global  |   event, queue    |  seconds   |
| uptime                 |  global  |      uptime       |  seconds   |
| pipeline_event         | pipeline | in, filtered, out |  events/s  |
| pipeline_event         | pipeline |   event, queue    |  seconds   |

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
  ./go.d.plugin -d -m logstash
  ```
