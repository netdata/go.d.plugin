<!--
title: "MongoDB monitoring with Netdata"
description: "Monitor the health and performance of MongoDB with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/mongodb/README.md"
sidebar_label: "mongodb-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# MongoDB collector

[MongoDB](https://www.mongodb.com/) is a source-available cross-platform document-oriented database program.

This module monitors one or more MongoDB instances, depending on your configuration. It collects information and
statistics about the server executing the following commands:

- [serverStatus](https://docs.mongodb.com/manual/reference/command/serverStatus/)
- [dbStats](https://docs.mongodb.com/manual/reference/command/dbStats/)
- [replSetGetStatus](https://www.mongodb.com/docs/manual/reference/command/replSetGetStatus/)

## Prerequisites

Create a read-only user for Netdata in the admin database.

- Authenticate as the admin user:

  ```bash
  use admin
  db.auth("admin", "<MONGODB_ADMIN_PASSWORD>")
  ```

- Create a user:

  ```bash
  db.createUser({
    "user":"netdata",
    "pwd": "<UNIQUE_PASSWORD>",
    "roles" : [
      {role: 'read', db: 'admin' },
      {role: 'clusterMonitor', db: 'admin'},
      {role: 'read', db: 'local' }
    ]
  })
  ```

## Metrics

- WireTiger metrics are available only if [WiredTiger](https://docs.mongodb.com/v6.0/core/wiredtiger/) is used as the
  storage engine.
- Sharding metrics are available on shards only
  for [mongos](https://www.mongodb.com/docs/manual/reference/program/mongos/).

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/mongodb/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/mongodb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata   # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/mongodb.conf
```

Needs only [connection URI](https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#connection-uri).
Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    uri: 'mongodb://user:password@localhost:27017'
    databases:
      include:
        - "* *"

  - name: remote
    dsn: 'mongodb://user:password@203.0.113.10:27017'
    databases:
      include:
        - "* *"
```

If no configuration is given, module will attempt to connect to mongodb daemon on `127.0.0.1:27017` address.

For all available options, see the `mongodb`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/mongodb.conf).

## Troubleshooting

To troubleshoot issues with the `mongodb` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m mongodb
  ```
