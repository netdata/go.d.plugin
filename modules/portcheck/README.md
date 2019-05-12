# portcheck

This module will monitor one or more TCP services for availability and response time.

It produces the following charts for every monitoring port:

1. **TCP Check Status** in boolean
 * success
 * failed
 * timeout

2. **Current State Duration** in seconds
 * time
 
3. **TCP Connection Latency** in ms
 * time

### configuration
 
Here is an example for 2 servers:

```yaml
jobs:
  - name: server1
    ports: [22, 23]
      
  - name: server2
    ports: [80, 81, 8080]
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/portcheck.conf).

---
