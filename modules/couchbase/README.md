<!--
title: "Couchbase monitoring with Netdata"
description: "Monitor the health and performance of Couchbase databases with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/couchbase/README.md"
sidebar_label: "Couchbase"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# Couchbase collector

[Couchbase](https://www.couchbase.com/) is an open source, distributed, JSON document database. It exposes a scale-out,
key-value store with managed cache for sub-millisecond data operations, purpose-built indexers for efficient queries and
a powerful query engine for executing SQL-like queries.

## Metrics

All metrics have "couchbase." prefix.

| Metric                            | Scope  |          Dimensions           |   Units    |
|-----------------------------------|:------:|:-----------------------------:|:----------:|
| bucket_quota_percent_used         | global | <i>a dimension per bucket</i> | percentage |
| bucket_ops_per_sec                | global | <i>a dimension per bucket</i> |   ops/s    |
| bucket_disk_fetches               | global | <i>a dimension per bucket</i> |  fetches   |
| bucket_item_count                 | global | <i>a dimension per bucket</i> |   items    |
| bucket_disk_used_stats            | global | <i>a dimension per bucket</i> |   bytes    |
| bucket_data_used                  | global | <i>a dimension per bucket</i> |   bytes    |
| bucket_mem_used                   | global | <i>a dimension per bucket</i> |   bytes    |
| bucket_vb_active_num_non_resident | global | <i>a dimension per bucket</i> |   items    |

## Charts

In this module 8 charts are supported because we collect basic stats from couchbase.

### Buckets Basic Stats

Collected from `/pools/default/buckets ` endpoint.

- Quota Percent Used Per Bucket in `%`
- Operations Per Second Per Bucket in `ops/s`
- Disk Fetches Per Bucket in `fetches`
- Item Count Per Bucket in `items`
- Disk Used Per Bucket in `bytes`
- Data Used Per Bucket in `bytes`
- Memory Used Per Bucket in `bytes`
- Number Of Non-Resident Items Per Bucket in `items`

## Configuration

Edit the `go.d/couchbase.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/couchbase.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8091
    username: admin
    password: admin-password

  - name: remote
    url: http://203.0.113.0:8091
    username: admin
    password: admin-passwor
```

For all available options, see the Couchbase
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/couchbase.conf).

## Troubleshooting

To troubleshoot issues with the `couchbase`, run the `go.d.plugin` with the debug option enabled. The output should give
you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m couchbase
  ```

