<!--
title: "MongoDB monitoring with Netdata"
description: "Monitor the health and performance of MongoDB with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/mongodb/README.md"
sidebar_label: "mongodb-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/Databases"
-->

# MongoDB monitoring with Netdata

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

All metrics have "mongodb." prefix.

- WireTiger metrics are available only if [WiredTiger](https://docs.mongodb.com/v6.0/core/wiredtiger/) is used as the
  storage engine.
- Sharding metrics are available on shards only
  for [mongos](https://www.mongodb.com/docs/manual/reference/program/mongos/).

Labels per scope:

- global: no labels.
- lock type: lock_type.
- commit type: commit_type.
- database: database.
- replica set member: repl_set_member.
- shard: shard_id.

| Metric                                         |       Scope        |                                                        Dimensions                                                        |     Units      |
|------------------------------------------------|:------------------:|:------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| operations_rate                                |       global       |                                                 reads, writes, commands                                                  |  operations/s  |
| operations_latency_time                        |       global       |                                                 reads, writes, commands                                                  |  milliseconds  |
| operations_by_type_rate                        |       global       |                                     insert, query, update, delete, getmore, command                                      |  operations/s  |
| document_operations_rate                       |       global       |                                           inserted, deleted, returned, updated                                           |  operations/s  |
| scanned_indexes_rate                           |       global       |                                                         scanned                                                          |   indexes/s    |
| scanned_documents_rate                         |       global       |                                                         scanned                                                          |  documents/s   |
| active_clients_count                           |       global       |                                                      reads, writes                                                       |    clients     |
| queued_operations_count                        |       global       |                                                      reads, writes                                                       |   operations   |
| lock_acquisitions_rate                         |     lock type      |                                    shared, exclusive, intent_shared, intent_exclusive                                    | acquisitions/s |
| cursors_open_count                             |       global       |                                                           open                                                           |    cursors     |
| cursors_open_no_timeout_count                  |       global       |                                                     open_no_timeout                                                      |    cursors     |
| cursors_opened_rate                            |       global       |                                                          opened                                                          |   cursors/s    |
| cursors_timed_out_rate                         |       global       |                                                        timed_out                                                         |   cursors/s    |
| cursors_by_lifespan_count                      |       global       |                                  le_1s, 1s_5s, 5s_15s, 15s_30s, 30s_1m, 1m_10m, ge_10m                                   |    cursors     |
| transactions_count                             |       global       |                                             active, inactive, open, prepared                                             |  transactions  |
| transactions_rate                              |       global       |                                          started, aborted, committed, prepared                                           | transactions/s |
| transactions_commits_rate                      |    commit type     |                                                      success, fail                                                       |   commits/s    |
| transactions_commits_duration_time             |    commit type     |                                                         commits                                                          |  milliseconds  |
| connections_usage                              |       global       |                                                     available, used                                                      |  connections   |
| connections_by_state_count                     |       global       |                      active, threaded, exhaust_is_master, exhaust_hello, awaiting_topology_changes                       |  connections   |
| connections_rate                               |       global       |                                                         created                                                          | connections/s  |
| asserts_rate                                   |       global       |                                     regular, warning, msg, user, tripwire, rollovers                                     |   asserts/s    |
| network_traffic_rate                           |       global       |                                                         in, out                                                          |    bytes/s     |
| network_requests_rate                          |       global       |                                                         requests                                                         |   requests/s   |
| network_slow_dns_resolutions_rate              |       global       |                                                         slow_dns                                                         | resolutions/s  |
| network_slow_ssl_handshakes_rate               |       global       |                                                         slow_ssl                                                         |  handshakes/s  |
| memory_resident_size                           |       global       |                                                           used                                                           |     bytes      |
| memory_virtual_size                            |       global       |                                                           used                                                           |     bytes      |
| memory_page_faults_rate                        |       global       |                                                         pgfaults                                                         |   pgfaults/s   |
| memory_tcmalloc_stats                          |       global       | allocated, central_cache_freelist, transfer_cache_freelist, thread_cache_freelists, pageheap_freelist, pageheap_unmapped |     bytes      |
| wiredtiger_concurrent_read_transactions_usage  |       global       |                                                     available, used                                                      |  transactions  |
| wiredtiger_concurrent_write_transactions_usage |       global       |                                                     available, used                                                      |  transactions  |
| wiredtiger_cache_usage                         |       global       |                                                           used                                                           |     bytes      |
| wiredtiger_cache_dirty_space_size              |       global       |                                                          dirty                                                           |     bytes      |
| wiredtiger_cache_io_rate                       |       global       |                                                      read, written                                                       |    pages/s     |
| wiredtiger_cache_evictions_rate                |       global       |                                                   unmodified, modified                                                   |    pages/s     |
| database_collection_count                      |      database      |                                                       collections                                                        |  collections   |
| database_indexes_count                         |      database      |                                                         indexes                                                          |    indexes     |
| database_views_count                           |      database      |                                                          views                                                           |     views      |
| database_documents_count                       |      database      |                                                        documents                                                         |   documents    |
| database_data_size                             |      database      |                                                        data_size                                                         |     bytes      |
| database_storage_size                          |      database      |                                                       storage_size                                                       |     bytes      |
| database_index_size                            |      database      |                                                        index_size                                                        |     bytes      |
| repl_set_member_state                          | replica set member |               primary, startup, secondary, recovering, startup2, unknown, arbiter, down, rollback, removed               |     state      |
| repl_set_member_health_status                  | replica set member |                                                         up, down                                                         |     status     |
| repl_set_member_replication_lag_time           | replica set member |                                                     replication_lag                                                      |  milliseconds  |
| repl_set_member_heartbeat_latency_time         | replica set member |                                                    heartbeat_latency                                                     |  milliseconds  |
| repl_set_member_ping_rtt_time                  | replica set member |                                                         ping_rtt                                                         |  milliseconds  |
| repl_set_member_uptime                         | replica set member |                                                          uptime                                                          |    seconds     |
| sharding_nodes_count                           |       global       |                                                shard_aware, shard_unaware                                                |     nodes      |
| sharding_sharded_databases_count               |       global       |                                                partitioned, unpartitioned                                                |   databases    |
| sharding_sharded_collections_count             |       global       |                                                partitioned, unpartitioned                                                |  collections   |
| sharding_shard_chunks_count                    |       shard        |                                                          chunks                                                          |     chunks     |

## Configuration

Edit the `go.d/mongodb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

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
