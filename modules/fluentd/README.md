# fluentd

This module will monitor one or more Fluentd servers depending on configuration.

The module gathers metrics from plugin endpoint provided by [in_monitor plugin](https://docs.fluentd.org/v1.0/articles/monitoring-rest-api).

**Requirements:**
 * `fluentd` with enabled monitoring agent

It produces the following charts:

1. **Plugin Retry Count**

2. **Plugin Buffer Queue Length**

3. **Plugin Buffer Total Size**

### configuration

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/fluentd.conf).
___

Needs only `url`.

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:24220

  - name: remote
    url: http://10.0.0.1:24220
```

Without configuration, module attempts to connect to `http://127.0.0.1:24220`.
___

**Filter plugins**: by default module collects statistics for all plugins.

To filter unwanted please configure `permit_plugin_id`:

```yaml
jobs:
  - name: local
    url: http://10.0.0.1:24220
    permit_plugin_id: '!monitor_agent !dummy *'
```

Syntax: [simple patterns](https://docs.netdata.cloud/libnetdata/simple_pattern/).

---
