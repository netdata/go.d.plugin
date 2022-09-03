<!--
title: "CoreDNS monitoring with Netdata"
description: "Monitor the health and performance of CoreDNS servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/coredns/README.md
sidebar_label: "CoreDNS"
-->

# CoreDNS monitoring with Netdata

[`CoreDNS`](https://coredns.io/) is a fast and flexible DNS server.

This module monitor one or more `CoreDNS` instances depending on configuration.

## Charts

It produces the following summary charts:

- Number Of DNS Requests in `requests/s`
- Number Of DNS Responses in `responses/s`
- Number Of Processed And Dropped DNS Requests in `requests/s`
- Number Of Dropped DNS Requests Because Of No Matching Zone in `requests/s`
- Number Of Panics in `panics/s`
- Number Of DNS Requests Per Transport Protocol in `requests/s`
- Number Of DNS Requests Per IP Family in `requests/s`
- Number Of DNS Requests Per Type in `requests/s`
- Number Of DNS Responses Per Rcode in `responses/s`

Per server charts (if configured):

- Number Of DNS Requests in `requests/s`
- Number Of DNS Responses in `responses/s`
- Number Of Processed And Dropped DNS Requests in `requests/s`
- Number Of DNS Requests Per Transport Protocol in `requests/s`
- Number Of DNS Requests Per IP Family in `requests/s`
- Number Of DNS Requests Per Type in `requests/s`
- Number Of DNS Responses Per Rcode in `responses/s`

Per zone charts (if configured):

- Number Of DNS Requests in `requests/s`
- Number Of DNS Responses in `responses/s`
- Number Of DNS Requests Per Transport Protocol in `requests/s`
- Number Of DNS Requests Per IP Family in `requests/s`
- Number Of DNS Requests Per Type in `requests/s`
- Number Of DNS Responses Per Rcode in `responses/s`

## Configuration

Edit the `go.d/coredns.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/coredns.conf
```

The module needs only the `url` to a CoreDNS `metrics-address`. Here is an example for several instances:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9153/metrics

  - name: remote
    url: http://203.0.113.10:9153/metrics
```

For all available options, please see the
module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/coredns.conf).

## Troubleshooting

To troubleshoot issues with the `coredns` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m coredns
  ```
