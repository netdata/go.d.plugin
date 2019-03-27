# kubelet

This module will monitor one or more kubelet instances.

It produces the following charts (if all data is available):

1. **API Server Audit Requests** in requests
 * rejected

2. **API Server Failed Data Encryption Key(DEK) Generation Operations** in events/s
 * failures
 
3. **API Server Latencies Of Data Encryption Key(DEK) Generation Operations** in observes/s
 * per bucket (5 µs, 10 µs, ..., 40960 µs, +Inf)
 
4. **API Server Latencies Of Data Encryption Key(DEK) Generation Operations Percentage** in %
 * per bucket (5 µs, 10 µs, ..., 40960 µs, +Inf)
 
5. **API Server Storage Envelope Transformation Cache Misses** in events/s
 * cache misses
 
6. **Number Of Containers Currently Running** in running containers
 * total
 
7. **Number Of Pods Currently Running** in running pods
 * total
 
7. **Runtime Operations By Type** in operations/s
 * per operation type
 
8. **Docker Operations By Type** in operations/s
 * per operation type
 
9. **Docker Operations Errors By Type** in operations/s
 * per operation error type
 
10. **Node Configuration-Related Error** in bool
 * experiencing error
 
11. **PLEG Relisting Interval Summary** in microseconds per quantile
 * 0.5
 * 0.9
 * 0.99
 
12. **PLEG Relisting Latency Summary** in microseconds per quantile
 * 0.5
 * 0.9
 * 0.99
 
13. **Token() Requests To The Alternate Token Source** in token requests/s
 * total
 * failed
 
14. **REST Client HTTP Requests By Status Code** in requests/s
 * per status code (200, 201, etc.)
 
15. **REST Client HTTP Requests By Method** in requests/s
 * per http method (GET, POST, etc.)
 
16. **Volume Manager State Of The World** per every plugin in state of the world
 * actual
 * desired
 

### configuration

Needs only `url` to kubelet metric-address.

Here is an example:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:10255/metrics
      
  - name: remote
    url : http://100.64.0.1:10255/metrics
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/k8s_kubelet.conf).

Without configuration, module attempts to connect to `http://127.0.0.1:10255/metrics`

---
