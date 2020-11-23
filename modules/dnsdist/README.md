<!--
title: "DNSdist monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/dnsdist/README.md
sidebar_label: "DNSdist"
-->

# DNSdist monitoring with Netdata

[`DNS dist`](https://dnsdist.org/) is a highly DNS-, DoS- and abuse-aware loadbalancer. 

This module monitors load-balancer performance and health metrics.

It collects metrics from [the internal webserver](https://dnsdist.org/guides/webserver.html).

Used endpoints:
- [/jsonstat?command=stats](https://dnsdist.org/guides/webserver.html).

## Requirements

For collecting metrics via HTTP, we need [enabled webserver](https://dnsdist.org/guides/webserver.html).

## Charts

1.  **Response latency**

    -   latency-slow
    -   latency100-1000
    -   latency50-100
    -   latency10-50
    -   latency1-10
    -   latency0-1

2.  **Cache performance**

    -   cache-hits
    -   cache-misses

3.  **ACL events**

    -   acl-drops
    -   rule-drop
    -   rule-nxdomain
    -   rule-refused

4.  **Noncompliant data**

    -   empty-queries
    -   no-policy
    -   noncompliant-queries
    -   noncompliant-responses

5.  **Queries**

    -   queries
    -   rdqueries
    -   rdqueries

6.  **Health**

    -   downstream-send-errors
    -   downstream-timeouts
    -   servfail-responses
    -   trunc-failures

## Configuration

Edit the `go.d/dnsdist.conf` configuration file using `edit-config` from the Agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dnsdist.conf
```

Needs `URL` (Complete address to access `DNSdist` stats.), the pair `user` and `password` to authenticate, and if
your `DNSdist` has API key, you will also need to write this parameter.

Here is a configuration example:

```yaml
jobs:
 - name: local
   url: 'http://127.0.0.1:8083'
   username: 'netdata'
   password: 'netdata'
   api_key: 'key'
```

For all available options, see the DNS dist collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsdist.conf).

## Troubleshooting

To troubleshoot issues with the ISC dhcpd collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m dnsdist
```
