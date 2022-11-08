<!--
title: "Kubelet monitoring with Netdata"
description: "Monitor the health and performance of Kubelet agents with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_kubelet/README.md"
sidebar_label: "Kubelet"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Container orchestrators/Kubernetes"
-->

# Kubelet monitoring with Netdata

[`Kubelet`](https://kubernetes.io/docs/concepts/overview/components/#kubelet) is an agent that runs on each node in the
cluster. It makes sure that containers are running in a pod.

This module will monitor one or more `kubelet` instances, depending on your configuration.

## Metrics

All metrics have "k8s_kubelet." prefix.

| Metric                                                  |     Scope      |                                                       Dimensions                                                        |       Units        |
|---------------------------------------------------------|:--------------:|:-----------------------------------------------------------------------------------------------------------------------:|:------------------:|
| apiserver_audit_requests_rejected                       |     global     |                                                        rejected                                                         |     requests/s     |
| apiserver_storage_data_key_generation_failures          |     global     |                                                        failures                                                         |      events/s      |
| apiserver_storage_data_key_generation_latencies         |     global     | 5_µs, 10_µs, 20_µs, 40_µs, 80_µs, 160_µs, 320_µs, 640_µs, 1280_µs, 2560_µs, 5120_µs, 10240_µs, 20480_µs, 40960_µs, +Inf |     observes/s     |
| apiserver_storage_data_key_generation_latencies_percent |     global     | 5_µs, 10_µs, 20_µs, 40_µs, 80_µs, 160_µs, 320_µs, 640_µs, 1280_µs, 2560_µs, 5120_µs, 10240_µs, 20480_µs, 40960_µs, +Inf |     percentage     |
| apiserver_storage_envelope_transformation_cache_misses  |     global     |                                                      cache misses                                                       |      events/s      |
| kubelet_containers_running                              |     global     |                                                          total                                                          | running_containers |
| kubelet_pods_running                                    |     global     |                                                          total                                                          |    running_pods    |
| kubelet_pods_log_filesystem_used_bytes                  |     global     |                                        <i>a dimension per namespace and pod</i>                                         |         B          |
| kubelet_runtime_operations                              |     global     |                                          <i>a dimension per operation type</i>                                          |    operations/s    |
| kubelet_runtime_operations_errors                       |     global     |                                          <i>a dimension per operation type</i>                                          |      errors/s      |
| kubelet_docker_operations                               |     global     |                                          <i>a dimension per operation type</i>                                          |    operations/s    |
| kubelet_docker_operations_errors                        |     global     |                                          <i>a dimension per operation type</i>                                          |      errors/s      |
| kubelet_node_config_error                               |     global     |                                                   experiencing_error                                                    |        bool        |
| kubelet_pleg_relist_interval_microseconds               |     global     |                                                     0.5, 0.9, 0.99                                                      |    microseconds    |
| kubelet_pleg_relist_latency_microseconds                |     global     |                                                     0.5, 0.9, 0.99                                                      |    microseconds    |
| kubelet_token_requests                                  |     global     |                                                      total, failed                                                      |  token_requests/s  |
| rest_client_requests_by_code                            |     global     |                                         <i>a dimension per HTTP status code</i>                                         |     requests/s     |
| rest_client_requests_by_method                          |     global     |                                           <i>a dimension per HTTP method</i>                                            |     requests/s     |
| volume_manager_total_volumes                            | volume manager |                                                     actual, desired                                                     |       state        |

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

