<!--
title: "systemd-logind monitoring with Netdata"
description: "Monitors number of sessions and users with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/logind/README.md"
sidebar_label: "logind-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/System metrics"
-->

# systemd-logind collector

[systemd-logind](https://www.freedesktop.org/software/systemd/man/systemd-logind.service.html) is a system service that
manages user logins.

Monitors number of sessions and users as reported by the `org.freedesktop.login1` DBus API.

## Requirements

- Works only on Linux systems.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/logind/metrics.csv) for a list of
metrics.

## Configuration

No configuration required.

## Troubleshooting

To troubleshoot issues with the `logind` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m logind
  ```
