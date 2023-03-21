<!--
title: "Kube-proxy monitoring with Netdata"
description: "Monitor the health and performance of Kube-proxy instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_kubeproxy/README.md"
sidebar_label: "Kube-proxy"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Container orchestrators/Kubernetes"
-->

# Kube-proxy collector

[`Kube-proxy`](https://kubernetes.io/docs/concepts/overview/components/#kube-proxy) is a network proxy that runs on each
node in your cluster, implementing part of the Kubernetes Service.

This module will monitor one or more `kube-proxy` instances, depending on your configuration.

## Metrics

All metrics have "k8s_kubeproxy." prefix.

| Metric                                         | Scope  |                                                   Dimensions                                                   |    Units     |
|------------------------------------------------|:------:|:--------------------------------------------------------------------------------------------------------------:|:------------:|
| kubeproxy_sync_proxy_rules                     | global |                                                sync_proxy_rules                                                |   events/s   |
| kubeproxy_sync_proxy_rules_latency_microsecond | global | 0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048, 4.096, 8.192, 16.384, +Inf |  observes/s  |
| kubeproxy_sync_proxy_rules_latency             | global | 0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048, 4.096, 8.192, 16.384, +Inf |  percentage  |
| rest_client_requests_by_code                   | global |                                    <i>a dimension per HTTP status code</i>                                     |  requests/s  |
| rest_client_requests_by_method                 | global |                                       <i>a dimension per HTTP method</i>                                       |  requests/s  |
| http_request_duration                          | global |                                                 0.5, 0.9, 0.99                                                 | microseconds |

## Configuration

Edit the `go.d/k8s_kubeproxy.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

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
  ./go.d.plugin -d -m k8s_kubeproxy
  ```
