<!--
title: "Solr monitoring with Netdata"
description: "Monitor the health and performance of Solr search servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/solr/README.md"
sidebar_label: "Solr"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Apps"
-->

# Solr collector

[`Solr`](https://lucene.apache.org/solr/) is an open-source enterprise-search platform, written in Java, from the Apache
Lucene project.

This module monitors `Solr` request handler statistics.

## Requirement

- `Solr` version 6.4+

## Metrics

All metrics have "solr." prefix.

| Metric                                     | Scope  |        Dimensions        |    Units     |
|--------------------------------------------|:------:|:------------------------:|:------------:|
| search_requests                            | global |          search          |  requests/s  |
| search_errors                              | global |          errors          |   errors/s   |
| search_errors_by_type                      | global | client, server, timeouts |   errors/s   |
| search_requests_processing_time            | global |           time           | milliseconds |
| search_requests_timings                    | global |  min, median, mean, max  | milliseconds |
| search_requests_processing_time_percentile | global |   p75, p95, p99, p999    | milliseconds |
| update_requests                            | global |          search          |  requests/s  |
| update_errors                              | global |          errors          |   errors/s   |
| update_errors_by_type                      | global | client, server, timeouts |   errors/s   |
| update_requests_processing_time            | global |           time           | milliseconds |
| update_requests_timings                    | global |  min, median, mean, max  | milliseconds |
| update_requests_processing_time_percentile | global |   p75, p95, p99, p999    | milliseconds |

## Configuration

Edit the `go.d/solr.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/solr.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://localhost:8983

  - name: remote
    url: http://203.0.113.10:8983

```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/solr.conf).

## Troubleshooting

To troubleshoot issues with the `solr` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m solr
  ```
