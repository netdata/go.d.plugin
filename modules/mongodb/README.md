<!--
title: "MongoDB monitoring with Netdata"
description: "Monitor the health and performance of MongoDB with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/mongodb/README.md
sidebar_label: "MongoDB"
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

## Charts

### Commands rate

- insert in `commands/s`
- query in `commands/s`
- update in `commands/s`
- delete in `commands/s`
- getmore in `commands/s`
- command in `commands/s`

### Active Clients

- readers in `clients`
- writers in `clients`

### Connections

- available in `connections/s`
- current in `connections/s`
- active in `connections/s`
- threaded in `connections/s`
- exhaustIsMaster in `connections/s`
- exhaustHello in `connections/s`
- awaiting topology changes in `connections/s`

### Memory

- resident in `MiB`
- virtual in `MiB`
- mapped in `MiB`
- mapped with journal in `MiB`

### Page faults

- Page Faults in `page faults/s`

### Asserts

- regular
- warning
- msg
- user
- tripwire
- rollovers

### Collections

- collections
- capped
- timeseries
- views
- internalCollections
- internalViews

### Tcmalloc generic metrics

- current_allocated_bytes in `MiB`
- heap_size in `MiB`

### Tcmalloc metrics

- Pageheap free in `KiB`
- Pageheap unmapped in `KiB`
- Total threaded cache in `KiB`
- Free in `KiB`
- Pageheap committed in `KiB`
- Pageheap total commit in `KiB`
- Pageheap decommit in `KiB`
- Pageheap reserve in `KiB`

### Network IO

- Bytes In in `bytes/s`
- Bytes Out in `bytes/s`

### Network Requests

- Requests in `requests/s`

### Current Queue Clients

- readers in `clients`
- writers in `clients`

### Command Metrics

- Eval in `commands`
- Eval Failed in `commands`
- Delete in `commands`
- Delete Failed in `commands`
- Count Failed in `commands`
- Create Indexes in `commands`
- Find And Modify in `commands`
- Insert Fail in `commands`

### Operations Latency

- Ops reads in `msec`
- Ops writes in `msec`
- disc in `msec`

### Current Transactions

- current active in `transactions`
- current inactive in `transactions`
- current open in `transactions`
- current prepared in `transactions`

### Global Locks

- Global Read Locks in `locks`
- Global Write Locks in `locks`
- Database Read Locks in `locks`
- Database Write Locks in `locks`
- Collection Read Locks in `locks`
- Collection Write Locks in `locks`

### Flow Control Stats

- timeAcquiringMicros in `number`
- isLaggedTimeMicros in `number`

### Wired Tiger Block Manager

- bytes read in `KiB`
- bytes read via memory map API in `KiB`
- bytes read via system call API in `KiB`
- bytes written in `KiB`
- bytes written for checkpoint in `KiB`
- bytes written via memory map API in `KiB`

- bytes allocated for updates in `KiB`
- bytes read into cache in `KiB`
- bytes written from cache in `KiB`

### Wired Tiger Capacity

- time waiting due to total capacity (usecs) in `usec`
- time waiting during checkpoint (usecs) in `usec`
- time waiting during eviction (usecs) in `usec`
- time waiting during logging (usecs) in `usec`
- time waiting during read (usecs) in `usec`

### Wired Tiger Connections

- memory allocations in `ops`
- memory frees in `ops`
- memory re-allocations in `ops`

### Wired Tiger Cursor

- open cursor count in `calls`
- cached cursor count in `calls`
- cursor bulk loaded cursor insert calls in `calls`
- cursor close calls that result in cache in `calls`
- cursor create calls in `calls`
- cursor insert calls in `calls`
- cursor modify calls in `calls`
- cursor next calls in `calls`
- cursor operation restarted in `calls`
- cursor prev calls in `calls`
- cursor remove calls in `calls`
- cursor remove key bytes removed in `calls`
- cursor reserve calls in `calls`
- cursor reset calls in `calls`
- cursor search calls in `calls`
- cursor search history store calls in `calls`
- cursor search near calls in `calls`
- cursor sweep buckets in `calls`
- cursor sweep cursors closed in `calls`
- cursor sweep cursors examined in `calls`
- cursor sweeps in `calls`
- cursor truncate calls in `calls`
- cursor update calls in `calls`
- cursor update value size change in `calls`

