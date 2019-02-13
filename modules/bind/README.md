# bind

This module will monitor one or more Bind(named) servers depending on configuration.

**Requirements:**
 * `bind` version 9.9+ with configured `statistics-channels`

It produces the following charts:

1. Received Requests by IP version (IPv4, IPv6)
2. Successful Queries
3. Recursive Clients
4. Queries by IP Protocol (TCP, UDP)
5. Queries Analysis
6. Received Updates
7. Query Failures
8. Query Failures Analysis
9. Server Statistics
10. Incoming Requests by OpCode
11. Incoming Requests by Query Type

Per View Statistics (the following set will be added for each bind view):
1. Resolver Active Queries
2. Resolver Statistics
3. Resolver Round Trip Timings
4. Resolver Requests by Query Type
4. Resolver Cache Hits

### configuration

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/bind.conf).
___

Needs only `url`.

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8653/json/v1

  - name: local
    url: http://127.0.0.1:8653/xml/v3
```

Without configuration, module will use `http://127.0.0.1:8653/json/v1`

**Views**: by default module doesn't collect views statistics.

To enable it please configure `permit_view`:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8653/json/v1
    permit_view: '!_* *'
```

Syntax: [simple patterns](https://docs.netdata.cloud/libnetdata/simple_pattern/).

### bind configuration
For detail information on how to get your bind installation ready, please refer to the [bind statistics channel developer comments](http://jpmens.net/2013/03/18/json-in-bind-9-s-statistics-server/) and to [bind documentation](https://ftp.isc.org/isc/bind/9.10.3/doc/arm/Bv9ARM.ch06.html#statistics) or [bind Knowledge Base article AA-01123](https://kb.isc.org/article/AA-01123/0).

Normally, you will need something like this in your `named.conf.options`:

```
statistics-channels {
        inet 127.0.0.1 port 8653 allow { 127.0.0.1; };
        inet ::1 port 8653 allow { ::1; };
};
```

(use the IPv4 or IPv6 line depending on what you are using, you can also use both)

Verify it works by running the following command:

```sh
curl "http://localhost:8653/json/v1/server"
```

---
