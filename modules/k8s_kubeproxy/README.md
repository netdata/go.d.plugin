<!--
title: "Kube-proxy monitoring with Netdata"
description: "Monitor the health and performance of Kube-proxy instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_kubeproxy/README.md
sidebar_label: "Kube-proxy"
-->

# Kube-proxy monitoring with Netdata

[`Kube-proxy`](https://kubernetes.io/docs/concepts/overview/components/#kube-proxy) is a network proxy that runs on each
node in your cluster, implementing part of the Kubernetes Service.

This module will monitor one or more `kube-proxy` instances, depending on your configuration.

## Charts

It produces the following charts:

- Sync Proxy Rules in `events/s`
- Sync Proxy Rules Latency in `observes/s`
- Sync Proxy Rules Latency Percentage in `%`
- REST Client HTTP Requests By Status Code in `requests/s`
- REST Client HTTP Requests By Method in `requests/s`
- HTTP Requests Duration in `microseconds`

## Configuration

Edit the `go.d/k8s_kubeproxy.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/k8s_kubeproxy.conf
```

Needs only `url` to `kube-proxy` metric-address. Here is an example for several instances:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:10249/metrics

  - name: remote
    url: http://203.0.113.1:10249/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/k8s_kubeproxy.conf).

## Troubleshooting

To troubleshoot issues with the `k8s_kubeproxy` collector, run the `go.d.plugin` with the debug option enabled. The
output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m k8s_kubeproxy
```
