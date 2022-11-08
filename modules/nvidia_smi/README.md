<!--
title: "Nvidia GPU monitoring with Netdata"
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/nvidia_smi/README.md"
description: "Monitors performance metrics using the nvidia-smi CLI tool."
sidebar_label: "nvidia_smi-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Devices"
-->

# Nvidia GPU monitoring with Netdata

Monitors performance metrics (memory usage, fan speed, pcie bandwidth utilization, temperature, etc.)
using the [nvidia-smi](https://developer.nvidia.com/nvidia-system-management-interface) CLI tool.

> **Warning**: under development, collects fewer metrics then python version.

## Metrics

All metrics have "nvidia_smi." prefix.

Labels per scope:

- gpu: product_name, product_brand.

| Metric                        | Scope |        Dimensions        |  Units  |
|-------------------------------|:-----:|:------------------------:|:-------:|
| gpu_pcie_bandwidth_usage      |  gpu  |          rx, tx          |   B/s   |
| gpu_fan_speed_perc            |  gpu  |        fan_speed         |    %    |
| gpu_utilization               |  gpu  |           gpu            |    %    |
| gpu_memory_utilization        |  gpu  |          memory          |    %    |
| gpu_decoder_utilization       |  gpu  |         decoder          |    %    |
| gpu_encoder_utilization       |  gpu  |         encoder          |    %    |
| gpu_frame_buffer_memory_usage |  gpu  |   free, used, reserved   |    B    |
| gpu_bar1_memory_usage         |  gpu  |        free, used        |    B    |
| gpu_temperature               |  gpu  |       temperature        | Celsius |
| gpu_clock_freq                |  gpu  | graphics, video, sm, mem |   MHz   |
| gpu_power_draw                |  gpu  |        power_draw        |  Watts  |
| gpu_performance_state         |  gpu  |          P0-P15          |  state  |

## Configuration

No configuration required.

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
