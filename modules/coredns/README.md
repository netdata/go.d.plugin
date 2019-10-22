# CoreDNS

This module will monitor one or more CoreDNS instances.

It produces the following charts:

Summary:

1.  **Number Of DNS Requests** in requests/s
2.  **Number Of DNS Responses** in responses/s
3.  **Number Of Processed And Dropped DNS Requests** in requests/s
4.  **Number Of Dropped DNS Requests Because Of No Matching Zone** in requests/s
5.  **Number Of Panics** in panics/s
6.  **Number Of DNS Requests Per Transport Protocol** in requests/s
7.  **Number Of DNS Requests Per IP Family** in requests/s
8.  **Number Of DNS Requests Per Type** in requests/s
9.  **Number Of DNS Responses Per Rcode** in responses/s

Per Server (if configured):

1.  **Number Of DNS Requests** in requests/s
2.  **Number Of DNS Responses** in responses/s
3.  **Number Of Processed And Dropped DNS Requests** in requests/s
4.  **Number Of DNS Requests Per Transport Protocol** in requests/s
5.  **Number Of DNS Requests Per IP Family** in requests/s
6.  **Number Of DNS Requests Per Type** in requests/s
7.  **Number Of DNS Responses Per Rcode** in responses/s

Per Zone (if configured):

1.  **Number Of DNS Requests** in requests/s
2.  **Number Of DNS Responses** in responses/s
3.  **Number Of DNS Requests Per Transport Protocol** in requests/s
4.  **Number Of DNS Requests Per IP Family** in requests/s
5.  **Number Of DNS Requests Per Type** in requests/s
6.  **Number Of DNS Responses Per Rcode** in responses/s

## Configuration

The module eeeds only the `url` to a CoreDNS `metrics-address`.

Here is an example:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:9153/metrics
      
  - name: remote
    url : http://100.64.0.1:9153/metrics
```

For all available options, please see the module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/coredns.conf).

Without configuration, the module attempts to connect to `http://127.0.0.1:9153/metrics`.
