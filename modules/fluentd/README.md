# fluentd

[`Fluentd`](https://www.fluentd.org/) is an open source data collector for unified logging layer.

This module will monitor one or more `Fluentd` servers depending on configuration. It gathers metrics from plugin endpoint provided by [in_monitor plugin](https://docs.fluentd.org/v1.0/articles/monitoring-rest-api).

## Requirements
-   `fluentd` with enabled monitoring agent

## Charts

It produces the following charts:

-   Plugin Retry Count in `count`

-   Plugin Buffer Queue Length in `queue length`

-   Plugin Buffer Total Size in `buffer`

## Configuration

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
