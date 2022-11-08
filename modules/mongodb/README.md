<!--
title: "MongoDB monitoring with Netdata"
description: "Monitor the health and performance of MongoDB with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/mongodb/README.md"
sidebar_label: "mongodb-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Databases"
-->

# MongoDB monitoring with Netdata

[`MongoDB`](https://www.mongodb.com/) MongoDB is a source-available cross-platform document-oriented database program.
Classified as a NoSQL database program, MongoDB uses JSON-like documents with optional schemas. MongoDB is developed by
MongoDB Inc. and licensed under the Server Side Public License (SSPL).

source: [`Wikipedia`](https://en.wikipedia.org/wiki/MongoDB)

---

This module monitors one or more `MongoDB` instances, depending on your configuration.

It collects information and statistics about the server executing the following commands:

- [`serverStatus`](https://docs.mongodb.com/manual/reference/command/serverStatus/#mongodb-dbcommand-dbcmd.serverStatus)
- [`dbStats`](https://docs.mongodb.com/manual/reference/command/dbStats/#dbstats)

## Prerequisites

Create a read-only user for Netdata in the admin database.

- Authenticate as the admin user.

  ```bash
  use admin
  db.auth("admin", "<MONGODB_ADMIN_PASSWORD>")
  ```

- Create a user.

  ```bash
  # MongoDB 2.x.
  db.addUser("netdata", "<UNIQUE_PASSWORD>", true)
  
  # MongoDB 3.x or higher.
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

- WireTiger metrics are available only if [WiredTiger](https://docs.mongodb.com/v5.0/core/wiredtiger/) is used as the
  storage engine.
- Sharding metris are available on shards only
  for [mongos](https://docs.mongodb.com/manual/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.process)

| Metric                        | Scope  |                                                                                                                                                                                          Dimensions                                                                                                                                                                                          |     Units      |
|-------------------------------|:------:|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| operations                    | global |                                                                                                                                                                       insert, query, update, delete, getmore, command                                                                                                                                                                        |     ops/s      |
| operations_latency            | global |                                                                                                                                                                                   reads, writes, commands                                                                                                                                                                                    |  milliseconds  |
| connections                   | global |                                                                                                                                                                                      current, available                                                                                                                                                                                      |  connections   |
| connections_rate              | global |                                                                                                                                                                                           created                                                                                                                                                                                            | connections/s  |
| connections_state             | global |                                                                                                                                                          active, threaded, exhaustIsMaster, exhaustHello, awaiting_topology_changes                                                                                                                                                          |  connections   |
| network_io                    | global |                                                                                                                                                                                           in, out                                                                                                                                                                                            |    bytes/s     |
| network_requests              | global |                                                                                                                                                                                           requests                                                                                                                                                                                           |   requests/s   |
| page_faults                   | global |                                                                                                                                                                                         page_faults                                                                                                                                                                                          | page_faults/s  |
| tcmalloc_generic              | global |                                                                                                                                                                                 current_allocated, heap_size                                                                                                                                                                                 |     bytes      |
| tcmalloc                      | global |                                                                                                                         pageheap_free, pageheap_unmapped, total_threaded_cache, free, pageheap_committed, pageheap_total_commit, pageheap_decommit, pageheap_reserve                                                                                                                         |     bytes      |
| asserts                       | global |                                                                                                                                                                       regular, warning, msg, user, tripwire, rollovers                                                                                                                                                                       |   asserts/s    |
| current_transactions          | global |                                                                                                                                                                               active, inactive, open, prepared                                                                                                                                                                               |  transactions  |
| shard_commit_types            | global |                                                                                                                no_shard_init, no_shard_successful, single_shard_init, single_shard_successful, shard_write_init, shard_write_successful, two_phase_init, two_phase_successful                                                                                                                |    commits     |
| active_clients                | global |                                                                                                                                                                                       readers, writers                                                                                                                                                                                       |    clients     |
| queued_operations             | global |                                                                                                                                                                                       readers, writers                                                                                                                                                                                       |   operation    |
| locks                         | global |                                                                                                                                                 global_read, global_write, database_read, database_write, collection_read, collection_write                                                                                                                                                  |   operation    |
| flow_control_timings          | global |                                                                                                                                                                                      acquiring, lagged                                                                                                                                                                                       |  milliseconds  |
| wiredtiger_blocks             | global |                                                                                                                                     read, read_via_memory_map_api, read_via_system_call_api, written, written_for_checkpoint, written_via_memory_map_api                                                                                                                                     |     bytes      |
| wiredtiger_cache              | global |                                                                                                                                                                  allocated_for_updates, read_into_cache, written_from_cache                                                                                                                                                                  |     bytes      |
| wiredtiger_capacity           | global |                                                                                                                                                    due_to_total_capacity, during_checkpoint, during_eviction, during_logging, during_read                                                                                                                                                    |      usec      |
| wiredtiger_connection         | global |                                                                                                                                                                   memory_allocations, memory_frees, memory_re_allocations                                                                                                                                                                    |     ops/s      |
| wiredtiger_cursor             | global | open_count, cached_count, bulk_loaded_insert_calls, close_calls_that_result_in_cache, create_calls, insert_calls, modify_calls, next_calls, operation_restarted, prev_calls, remove_calls, reserve_calls, cursor_reset_calls, search_calls, search_history_store_calls, search_near_calls, sweep_buckets, sweep_cursors_closed, sweep_cursors_examined, sweeps, truncate_calls, update_calls |    calls/s     |
| wiredtiger_lock               | global |                                                                                   checkpoint, dhandle_read, dhandle_write, durable_timestamp_queue_read, durable_timestamp_queue_write, metadata, read_timestamp_queue_read, read_timestamp_queue_write, schema, table_read, table_write, txn_global_read                                                                                    |     ops/s      |
| wiredtiger_lock_duration      | global |          checkpoint, checkpoint_internal_thread, dhandle_application_thread, dhandle_internal_thread, durable_timestamp_queue_application_thread, durable_timestamp_queue_internal_thread, metadata_application_thread, metadata_internal_thread, read_timestamp_queue_application_thread, read_timestamp_queue_internal_thread, schema_application_thread, schema_internal_thread           |   operation    |
| wiredtiger_log_ops            | global |                                                                                                                                                             flush, force_write, force_write_skipped, scan, sync, sync_dir, write                                                                                                                                                             |     ops/s      |
| wiredtiger_transactions       | global |                                                                                                                                              prepared, query_timestamp, rollback_to_stable, set_timestamp, begins, sync, committed, rolled back                                                                                                                                              | transactions/s |
| database_collections          | global |                                                                                                                                                                               <i>a dimension per database</i>                                                                                                                                                                                |  collections   |
| database_indexes              | global |                                                                                                                                                                               <i>a dimension per database</i>                                                                                                                                                                                |    indexes     |
| database_views                | global |                                                                                                                                                                               <i>a dimension per database</i>                                                                                                                                                                                |     views      |
| database_documents            | global |                                                                                                                                                                               <i>a dimension per database</i>                                                                                                                                                                                |   documents    |
| database_storage_size         | global |                                                                                                                                                                               <i>a dimension per database</i>                                                                                                                                                                                |     bytes      |
| replication_lag               | global |                                                                                                                                                                          <i>a dimension per replication member</i>                                                                                                                                                                           |  milliseconds  |
| replication_heartbeat_latency | global |                                                                                                                                                                          <i>a dimension per replication member</i>                                                                                                                                                                           |  milliseconds  |
| replication_node_ping         | global |                                                                                                                                                                          <i>a dimension per replication member</i>                                                                                                                                                                           |  milliseconds  |
| shard_nodes_count             | global |                                                                                                                                                                                  shard_aware, shard_unaware                                                                                                                                                                                  |     nodes      |
| shard_databases_status        | global |                                                                                                                                                                                 partitioned, un-partitioned                                                                                                                                                                                  |   databases    |
| chunks                        | global |                                                                                                                                                                                 <i>a dimension per shard</i>                                                                                                                                                                                 |     chunks     |

## Configuration

Edit the `go.d/mongodb.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata   # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/mongodb.conf
```

Sample using connection string:

**This is the preferred way**

```yaml
uri: 'mongodb://localhost:27017'
```

If no configuration is given, module will attempt to connect to mongodb daemon on `127.0.0.1:27017` address

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
