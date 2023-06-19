# Nvidia GPU collector

## Overview

This collector monitors GPUs performance metrics using
the [nvidia-smi](https://developer.nvidia.com/nvidia-system-management-interface) CLI tool.

> **Warning**: under development, [loop mode](https://github.com/netdata/netdata/issues/14522) not implemented yet.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### gpu

These metrics refer to the GPU.

Labels:

| Label        | Description                                   |
|--------------|-----------------------------------------------|
| uuid         | GPU id (e.g. 00000000:00:04.0)                |
| product_name | GPU product name (e.g. NVIDIA A100-SXM4-40GB) |

Metrics:

| Metric                                    |        Dimensions        |  Unit   | XML/CSV |
|-------------------------------------------|:------------------------:|:-------:|:-------:|
| nvidia_smi.gpu_pcie_bandwidth_usage       |          rx, tx          |   B/s   |   + -   |
| nvidia_smi.gpu_pcie_bandwidth_utilization |          rx, tx          |    %    |   + -   |
| nvidia_smi.gpu_fan_speed_perc             |        fan_speed         |    %    |   + +   |
| nvidia_smi.gpu_utilization                |           gpu            |    %    |   + +   |
| nvidia_smi.gpu_memory_utilization         |          memory          |    %    |   + +   |
| nvidia_smi.gpu_decoder_utilization        |         decoder          |    %    |   + -   |
| nvidia_smi.gpu_encoder_utilization        |         encoder          |    %    |   + -   |
| nvidia_smi.gpu_frame_buffer_memory_usage  |   free, used, reserved   |    B    |   + +   |
| nvidia_smi.gpu_bar1_memory_usage          |        free, used        |    B    |   + -   |
| nvidia_smi.gpu_temperature                |       temperature        | Celsius |   + +   |
| nvidia_smi.gpu_voltage                    |         voltage          |    V    |   + -   |
| nvidia_smi.gpu_clock_freq                 | graphics, video, sm, mem |   MHz   |   + +   |
| nvidia_smi.gpu_power_draw                 |        power_draw        |  Watts  |   + +   |
| nvidia_smi.gpu_performance_state          |          P0-P15          |  state  |   + +   |
| nvidia_smi.gpu_mig_mode_current_status    |    enabled, disabled     | status  |   + -   |
| nvidia_smi.gpu_mig_devices_count          |           mig            | devices |   + -   |

### mig

These metrics refer to the Multi-Instance GPU (MIG).

Labels:

| Label           | Description                                   |
|-----------------|-----------------------------------------------|
| uuid            | GPU id (e.g. 00000000:00:04.0)                |
| product_name    | GPU product name (e.g. NVIDIA A100-SXM4-40GB) |
| gpu_instance_id | GPU instance id (e.g. 1)                      |

Metrics:

| Metric                                       |      Dimensions      | Unit | XML/CSV |
|----------------------------------------------|:--------------------:|:----:|:-------:|
| nvidia_smi.gpu_mig_frame_buffer_memory_usage | free, used, reserved |  B   |   + -   |
| nvidia_smi.gpu_mig_bar1_memory_usage         |      free, used      |  B   |   + -   |

## Setup

### Prerequisites

#### Enable in go.d.conf.

This collector is disabled by default. You need to explicitly enable it in the `go.d.conf` file.

### Configuration

#### File

The configuration file name is `go.d/nvidia_smi.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/nvidia_smi.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                                                                            |  Default   | Required |
|:-------------------:|--------------------------------------------------------------------------------------------------------------------------------------------------------|:----------:|:--------:|
|    update_every     | Data collection frequency.                                                                                                                             |     10     |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                                                                     |     0      |          |
|     binary_path     | Path to nvidia_smi binary. The default is "nvidia_smi" and the executable is looked for in the directories specified in the PATH environment variable. | nvidia_smi |          |
|       timeout       | nvidia_smi binary execution timeout.                                                                                                                   |     2      |          |
|   use_csv_format    | Used format when requesting GPU information. XML is used if set to 'no'.                                                                               |    yes     |          |

</details>

##### use_csv_format

This module supports data collection in CSV and XML formats. The default is CSV.

- XML provides more metrics, but requesting GPU information consumes more CPU, especially if there are multiple GPUs
  in the system.
- CSV provides fewer metrics, but is much lighter than XML in terms of CPU usage.

#### Examples

##### XML format

Use XML format when requesting GPU information.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: nvidia_smi
    use_csv_format: no
```

</details>

##### Custom binary path

The executable is not in the directories specified in the PATH environment variable.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: nvidia_smi
    binary_path: /usr/local/sbin/nvidia_smi
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `nvidia_smi` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m nvidia_smi
  ```
