# DNS queries monitoring with Netdata

This module provides DNS query RTT in milliseconds.

## Charts

It produces only one chart:

-   Query Time in `milliseconds`

## Configuration

Here is an example:

```yaml
jobs:
  - name: job1
    domains:
      - google.com
      - github.com
      - reddit.com
    servers:
      - 8.8.8.8
      - 8.8.4.4
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dns_query.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m dns_query
