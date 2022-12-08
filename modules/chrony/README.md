<!--
title: "Chrony monitoring with Netdata"
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/chrony/README.md"
sidebar_label: "Chrony"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Apps"
-->

# Chrony monitoring with Netdata

[chrony](https://chrony.tuxfamily.org/) is a versatile implementation of the Network Time Protocol (NTP).

This module monitors the system's clock performance and peers activity status using Chrony communication protocol v6.

## Metrics

All metrics have "chrony." prefix.

| Metric               | Scope  |                        Dimensions                        |   Units   |
|----------------------|:------:|:--------------------------------------------------------:|:---------:|
| stratum              | global |                         stratum                          |   level   |
| current_correction   | global |                    current_correction                    |  seconds  |
| root_delay           | global |                        root_delay                        |  seconds  |
| root_dispersion      | global |                     root_dispersion                      |  seconds  |
| last_offset          | global |                          offset                          |  seconds  |
| rms_offset           | global |                          offset                          |  seconds  |
| frequency            | global |                        frequency                         |    ppm    |
| residual_frequency   | global |                    residual_frequency                    |    ppm    |
| skew                 | global |                           skew                           | frequency |
| update_interval      | global |                     update_interval                      |  seconds  |
| ref_measurement_time | global |                   ref_measurement_time                   |  seconds  |
| leap_status          | global |   normal, insert_second, delete_second, unsynchronised   |  status   |
| activity             | global | online, offline, burst_online, burst_offline, unresolved |  sources  |

## Configuration

Edit the `go.d/chrony.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata   # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/chrony.conf
```

Configuration example:

```yaml
jobs:
  - name: local
    address: '127.0.0.1:323'
    timeout: 1

  - name: remote
    address: '203.0.113.0:323'
    timeout: 3
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/chrony.conf).

---

## Troubleshooting

To troubleshoot issues with the `chrony` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m chrony
  ```
