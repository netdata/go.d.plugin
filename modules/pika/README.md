<!--
title: "Pika monitoring with Netdata"
description: "Monitor the health and performance of Pika storage services with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/pika/README.md"
sidebar_label: "Pika"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Storage"
-->

# Pika collector

[`Pika`](https://github.com/Qihoo360/pika#introduction%E4%B8%AD%E6%96%87) is a persistent huge storage service,
compatible with the vast majority of redis interfaces (details), including string, hash, list, zset, set and management
interfaces.

---

This module monitors one or more `Pika` instances, depending on your configuration.

It collects information and statistics about the server executing the following commands:

- [`INFO ALL`](https://github.com/Qihoo360/pika/wiki/pika-info%E4%BF%A1%E6%81%AF%E8%AF%B4%E6%98%8E)

## Metrics

All metrics have "pika." prefix.

| Metric                        | Scope  |           Dimensions            |    Units    |
|-------------------------------|:------:|:-------------------------------:|:-----------:|
| connections                   | global |            accepted             | connections |
| clients                       | global |            connected            |   clients   |
| memory                        | global |              used               |    bytes    |
| connected_replicas            | global |            connected            |  replicas   |
| commands                      | global |            processed            | commands/s  |
| commands_calls                | global | <i>a dimension per command</i>  |   calls/s   |
| database_strings_keys         | global | <i>a dimension per database</i> |    keys     |
| database_strings_expires_keys | global | <i>a dimension per database</i> |    keys     |
| database_strings_invalid_keys | global | <i>a dimension per database</i> |    keys     |
| database_hashes_keys          | global | <i>a dimension per database</i> |    keys     |
| database_hashes_expires_keys  | global | <i>a dimension per database</i> |    keys     |
| database_hashes_invalid_keys  | global | <i>a dimension per database</i> |    keys     |
| database_lists_keys           | global | <i>a dimension per database</i> |    keys     |
| database_lists_expires_keys   | global | <i>a dimension per database</i> |    keys     |
| database_lists_invalid_keys   | global | <i>a dimension per database</i> |    keys     |
| database_zsets_keys           | global | <i>a dimension per database</i> |    keys     |
| database_zsets_expires_keys   | global | <i>a dimension per database</i> |    keys     |
| database_zsets_invalid_keys   | global | <i>a dimension per database</i> |    keys     |
| database_sets_keys            | global | <i>a dimension per database</i> |    keys     |
| database_sets_expires_keys    | global | <i>a dimension per database</i> |    keys     |
| database_sets_invalid_keys    | global | <i>a dimension per database</i> |    keys     |
| uptime                        | global |             uptime              |   seconds   |

## Configuration

Edit the `go.d/pika.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/pika.conf
```

There are two connection types: by tcp socket and by unix socket.

```cmd
# by tcp socket
redis://<user>:<password>@<host>:<port>

# by unix socket
unix://<user>:<password>@</path/to/pika.sock
```

Needs only `address`, here is an example with two jobs:

```yaml
jobs:
  - name: local
    address: 'redis://@127.0.0.1:6379'

  - name: remote
    address: 'redis://user:password@203.0.113.0:6379'
```

For all available options, see the `pika`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/pika.conf).

## Troubleshooting

To troubleshoot issues with the `pika` collector, run the `go.d.plugin` with the debug option enabled. The output should
give you clues as to why the collector isn't working.
'
First, navigate to your plugins' directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m pika
```
