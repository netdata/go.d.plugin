# NVMe devices collector

## Overview

NVMe (nonvolatile memory express) is a new storage access and transport protocol for flash and next-generation
solid-state drives (SSDs) that delivers the highest throughput and fastest response times yet for all types of
enterprise workloads.

This collector monitors the health of NVMe devices using the command line
tool [nvme](https://github.com/linux-nvme/nvme-cli#nvme-cli), which can only be run by the root user. It uses `sudo` and
assumes it is set up so that the netdata user can execute `nvme` as root without a password.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### device

These metrics refer to the NVME device.

Labels:

| Label  | Description      |
|--------|------------------|
| device | NVMe device name |

Metrics:

| Metric                                          |                                                           Dimensions                                                           |     Unit      |
|-------------------------------------------------|:------------------------------------------------------------------------------------------------------------------------------:|:-------------:|
| nvme.device_estimated_endurance_perc            |                                                              used                                                              |       %       |
| nvme.device_available_spare_perc                |                                                             spare                                                              |       %       |
| nvme.device_composite_temperature               |                                                          temperature                                                           |    celsius    |
| nvme.device_io_transferred_count                |                                                         read, written                                                          |     bytes     |
| nvme.device_power_cycles_count                  |                                                             power                                                              |    cycles     |
| nvme.device_power_on_time                       |                                                            power-on                                                            |    seconds    |
| nvme.device_critical_warnings_state             | available_spare, temp_threshold, nvm_subsystem_reliability, read_only, volatile_mem_backup_failed, persistent_memory_read_only |     state     |
| nvme.device_unsafe_shutdowns_count              |                                                             unsafe                                                             |   shutdowns   |
| nvme.device_media_errors_rate                   |                                                             media                                                              |   errors/s    |
| nvme.device_error_log_entries_rate              |                                                           error_log                                                            |   entries/s   |
| nvme.device_warning_composite_temperature_time  |                                                             wctemp                                                             |    seconds    |
| nvme.device_critical_composite_temperature_time |                                                             cctemp                                                             |    seconds    |
| nvme.device_thermal_mgmt_temp1_transitions_rate |                                                             temp1                                                              | transitions/s |
| nvme.device_thermal_mgmt_temp2_transitions_rate |                                                             temp2                                                              | transitions/s |
| nvme.device_thermal_mgmt_temp1_time             |                                                             temp1                                                              |    seconds    |
| nvme.device_thermal_mgmt_temp2_time             |                                                             temp2                                                              |    seconds    |

## Setup

### Prerequisites

#### Install nvme-cli

See [Distro Support](https://github.com/linux-nvme/nvme-cli#distro-support). Install `nvme-cli` using your
distribution's package manager.

#### Allow netdata to execute nvme

Add the netdata user to `/etc/sudoers` (use `which nvme` to find the full path to the binary):

```bash
netdata ALL=(root) NOPASSWD: /usr/sbin/nvme
```

### Configuration

#### File

The configuration file name is `go.d/nvme.conf`.

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
sudo ./edit-config go.d/nvme.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                                                                | Default | Required |
|:-------------------:|--------------------------------------------------------------------------------------------------------------------------------------------|:-------:|:--------:|
|    update_every     | Data collection frequency.                                                                                                                 |   10    |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                                                         |    0    |          |
|     binary_path     | Path to nvme binary. The default is "nvme" and the executable is looked for in the directories specified in the PATH environment variable. |  nvme   |          |
|       timeout       | nvme binary execution timeout.                                                                                                             |    2    |          |

</details>

#### Examples

##### Custom binary path

The executable is not in the directories specified in the PATH environment variable.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: nvme
    binary_path: /usr/local/sbin/nvme
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `nvme` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m nvme
  ```
