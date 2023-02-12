<!--
title: "NTP daemon monitoring with Netdata"
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/ntpd/README.md"
sidebar_label: "NTP daemon"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/Apps"
-->

# NTP daemon monitoring with Netdata

Monitors the system variables of the local `ntpd` daemon (optional incl. variables of the polled peers) using the NTP
Control Message Protocol via UDP socket, similar to `ntpq`,
the [standard NTP query program](http://doc.ntp.org/current-stable/ntpq.html).

## Metrics

All metrics have "ntpd." prefix.

Labels per scope:

- global: no labels.
- peer: peer_address.

| Metric          | Scope  |    Dimensions    |    Units     |
|-----------------|:------:|:----------------:|:------------:|
| sys_offset      | global |      offset      | milliseconds |
| sys_jitter      | global |  system, clock   | milliseconds |
| sys_frequency   | global |    frequency     |     ppm      |
| sys_wander      | global |      clock       |     ppm      |
| sys_rootdelay   | global |      delay       | milliseconds |
| sys_rootdisp    | global |    dispersion    | milliseconds |
| sys_stratum     | global |     stratum      |   stratum    |
| sys_tc          | global | current, minimum |     log2     |
| sys_precision   | global |    precision     |     log2     |
| peer_offset     |  peer  |      offset      | milliseconds |
| peer_delay      |  peer  |      delay       | milliseconds |
| peer_dispersion |  peer  |    dispersion    | milliseconds |
| peer_jitter     |  peer  |      jitter      | milliseconds |
| peer_xleave     |  peer  |      xleave      | milliseconds |
| peer_rootdelay  |  peer  |    rootdelay     | milliseconds |
| peer_rootdisp   |  peer  |    dispersion    | milliseconds |
| peer_stratum    |  peer  |     stratum      |   stratum    |
| peer_hmode      |  peer  |      hmode       |    hmode     |
| peer_pmode      |  peer  |      pmode       |    pmode     |
| peer_hpoll      |  peer  |      hpoll       |     log2     |
| peer_ppoll      |  peer  |      ppoll       |     log2     |
| peer_precision  |  peer  |    precision     |     log2     |

## Configuration

Edit the `go.d/ntpd.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory, if different
sudo ./edit-config go.d/ntpd.conf

```

Configuration example:

```yaml
jobs:
  - name: local
    address: '127.0.0.1:123'
    collect_peers: no

  - name: remote
    address: '203.0.113.0:123'
    timeout: 3
    collect_peers: no
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/ntpd.conf).

---

## Troubleshooting

To troubleshoot issues with the `ntpd` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m ntpd
  ```
