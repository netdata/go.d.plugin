# k8s_kubeproxy

This module will monitor one or more kube-proxy instances.


It produces the following charts:

1. **Sync Proxy Rules** in events/s
 * sync proxy rules

2. **Sync Proxy Rules Latency In Microseconds** observes per bucket
 * bucket 1000
 * bucket 2000
 * bucket 4000
 * bucket 8000
 * bucket 16000
 * bucket 32000
 * bucket 64000
 * bucket 128000
 * bucket 256000
 * bucket 512000
 * bucket 1024000
 * bucket 2048000
 * bucket 4096000
 * bucket 8192000
 * bucket 16384000
 * bucket +Inf
 
3. **HTTP Requests By Status Code** in requests/s
 * per code (200, 201, 404, etc.)
 
4. **HTTP Requests By Method** in requests/s
 * per code (GET, POST, etc.)

### configuration

Needs only `url` to kube-proxy metric-address.

Here is an example:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:10249/metrics
      
  - name: remote
    url : http://100.64.0.1:10249/metrics
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/k8s_kubeproxy.conf).

Without configuration, module attempts to connect to `http://127.0.0.1:10249/metrics`

---
