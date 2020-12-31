<!--
title: "Apache CouchDB monitoring with Netdata"
description: "Monitor the health and performance of CouchDB databases with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/couchdb/README.md
sidebar_label: "CouchDB"
-->

# Apache CouchDB monitoring with Netdata

Monitors vital statistics of a local Apache CouchDB 2.x server, including:

- Overall server reads/writes
- HTTP traffic breakdown
    - Request methods (`GET`, `PUT`, `POST`, etc.)
    - Success response status codes (`200`, `201`, etc.)
    - Response status code classes (`2xx`, `3xx`, etc.)
- Active server tasks
- Replication status (CouchDB 2.1 and up only)
- Erlang VM stats
- Optional per-database statistics: sizes, # of docs, # of deleted docs

This module monitors one or more `CouchDB` instances, depending on your configuration.

## Configuration

Edit the `go.d/couchdb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

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
  ...
  databases: 'db1 db2 db3 ...'
```

For all available options, see the CouchDB
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/couchdb.conf).

## Troubleshooting

To troubleshoot issues with the `couchdb` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m couchdb
```
