<!--
title: "Nvidia GPU monitoring with Netdata"
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/nvidia_smi/README.md"
description: "Monitors performance metrics using the nvidia-smi CLI tool."
sidebar_label: "nvidia_smi-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Devices"
-->

# Nvidia GPU collector

Monitors performance metrics (memory usage, fan speed, pcie bandwidth utilization, temperature, etc.)
using the [nvidia-smi](https://developer.nvidia.com/nvidia-system-management-interface) CLI tool.

> **Warning**: under development, collects fewer metrics then python version.

## Metrics

All metrics have "nvidia_smi." prefix.

Labels per scope:

- gpu: uuid, product_name.
- mig: gpu_uuid, gpu_product_name, gpu_instance_id

| Metric                            | Scope |        Dimensions        |  Units  | XML | CSV |
|-----------------------------------|:-----:|:------------------------:|:-------:|:---:|:---:|
| gpu_pcie_bandwidth_usage          |  gpu  |          rx, tx          |   B/s   | yes | no  |
| gpu_pcie_bandwidth_utilization    |  gpu  |          rx, tx          |    %    | yes | no  |
| gpu_fan_speed_perc                |  gpu  |        fan_speed         |    %    | yes | yes |
| gpu_utilization                   |  gpu  |           gpu            |    %    | yes | yes |
| gpu_memory_utilization            |  gpu  |          memory          |    %    | yes | yes |
| gpu_decoder_utilization           |  gpu  |         decoder          |    %    | yes | no  |
| gpu_encoder_utilization           |  gpu  |         encoder          |    %    | yes | no  |
| gpu_frame_buffer_memory_usage     |  gpu  |   free, used, reserved   |    B    | yes | yes |
| gpu_bar1_memory_usage             |  gpu  |        free, used        |    B    | yes | no  |
| gpu_temperature                   |  gpu  |       temperature        | Celsius | yes | yes |
| gpu_voltage                       |  gpu  |         voltage          |    V    | yes | no  |
| gpu_clock_freq                    |  gpu  | graphics, video, sm, mem |   MHz   | yes | yes |
| gpu_power_draw                    |  gpu  |        power_draw        |  Watts  | yes | yes |
| gpu_performance_state             |  gpu  |          P0-P15          |  state  | yes | yes |
| gpu_mig_mode_current_status       |  gpu  |    enabled, disabled     | status  | yes | no  |
| gpu_mig_devices_count             |  gpu  |           mig            | devices | yes | no  |
| gpu_mig_frame_buffer_memory_usage |  mig  |   free, used, reserved   |    B    | yes | no  |
| gpu_mig_bar1_memory_usage         |  mig  |        free, used        |    B    | yes | no  |

## Configuration

This module supports data collection in CSV and XML formats. The default is CSV.

- XML provides more metrics, but requesting GPU information consumes more CPU, especially if there are multiple GPU
  cards in the system.
- CSV provides fewer metrics, but is much lighter than XML in terms of CPU usage.

The format can be changed in the configuration file.

Edit the `go.d/nvidia_smi.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```yaml
jobs:
  - name: nvidia_smi
    use_csv_format: no # set to 'no' to use the XML format.
```

## Troubleshooting

To troubleshoot issues with the `nvidia_smi` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m nvidia_smi
  ```
