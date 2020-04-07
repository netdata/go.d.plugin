# CoreDNS monitoring with Netdata

[`CoreDNS`](https://coredns.io/) is a fast and flexible DNS server.

This module monitor one or more `CoreDNS` instances depending on configuration.

## Charts

It produces the following summary charts:

-   Number Of DNS Requests in `requests/s`
-   Number Of DNS Responses in `responses/s`
-   Number Of Processed And Dropped DNS Requests in `requests/s`
-   Number Of Dropped DNS Requests Because Of No Matching Zone in `requests/s`
-   Number Of Panics in `panics/s`
-   Number Of DNS Requests Per Transport Protocol in `requests/s`
-   Number Of DNS Requests Per IP Family in `requests/s`
-   Number Of DNS Requests Per Type in `requests/s`
-   Number Of DNS Responses Per Rcode in `responses/s`

Per server charts (if configured):

-   Number Of DNS Requests in `requests/s`
-   Number Of DNS Responses in `responses/s`
-   Number Of Processed And Dropped DNS Requests in `requests/s`
-   Number Of DNS Requests Per Transport Protocol in `requests/s`
-   Number Of DNS Requests Per IP Family in `requests/s`
-   Number Of DNS Requests Per Type in `requests/s`
-   Number Of DNS Responses Per Rcode in `responses/s`

Per zone charts (if configured):

-   Number Of DNS Requests in `requests/s`
-   Number Of DNS Responses in `responses/s`
-   Number Of DNS Requests Per Transport Protocol in `requests/s`
-   Number Of DNS Requests Per IP Family in `requests/s`
-   Number Of DNS Requests Per Type in `requests/s`
-   Number Of DNS Responses Per Rcode in `responses/s`

## Configuration

Edit the `go.d/coredns.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/coredns.conf
```

The module needs only the `url` to a CoreDNS `metrics-address`. Here is an example for several instances:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:9153/metrics
      
  - name: remote
    url : http://203.0.113.10:9153/metrics
```

For all available options, please see the module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/coredns.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m coredns
