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

### Default charts

Works with local and cloud hosted [`Atlas`](https://www.mongodb.com/cloud/atlas) database servers

#### Operations by type

- insert in `ops/s`
- query in `ops/s`
- update in `ops/s`
- delete in `ops/s`
- getmore in `ops/s`
- command in `ops/s`

#### Operations Latency

- ops reads in `milliseconds`
- ops writes in `milliseconds`
- commands in `milliseconds`

#### Connections

- available in `connections`
- current in `connections`

#### Connections Rate

- created in `connections/s`

#### Connections By State

- active in `connections`
- threaded in `connections`
- exhaustIsMaster in `connections`
- exhaustHello in `connections`
- awaiting topology changes in `connections`

#### Network IO

- bytes in in `bytes/s`
- bytes out in `bytes/s`

#### Network Requests

- Requests in `requests/s`

#### Memory

- resident in `MiB`
- virtual in `MiB`
- mapped in `MiB`
- mapped with journal in `MiB`

#### Page faults

- Page Faults in `page faults/s`

#### Raised assertions

- regular in `asserts/s`
- warning in `asserts/s`
- msg in `asserts/s`
- user in `asserts/s`
- tripwire in `asserts/s`
- rollovers in `asserts/s`

#### Current Transactions

- current active in `transactions`
- current inactive in `transactions`
- current open in `transactions`
- current prepared in `transactions`

### Optional charts:

Depending on the database server version and configuration Mongo reports slightly different statistics. We use
the [`serverStatus`](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#serverstatus)
command to monitor the database. Based on the command output the following may be included:

#### Active Clients

- readers in `clients`
- writers in `clients`

if serverStatus
reports [global locks active clients](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#globallock)

#### Queued operations because of a lock

- readers in `clients`
- writers in `clients`

if serverStatus
reports [global locks current queue clients](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#globallock)

#### Tcmalloc generic metrics

- current_allocated_bytes in `bytes`
- heap_size in `bytes`

if serverStatus reports tcmalloc stats

#### Tcmalloc metrics

- pageheap free in `bytes`
- pageheap unmapped in `bytes`
- total threaded cache in `bytes`
- free in `bytes`
- pageheap committed in `bytes`
- pageheap total commit in `bytes`
- pageheap decommit in `bytes`
- pageheap reserve in `bytes`

if serverStatus reports tcmalloc stats

#### Acquired locks

- global read locks in `locks/s`
- global write locks in `locks/s`
- database read locks in `locks/s`
- database write locks in `locks/s`
- collection read locks in `locks/s`
- collection write locks in `locks/s`

#### Flow Control Stats

- acquiring time in `milliseconds`
- lagged time in `milliseconds`

if serverStatus reports [flow control](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#flowcontrol)

### Wired Tiger Charts

Available only if [WiredTiger](https://docs.mongodb.com/v5.0/core/wiredtiger/)
is used as the storage engine.

#### [Wired Tiger Block Manager](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.block-manager)

- bytes read in `bytes`
- bytes read via memory map API in `bytes`
- bytes read via system call API in `bytes`
- bytes written in `bytes`
- bytes written for checkpoint in `bytes`
- bytes written via memory map API in `bytes`

#### [Wired Tiger Cache](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.cache)

- bytes allocated for updates in `bytes`
- bytes read into cache in `bytes`
- bytes written from cache in `bytes`

#### Wired Tiger Capacity

- time waiting due to total capacity in `usec`
- time waiting during checkpoint in `usec`
- time waiting during eviction in `usec`
- time waiting during logging in `usec`
- time waiting during read in `usec`

#### [Wired Tiger Connection](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.connection)

- memory allocations in `ops/s`
- memory frees in `ops/s`
- memory re-allocations in `ops/s`

#### [Wired Tiger Cursor](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.cursor)

- open cursor count in `calls/s`
- cached cursor count in `calls/s`
- cursor bulk loaded cursor insert calls in `calls/s`
- cursor close calls that result in cache in `calls/s`
- cursor create calls in `calls/s`
- cursor insert calls in `calls/s`
- cursor modify calls in `calls/s`
- cursor next calls in `calls/s`
- cursor operation restarted in `calls/s`
- cursor prev calls in `calls/s`
- cursor remove calls in `calls/s`
- cursor remove key bytes removed in `calls/s`
- cursor reserve calls in `calls/s`
- cursor reset calls in `calls/s`
- cursor search calls in `calls/s`
- cursor search history store calls in `calls/s`
- cursor search near calls in `calls/s`
- cursor sweep buckets in `calls/s`
- cursor sweep cursors closed in `calls/s`
- cursor sweep cursors examined in `calls/s`
- cursor sweeps in `calls/s`
- cursor truncate calls in `calls/s`
- cursor update calls in `calls/s`
- cursor update value size change in `calls/s`

#### Wired Tiger Lock

- checkpoint lock acquisitions in `ops/s`
- dhandle read lock acquisitions in `ops/s`
- dhandle write lock acquisitions in `ops/s`
- durable timestamp queue read lock acquisitions in `ops/s`
- durable timestamp queue write lock acquisitions in `ops/s`
- metadata lock acquisitions in `ops/s`
- read timestamp queue read lock acquisitions in `ops/s`
- read timestamp queue write lock acquisitions in `ops/s`
- schema lock acquisitions in `ops/s`
- table read lock acquisitions in `ops/s`
- table write lock acquisitions in `ops/s`
- txn global read lock acquisitions in `ops/s`

#### Wired Tiger Lock Duration

- checkpoint lock application thread wait time in `usec`
- checkpoint lock internal thread wait time in `usec`
- dhandle lock application thread time waiting in `usec`
- dhandle lock internal thread time waiting in `usec`
- durable timestamp queue lock application thread time waiting in `usec`
- durable timestamp queue lock internal thread time waiting in `usec`
- metadata lock application thread wait time in `usec`
- metadata lock internal thread wait time in `usec`
- read timestamp queue lock application thread time waiting in `usec`
- read timestamp queue lock internal thread time waiting in `usec`
- schema lock application thread wait time in `usec`
- schema lock internal thread wait time in `usec`
- table lock application thread time waiting for the table lock in `usec`
- table lock internal thread time waiting for the table lock in `usec`
- txn global lock application thread time waiting in `usec`
- txn global lock internal thread time waiting in `usec`

#### [Wired Tiger Log Operations](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.log)

- log flush operations in `ops/s`
- log force write operations in `ops/s`
- log force write operations skipped in `ops/s`
- log scan operations in `ops/s`
- log sync operations in `ops/s`
- log sync_dir operations in `ops/s`
- log write operations in `ops/s`

#### [Wired Tiger Log Operations IO](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.log)

- log bytes of payload data in `bytes/s`
- log bytes written in `bytes/s`
- logging bytes consolidated in `bytes/s`
- total log buffer size in `bytes/s`

#### [Wired Tiger Log Transactions](https://docs.mongodb.com/v5.0/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.wiredTiger.log)

- prepared transactions in `transactions/s`
- query timestamp calls in `transactions/s`
- rollback to stable calls in `transactions/s`
- set timestamp calls in `transactions/s`
- transaction begins in `transactions/s`
- transaction sync calls in `transactions/s`
- transactions committed in `transactions/s`
- transactions rolled back in `transactions/s`

### [Database Statistics](https://docs.mongodb.com/manual/reference/command/dbStats/#dbstats)

Stack chart that include per database metrics for collections, indexes, views, documents and disk size. Because
internally the module has to run the
`dbstats` commands once per database and each server may have a large number of databases, these charts are disabled by
default. To enabled them please use
a [matcher](https://learn.netdata.cloud/docs/agent/collectors/go.d.plugin/pkg/matcher)
in the configuration file to enable the collection for specific or all databases.

For example, to exclude the admin and config databases use:

```yaml
 databases:
   includes:
     - "* *"
   excludes:
     - "= config"
     - "= admin"
```

#### [Collections](https://docs.mongodb.com/manual/reference/command/dbStats/#mongodb-data-dbStats.collections)

Number of collections per database in `collections`

#### [Indexes](https://docs.mongodb.com/manual/reference/command/dbStats/#mongodb-data-dbStats.indexes)

Number of indexes per database in `indexes`

#### [Views](https://docs.mongodb.com/manual/reference/command/dbStats/#mongodb-data-dbStats.views)

Number of views per database in `views`

#### [Documents](https://docs.mongodb.com/manual/reference/command/dbStats/#mongodb-data-dbStats.objects)

Number of documents per database in `documents`

#### [Disk Size](https://docs.mongodb.com/manual/reference/command/dbStats/#mongodb-data-dbStats.dataSize)

Total space allocated to all collections in the database for document storage(including free space) in `bytes`

### Replication

Charts are available on replica sets only.

#### Replication lag

Time when last entry from the oplog was applied for every member in `milliseconds`

#### Replication set member heartbeat latency

Time when last heartbeat was received from replica set member `milliseconds`

#### Ping

Round trip time between members and local instance in `milliseconds`

### Sharding

Charts are available on shards only for [mongos](https://docs.mongodb.com/manual/reference/command/serverStatus/#mongodb-serverstatus-serverstatus.process)

#### Transaction Commit Types

Number of initiated and successful `no shard`, `single shard`, `shard write` ,`two phase`commits in `commits`

#### Shard Nodes

Number of shard aware vs shard unaware nodes in `nodes`

#### Databases Sharding Status

Partitioned vs UnPartitioned databases in `databases`

#### Collections Sharding Status

Partitioned vs UnPartitioned collections in `databases`

#### Chucks Per Node

Number of chucks per node in `chunks`

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
uri: 'mongodb://localhost:27017'
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
