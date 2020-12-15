<!--
title: "Couchbase monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/couchbase/README.md
sidebar_label: "couchbase"
-->

# Couchbase monitoring with Netdata

Couchbase Server is an open source, distributed, JSON document database. It exposes a scale-out, key-value store with managed cache for sub-millisecond data operations, purpose-built indexers for efficient queries and a powerful query engine for executing SQL-like queries.

## Configuration

Edit the `go.d/couchbase.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/couchbase.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with one endpoints:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8091
    username: admin
    password: admin-password

```

For all available options, see the Couchbase
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/couchbase.conf).



## Troubleshooting

To troubleshoot issues with the `couchbase`, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m couchbase
```
