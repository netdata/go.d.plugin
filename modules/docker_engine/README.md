# docker_engine

This module will monitor one or more Docker engines.

**Requirements:**
 * docker with enabled [`metric-address`](https://docs.docker.com/config/thirdparty/prometheus/)


It produces the following charts:

1. **Container Actions** in actions/s
 * changes
 * commits
 * create
 * delete
 * start

2. **Container States** in number of containers in state
 * running
 * paused
 * stopped

### configuration

Needs only `url` to docker metric-address.

Here is an example for 2 docker instances:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:9323/metrics
      
  - name: remote
    url : http://100.64.0.1:9323/metrics
```

Without configuration, module attempts to connect to `http://127.0.0.1:9323/metrics`

---
