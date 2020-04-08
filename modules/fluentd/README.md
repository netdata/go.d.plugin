# Fluentd monitoring with Netdata

[`Fluentd`](https://www.fluentd.org/) is an open source data collector for unified logging layer.

This module will monitor one or more `Fluentd` servers, depending on your configuration. It gathers metrics from plugin endpoint provided by [in_monitor plugin](https://docs.fluentd.org/v1.0/articles/monitoring-rest-api).

## Requirements

-   `fluentd` with enabled monitoring agent

## Charts

It produces the following charts:

-   Plugin Retry Count in `count`
-   Plugin Buffer Queue Length in `queue length`
-   Plugin Buffer Total Size in `buffer`

## Configuration

Edit the `go.d/fluentd.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/fluentd.conf
```

Needs only `url`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:24220

  - name: local_with_filtering
    url: http://127.0.0.1:24220
    permit_plugin_id: '!monitor_agent !dummy *'

  - name: remote
    url: http://203.0.113.10:24220
```

By default this module collects statistics for all plugins. Filter plugins syntax: [simple patterns](https://docs.netdata.cloud/libnetdata/simple_pattern/).

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/fluentd.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m fluentd
