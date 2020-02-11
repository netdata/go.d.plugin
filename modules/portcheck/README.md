# Any TCP endpoint monitoring with Netdata

This module will monitors one or more TCP services availability and response time.

## Charts

It produces the following charts for every monitoring port:

-   TCP Check Status in `boolean`
-   Current State Duration in `seconds`
-   TCP Connection Latency in `ms`

## Configuration

Edit the `go.d/portcheck.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/portcheck.conf
```
 
Here is an example for 2 servers:

```yaml
jobs:
  - name: server1
    host: 127.0.0.1
    ports: 
      - 22
      - 23
      
  - name: server2
    host: 203.0.113.10
    ports:
      - 80
      - 81
      - 8081
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/portcheck.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m portcheck
