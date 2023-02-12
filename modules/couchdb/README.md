<!--
title: "Apache CouchDB monitoring with Netdata"
description: "Monitor the health and performance of CouchDB databases with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/couchdb/README.md"
sidebar_label: "CouchDB"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# Apache CouchDB monitoring with Netdata

Monitors vital statistics of a local Apache CouchDB 2.x server.

This module monitors one or more `CouchDB` instances, depending on your configuration.

## Metrics

All metrics have "couchdb." prefix.

| Metric                | Scope  |                                                    Dimensions                                                     |    Units    |
|-----------------------|:------:|:-----------------------------------------------------------------------------------------------------------------:|:-----------:|
| activity              | global |                                          db_reads, db_writes, view_reads                                          | requests/s  |
| request_methods       | global |                                    copy, delete, get, head, options, post, put                                    | requests/s  |
| response_codes        | global | 200, 201, 202, 204, 206, 301, 302, 304, 400, 401, 403, 404, 406, 409, 412, 413, 414, 415, 416, 417, 500, 501, 503 | responses/s |
| response_code_classes | global |                                                2xx, 3xx, 4xx, 5xx                                                 | responses/s |
| active_tasks          | global |                               indexer, db_compaction, replication, view_compaction                                |    tasks    |
| replicator_jobs       | global |                               running, pending, crashed, internal_replication_jobs                                |    jobs     |
| open_files            | global |                                                       files                                                       |    files    |
| erlang_vm_memory      | global |                                      atom, binaries, code, ets, procs, other                                      |      B      |
| proccounts            | global |                                                os_procs, erl_procs                                                |  processes  |
| peakmsgqueue          | global |                                                     peak_size                                                     |  messages   |
| reductions            | global |                                                    reductions                                                     | reductions  |
| db_sizes_file         | global |                                          <i>a dimension per database</i>                                          |     KiB     |
| db_sizes_external     | global |                                          <i>a dimension per database</i>                                          |     KiB     |
| db_sizes_active       | global |                                          <i>a dimension per database</i>                                          |     KiB     |
| db_doc_count          | global |                                          <i>a dimension per database</i>                                          |    docs     |
| db_doc_del_count      | global |                                          <i>a dimension per database</i>                                          |    docs     |

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
