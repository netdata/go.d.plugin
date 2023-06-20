# systemd-logind collector

## Overview

[systemd-logind](https://www.freedesktop.org/software/systemd/man/systemd-logind.service.html) is a system service that
manages user logins.

This collector monitors number of sessions and users as reported by the `org.freedesktop.login1` DBus API.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                |                 Dimensions                  |   Unit   |
|-----------------------|:-------------------------------------------:|:--------:|
| logind.sessions       |                remote, local                | sessions |
| logind.sessions_type  |          console, graphical, other          | sessions |
| logind.sessions_state |           online, closing, active           | sessions |
| logind.users_state    | offline, closing, online, lingering, active |  users   |

## Setup

### Prerequisites

No action required.

### Configuration

No configuration required.

## Troubleshooting

### Debug mode

To troubleshoot issues with the `logind` collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

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
