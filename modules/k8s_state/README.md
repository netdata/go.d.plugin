<!--
title: "Kubernetes cluster state monitoring with Netdata"
description: "Monitor the state of your Kubernetes clusters with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_state/README.md
sidebar_label: "Kubernetes cluster state"
-->

# Kubernetes cluster state monitoring with Netdata

[Kubernetes](https://kubernetes.io/) is an open-source container orchestration system for automating software
deployment, scaling, and management.

This module collects health metrics for the following Kubernetes resources:

- [Nodes](https://kubernetes.io/docs/concepts/architecture/nodes/).
- [Pods](https://kubernetes.io/docs/concepts/workloads/pods/).

## Metrics

| Label               |                Description                |
|---------------------|:-----------------------------------------:|
| k8s_kind            | The REST resource this object represents. |
| k8s_cluster_id      |                     %                     |
| k8s_cluster_name    |                     %                     |
| k8s_node_name       |                     %                     |
| k8s_namespace       |                     %                     |
| k8s_controller_kind |                     %                     |
| k8s_controller_name |                     %                     |
| k8s_pod_uid         |                     %                     |
| k8s_pod_name        |                     %                     |
| k8s_qos_class       |                     %                     |
| k8s_container_id    |                     %                     |
| k8s_container_name  |                     %                     |

### Nodes

| Metric                            |   Units    | Source     |
|-----------------------------------|:----------:|------------|
| node_allocatable_cpu_utilization  |     %      | calculated |
| node_allocatable_cpu_used         |  millicpu  | calculated |
| node_allocatable_mem_utilization  |     %      | calculated |
| node_allocatable_mem_used         |   bytes    | calculated |
| node_allocatable_pods_utilization |     %      | calculated |
| node_allocatable_pods_usage       |    pods    | calculated |
| node_condition                    |   status   | desc       |
| node_pods_readiness               |     %      | desc       |
| node_pods_readiness_state         |    pods    | desc       |
| node_pods_condition               |    pods    | desc       |
| node_pods_phase                   |    pods    | desc       |
| node_containers                   | containers | desc       |
| node_containers_state             | containers | desc       |
| node_init_containers_state        | containers | desc       |
| node_age                          |  seconds   | desc       |

### Pods

- Pod

| Metric                    |   Units    | Description |
|---------------------------|:----------:|-------------|
| pod_allocated_cpu         |     %      | desc        |
| pod_allocated_cpu_used    |  millicpu  | desc        |
| pod_allocated_mem         |     %      | desc        |
| pod_allocated_mem_used    |   bytes    | desc        |
| pod_readiness_state       |   state    | desc        |
| pod_condition             |   state    | desc        |
| pod_phase                 |   state    | desc        |
| pod_age                   |  seconds   | desc        |
| pod_containers            |  seconds   | desc        |
| pod_containers_state      | containers | desc        |
| pod_init_containers_state | containers | desc        |

- Pod Containers

| Metric                                |  Units   | Description |
|---------------------------------------|:--------:|-------------|
| pod_container_readiness_state         |  state   | desc        |
| pod_container_restarts                | restarts | desc        |
| pod_container_state                   |  state   | desc        |
| pod_container_waiting_state_reason    |  state   | desc        |
| pod_container_terminated_state_reason |  state   | desc        |

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

To troubleshoot issues with the `k8s_state` collector, run the `go.d.plugin` with the debug option enabled. The
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
./go.d.plugin -d -m k8s_state
```
