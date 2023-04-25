<!--
title: "Kubelet monitoring with Netdata"
description: "Monitor the health and performance of Kubelet agents with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_kubelet/README.md"
sidebar_label: "Kubelet"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Container orchestrators/Kubernetes"
-->

# Kubelet collector

[`Kubelet`](https://kubernetes.io/docs/concepts/overview/components/#kubelet) is an agent that runs on each node in the
cluster. It makes sure that containers are running in a pod.

This module will monitor one or more `kubelet` instances, depending on your configuration.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/k8s_kubelet/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/k8s_kubelet.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

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
  ./go.d.plugin -d -m k8s_kubelet
  ```

