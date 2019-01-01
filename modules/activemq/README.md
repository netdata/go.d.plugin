# activemq

This plugin collects queues and topics metrics using ActiveMQ Console API.

It produces following charts per queue and per topic:

1. **Messages** in messages/s
 * enqueued
 * dequeued
 * unprocessed
 
2. **Consumers** in consumers
 * consumers
 

### configuration

Here is an example for 2 servers:

```yaml
jobs:
  - name: job1
    url: http://127.0.0.1:8161
    webadmin: admin
    
  - name: remote
    url: http://100.127.0.1:8161
    webadmin: admin
```

Without configuration, module won't work.

---