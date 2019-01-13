# consul

This module will monitor consul health checks.

It produces the following charts:

1. **Service Checks** in status
 * check id

2. **Unbound Checks** in status
 * check id


### configuration

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:8500
      
  - name: remote
    url : http://100.64.0.1:8500
```

Without configuration, module attempts to connect to `http://127.0.0.1:8500`

---