### Wired Tiger Lock

- checkpoint lock acquisitions in `ops`
- dhandle read lock acquisitions in `ops`
- dhandle write lock acquisitions in `ops`
- durable timestamp queue read lock acquisitions in `ops`
- durable timestamp queue write lock acquisitions in `ops`
- metadata lock acquisitions in `ops`
- read timestamp queue read lock acquisitions in `ops`
- read timestamp queue write lock acquisitions in `ops`
- schema lock acquisitions in `ops`
- table read lock acquisitions in `ops`
- table write lock acquisitions in `ops`
- txn global read lock acquisitions in `ops`

### Wired Tiger Lock Duration

- checkpoint lock application thread wait time (usecs) in `usec`
- checkpoint lock internal thread wait time (usecs) in `usec`
- dhandle lock application thread time waiting (usecs) in `usec`
- dhandle lock internal thread time waiting (usecs) in `usec`
- durable timestamp queue lock application thread time waiting (usecs) in `usec`
- durable timestamp queue lock internal thread time waiting (usecs) in `usec`
- metadata lock application thread wait time (usecs) in `usec`
- metadata lock internal thread wait time (usecs) in `usec`
- read timestamp queue lock application thread time waiting (usecs) in `usec`
- read timestamp queue lock internal thread time waiting (usecs) in `usec`
- schema lock application thread wait time (usecs) in `usec`
- schema lock internal thread wait time (usecs) in `usec`
- table lock application thread time waiting for the table lock (usecs) in `usec`
- table lock internal thread time waiting for the table lock (usecs) in `usec`
- txn global lock application thread time waiting (usecs) in `usec`
- txn global lock internal thread time waiting (usecs) in `usec`

### Wired Tiger Log Operations

- log flush operations in `ops`
- log force write operations in `ops`
- log force write operations skipped in `ops`
- log scan operations in `ops`
- log sync operations in `ops`
- log sync_dir operations in `ops`
- log write operations in `ops`

### Wired Tiger Log Operations

- log bytes of payload data in `bytes`
- log bytes written in `bytes`
- logging bytes consolidated in `bytes`
- total log buffer size in `bytes`

### Wired Tiger Log Transactions

- prepared transactions in `transactions`
- query timestamp calls in `transactions`
- rollback to stable calls in `transactions`
- set timestamp calls in `transactions`
- transaction begins in `transactions`
- transaction sync calls in `transactions`
- transactions committed in `transactions`
- transactions rolled back in `transactions`

## Prerequisites

Create a read-only user for Netdata in the admin database.

1. Authenticate as the admin user.

```
use admin
db.auth("admin", "<MONGODB_ADMIN_PASSWORD>")
```

2. Create a user.

```
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

## Configuration

Edit the `go.d/mongodb.conf` configuration file using `edit-config` from the
Netdata [config directory](/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata   # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/mongodb.conf
```

Sample using connection string:

**This is the preferred way**

```yaml
connectionStr: 'mongodb://mongodb0.example.com:27017'
```

Sample using local database without authentication:

```yaml
local:
  name: 'admin'
  host: '127.0.0.1'
  port: 27017
```

Sample using authentication:

```yaml
local:
  name: 'admin'
  authdb: 'admin'
  host: '127.0.0.1'
  port: 27017
  user: 'netdata'
  pass: '<password>'
```

If no configuration is given, module will attempt to connect to mongodb daemon on `127.0.0.1:27017` address

For all available options, see the `mongodb`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/mongodb.conf).

## Troubleshooting

To troubleshoot issues with the `mongodb` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m mongodb
```

