# logstash

This module will monitor one or more logstash servers depending on configuration.

Servers can be either local or remote.

It produces following charts:

1. **JVM Threads** in count
 * threads

2. **JVM Heap Memory Percentage** in percent
 * in use

3. **JVM Heap Memory** in KiB
 * used
 * committed

4. **JVM Pool Survivor Memory** in KiB
 * used
 * committed

5. **JVM Pool Old Memory** in KiB
 * used
 * committed

6. **JVM Pool Young Memory** in KiB
 * used
 * committed

7. **Garbage Collection Count** in counts/s
 * young
 * old

8. **Time Spent On Garbage Collection** in ms
 * young
 * old

9. **Uptime** in time
 * uptime

### configuration
Detailed logstash configuration with all available parameters can be found in ['go.d/logstash.conf'](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/logstash.conf)

Here is a simple example for local and remote server:

```yaml
jobs:
  - name: local
    url : http://localhost:9600

  - name: remote
    url : http://10.0.0.1:9600
```

---
