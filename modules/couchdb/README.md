<!--
title: "Apache CouchDB monitoring with Netdata"
description: "Monitor the health and performance of CouchDB databases with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/couchdb/README.md"
sidebar_label: "CouchDB"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# Apache CouchDB collector

Monitors vital statistics of a local Apache CouchDB 2.x server.

This module monitors one or more `CouchDB` instances, depending on your configuration.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/couchdb/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/couchdb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata   # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/couchdb.conf
```

Sample for a local server running on port 5984:

```yaml
local:
  user: 'admin'
  pass: 'password'
  node: 'couchdb@127.0.0.1'
```

Be sure to specify a correct admin-level username and password.

You may also need to change the `node` name; this should match the value of `-name NODENAME` in your
CouchDB's `etc/vm.args` file. Typically, this is of the form `couchdb@fully.qualified.domain.name` in a cluster,
or `couchdb@127.0.0.1` / `couchdb@localhost` for a single-node server.

If you want per-database statistics, these need to be added to the configuration, separated by spaces:

```yaml
local:
  databases: 'db1 db2 db3 ...'
```

For all available options, see the CouchDB
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/couchdb.conf).

## Troubleshooting

To troubleshoot issues with the `couchdb` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m couchdb
  ```
