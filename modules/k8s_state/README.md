<!--
title: "Kubernetes cluster state monitoring with Netdata"
description: "Monitor the state of your Kubernetes clusters with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/k8s_state/README.md"
sidebar_label: "Kubernetes cluster state"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Container orchestrators/Kubernetes"
-->

# Kubernetes cluster state monitoring with Netdata

[Kubernetes](https://kubernetes.io/) is an open-source container orchestration system for automating software
deployment, scaling, and management.

This module collects health metrics for the following Kubernetes resources:

- [Nodes](https://kubernetes.io/docs/concepts/architecture/nodes/).
- [Pods](https://kubernetes.io/docs/concepts/workloads/pods/).

## Requirements

- Only works when Netdata is running inside a Kubernetes cluster.
- RBAC: needs **list**, **watch** verbs for **pod** and **node** resources.
- RBAC: needs **get** verb for **namespace** resource.

## Metrics

All metrics have "k8s_state." prefix.

### Node

| Metric                                    |                           Dimensions                            |   Units    |
|-------------------------------------------|:---------------------------------------------------------------:|:----------:|
| node_allocatable_cpu_requests_utilization |                            requests                             |     %      |
| node_allocatable_cpu_requests_used        |                            requests                             |  millicpu  |
| node_allocatable_cpu_limits_utilization   |                             limits                              |     %      |
| node_allocatable_cpu_limits_used          |                             limits                              |  millicpu  |
| node_allocatable_mem_requests_utilization |                            requests                             |     %      |
| node_allocatable_mem_requests_used        |                            requests                             |   bytes    |
| node_allocatable_mem_limits_utilization   |                             limits                              |     %      |
| node_allocatable_mem_limits_used          |                             limits                              |   bytes    |
| node_allocatable_pods_utilization         |                            allocated                            |     %      |
| node_allocatable_pods_usage               |                      available, allocated                       |    pods    |
| node_condition                            |                    <i>added dynamically</i>                     |   status   |
| node_schedulability                       |                   schedulable, unschedulable                    |   state    |
| node_pods_readiness                       |                              ready                              |     %      |
| node_pods_readiness_state                 |                         ready, unready                          |    pods    |
| node_pods_condition                       | pod_ready, pod_scheduled,<br/>pod_initialized, containers_ready |    pods    |
| node_pods_phase                           |               running, failed, succeeded, pending               |    pods    |
| node_containers                           |                   containers, init_containers                   | containers |
| node_containers_state                     |                  running, waiting, terminated                   | containers |
| node_init_containers_state                |                  running, waiting, terminated                   | containers |
| node_age                                  |                               age                               |  seconds   |

### Pod

| Metric                                |                           Dimensions                            |   Units    |
|---------------------------------------|:---------------------------------------------------------------:|:----------:|
| pod_cpu_requests_used                 |                            requests                             |  millicpu  |
| pod_cpu_limits_used                   |                             limits                              |  millicpu  |
| pod_mem_requests_used                 |                            requests                             |   bytes    |
| pod_mem_limits_used                   |                             limits                              |   bytes    |
| pod_condition                         | pod_ready, pod_scheduled,<br/>pod_initialized, containers_ready |   state    |
| pod_phase                             |               running, failed, succeeded, pending               |   state    |
| pod_age                               |                               age                               |  seconds   |
| pod_containers                        |                   containers, init_containers                   | containers |
| pod_containers_state                  |                  running, waiting, terminated                   | containers |
| pod_init_containers_state             |                  running, waiting, terminated                   | containers |

### Pod container

| Metric                                |                           Dimensions                           |   Units    |
|---------------------------------------|:--------------------------------------------------------------:|:----------:|
| pod_container_readiness_state         |                             ready                              |   state    |
| pod_container_restarts                |                            restarts                            | restarts/s |
| pod_container_state                   |                  running, waiting, terminated                  |   state    |
| pod_container_waiting_state_reason    |                    <i>added dynamically</i>                    |   state    |
| pod_container_terminated_state_reason |                    <i>added dynamically</i>                    |   state    |

## Labels

- 'k8s_cluster_id' value is 'kube-system' namespace UID.
- 'k8s_cluster_name' currently only appears when running on [GKE](https://cloud.google.com/kubernetes-engine).

| Label               | Node | Pod | Container |
|---------------------|:----:|:---:|:---------:|
| k8s_kind            | yes  | yes |    yes    |
| k8s_cluster_id      | yes  | yes |    yes    |
| k8s_cluster_name    | yes  | yes |    yes    |
| k8s_node_name       | yes  | yes |    yes    |
| k8s_namespace       |      | yes |    yes    |
| k8s_controller_kind |      | yes |    yes    |
| k8s_controller_name |      | yes |    yes    |
| k8s_pod_uid         |      | yes |    yes    |
| k8s_pod_name        |      | yes |    yes    |
| k8s_qos_class       |      | yes |    yes    |
| k8s_container_id    |      |     |    yes    |
| k8s_container_name  |      |     |    yes    |

## Configuration

No configuration is needed. This module is enabled when you install Netdata
using [netdata/helmchart](https://github.com/netdata/helmchart#netdata-helm-chart-for-kubernetes-deployments).

## Troubleshooting

To troubleshoot issues with the `k8s_state` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m k8s_state
  ```
