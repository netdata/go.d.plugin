<!--
title: "CockroachDB monitoring with Netdata"
description: "Monitor the health and performance of CockroachDB databases with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/cockroachdb/README.md"
sidebar_label: "CockroachDB"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# CockroachDB collector

[`CockroachDB`](https://www.cockroachlabs.com/)  is the SQL database for building global, scalable cloud services that
survive disasters.

This module will monitor one or more `CockroachDB` databases, depending on your configuration.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/cockroachdb/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/cockroachdb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/cockroachdb.conf
```

Needs only `url` to server's `_status/vars`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8080/_status/vars

  - name: remote
    url: http://203.0.113.10:8080/_status/vars
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/cockroachdb.conf).

## Update every

Default `update_every` is 10 seconds because `CockroachDB` default sampling interval is 10 seconds, and it is not user
configurable. It doesn't make sense to decrease the value.

## Troubleshooting

To troubleshoot issues with the `cockroachdb` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m cockroachdb
  ```

