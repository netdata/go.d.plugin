# DNS queries monitoring with Netdata

This module provides DNS query RTT in milliseconds.

## Charts

It produces only one chart:

-   Query Time in `milliseconds`

## Configuration

Edit the `go.d/dns_query.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dns_query.conf
```

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
