# dns_query

This module provides DNS query time in milliseconds.

It produces one chart:

1. **Query Time** in milliseconds
 * server1
 * server2
 ...

### configuration

Module specific options:
 * domains     - list of domains.
 * servers     - list of servers.
 * port        - server port. Default is 53.
 * network     - network transport. Default is upd. Supported options: udp, tcp, tcp-tls.
 * record_type - query record type. Default is A. Supported options: A, AAAA, CNAME, MX, NS, PTR, TXT, SOA, SPF, TXT, SRV.
 * timeout     - query read timeout. Default is 2 seconds.

Mandatory options: `domains` and `servers`. All other are optional.

Here is an example:

```yaml
jobs:
  - name: job1
    domains :
      - google.com
      - github.com
      - reddit.com
    servers:
      - 8.8.8.8
      - 8.8.4.4
```


Without configuration, module won't work.

---