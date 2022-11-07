<!--
title: "NVMe devices monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/nvme/README.md
description: "Monitors NVMe devices health metrics using the nvme CLI tool."
sidebar_label: "NVMe devices"
-->

# NVMe devices monitoring with Netdata

Monitors health metrics (estimated endurance, space capacity, critical warnings, temperature, etc.) using
the [nvme](https://github.com/linux-nvme/nvme-cli#nvme-cli) CLI tool.

## Requirements

The module uses `nvme`, which can only be executed by root. It uses `sudo` and assumes that it is configured such that
the netdata user can execute `nvme` as root without a password.

- Add the following line to the `/etc/sudoers` file:

  ```bash
  netdata ALL=(root)       NOPASSWD: /usr/sbin/nvme
  ```
  Use `which nvme` to find the full path to the binary.

- Reset Netdata's systemd
  unit [CapabilityBoundingSet](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#Capabilities) (Linux
  distributions with systemd). As the root user, do the following:

  ```bash
  mkdir /etc/systemd/system/netdata.service.d
  echo -e '[Service]\nCapabilityBoundingSet=~' | tee /etc/systemd/system/netdata.service.d/unset-capability-bounding-set.conf
  systemctl daemon-reload
  systemctl restart netdata.service
  ```

  The default CapabilityBoundingSet doesn't allow using sudo, and is quite strict in general. Resetting is not
  optimal, but a next-best solution given the inability to execute nvme using sudo.

## Metrics

All metrics have "nvme." prefix.

Labels per scope:

- device: device.

| Metric                                     | Scope  |                                                           Dimensions                                                           |    Units    |
|--------------------------------------------|:------:|:------------------------------------------------------------------------------------------------------------------------------:|:-----------:|
| device_estimated_endurance_perc            | device |                                                              used                                                              |      %      |
| device_available_spare_perc                | device |                                                             spare                                                              |      %      |
| device_composite_temperature               | device |                                                          temperature                                                           |   celsius   |
| device_io_transferred_count                | device |                                                         read, written                                                          |    bytes    |
| device_power_cycles_count                  | device |                                                             power                                                              |   cycles    |
| device_power_on_time                       | device |                                                            power-on                                                            |   seconds   |
| device_critical_warnings_state             | device | available_spare, temp_threshold, nvm_subsystem_reliability, read_only, volatile_mem_backup_failed, persistent_memory_read_only |    state    |
| device_unsafe_shutdowns_count              | device |                                                             unsafe                                                             |  shutdowns  |
| device_media_errors_rate                   | device |                                                             media                                                              |   errors    |
| device_error_log_entries_rate              | device |                                                           error_log                                                            |   entries   |
| device_warning_composite_temperature_time  | device |                                                             wctemp                                                             |   seconds   |
| device_critical_composite_temperature_time | device |                                                             cctemp                                                             |   seconds   |
| device_thermal_mgmt_temp1_transitions_rate | device |                                                             temp1                                                              | transitions |
| device_thermal_mgmt_temp2_transitions_rate | device |                                                             temp2                                                              | transitions |
| device_thermal_mgmt_temp1_time             | device |                                                             temp1                                                              |   seconds   |
| device_thermal_mgmt_temp2_time             | device |                                                             temp2                                                              |   seconds   |

## Configuration

No configuration required.

## Troubleshooting

To troubleshoot issues with the `nvme` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m nvme
  ```
