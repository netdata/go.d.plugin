<!--
title: "Kubelet monitoring with Netdata"
description: "Monitor the health and performance of Kubelet agents with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_kubelet/README.md
sidebar_label: "Kubelet"
-->

# Kubelet monitoring with Netdata

[`Kubelet`](https://kubernetes.io/docs/concepts/overview/components/#kubelet) is an agent that runs on each node in the
cluster. It makes sure that containers are running in a pod.

This module will monitor one or more `kubelet` instances, depending on your configuration.

## Charts

It produces the following charts:

- API Server Audit Requests in `requests/s`
- API Server Failed Data Encryption Key(DEK) Generation Operations in `events/s`
- API Server Latencies Of Data Encryption Key(DEK) Generation Operations in `observes/s`
- API Server Latencies Of Data Encryption Key(DEK) Generation Operations Percentage in `%`
- API Server Storage Envelope Transformation Cache Misses` in `events/s`
- Number Of Containers Currently Running in `containers`
- Number Of Pods Currently Running in `pods`
- Bytes Used By The Pod Logs On The Filesystem in `bytes`
- Runtime Operations By Type in `operations/s`
- Docker Operations By Type in `operations/s`
- Docker Operations Errors By Type in `operations/s`
- Node Configuration-Related Error in `bool`
- PLEG Relisting Interval Summary in `microseconds`
- PLEG Relisting Latency Summary in `microseconds`
- Token Requests To The Alternate Token Source in `requests/s`
- REST Client HTTP Requests By Status Code in `requests/s`
- REST Client HTTP Requests By Method in `requests/s`

Per every plugin:

- Volume Manager State Of The World in `state`

## Configuration

Edit the `go.d/k8s_kubelet.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/k8s_kubelet.conf
```

Needs only `url` to `kubelet` metric-address. Here is an example for 2 instances:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:10255/metrics

  - name: remote
    url: http://203.0.113.10:10255/metrics
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/k8s_kubelet.conf).

## Troubleshooting

To troubleshoot issues with the `k8s_kubelet` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m k8s_kubelet
```
