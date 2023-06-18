# Systemd units state collector

## Overview

[Systemd](https://www.freedesktop.org/wiki/Software/systemd/) is a suite of basic building blocks for a Linux system.

This collector monitors Systemd units state. Works only on Linux systems.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### unit

These metrics refer to the systemd unit.

Labels:

| Label     | Description       |
|-----------|-------------------|
| unit_name | systemd unit name |

Metrics:

| Metric                       |                     Dimensions                     | Unit  |
|------------------------------|:--------------------------------------------------:|:-----:|
| systemd.service_unit_state   | active, inactive, activating, deactivating, failed | state |
| systemd.socket_unit_state    | active, inactive, activating, deactivating, failed | state |
| systemd.target_unit_state    | active, inactive, activating, deactivating, failed | state |
| systemd.path_unit_state      | active, inactive, activating, deactivating, failed | state |
| systemd.device_unit_state    | active, inactive, activating, deactivating, failed | state |
| systemd.mount_unit_state     | active, inactive, activating, deactivating, failed | state |
| systemd.automount_unit_state | active, inactive, activating, deactivating, failed | state |
| systemd.swap_unit_state      | active, inactive, activating, deactivating, failed | state |
| systemd.timer_unit_state     | active, inactive, activating, deactivating, failed | state |
| systemd.scope_unit_state     | active, inactive, activating, deactivating, failed | state |
| systemd.slice_unit_state     | active, inactive, activating, deactivating, failed | state |

## Setup

### Prerequisites

#### Enable in go.d.conf.

This collector is disabled by default. You need to explicitly enable it in the `go.d.conf` file.

### Configuration

#### File

The configuration file name is `go.d/systemdunits.conf`.

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
sudo ./edit-config go.d/systemdunits.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                                     |  Default  | Required |
|:-------------------:|-----------------------------------------------------------------------------------------------------------------|:---------:|:--------:|
|    update_every     | Data collection frequency.                                                                                      |     1     |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                              |     0     |          |
|       include       | Systemd units filter. Pattern syntax is [shell file name pattern](https://golang.org/pkg/path/filepath/#Match). | *.service |          |
|       timeout       | System bus requests timeout.                                                                                    |     1     |          |

</details>

#### Examples

##### Service units

Collect state of all service type units.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: service
    include:
      - '*.service'
```

</details>

##### One specific unit

Collect state of one specific unit.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my-specific-service
    include:
      - 'my-specific.service'
```

</details>

##### All unit types

Collect state of all units.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my-specific-service-unit
    include:
      - '*'
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collect state of all service and socket type units.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: service
    include:
      - '*.service'

  - name: socket
    include:
      - '*.socket'
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `systemdunits` collector, run the `go.d.plugin` with the debug option enabled. The
output
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
  ./go.d.plugin -d -m systemdunits
  ```
