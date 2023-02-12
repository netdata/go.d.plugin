<!--
title: "FreeRADIUS monitoring with Netdata"
description: "Monitor the health and performance of FreeRADIUS servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/freeradius/README.md"
sidebar_label: "FreeRADIUS"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Networking"
-->

# FreeRADIUS monitoring with Netdata

[`FreeRADIUS`](https://freeradius.org/) is a modular, high performance free RADIUS suite.

This module will monitor one or more `FreeRADIUS` servers, depending on your configuration.

## Requirements

- `FreeRADIUS` with enabled status feature.

The configuration for the status server is automatically created in the sites-available directory. By default, server is
enabled and can be queried from every client.

To enable status feature do the following:

- `cd sites-enabled`
- `ln -s ../sites-available/status status`
- restart FreeRADIUS server

## Metrics

All metrics have "freeradius." prefix.

| Metric                                | Scope  |                      Dimensions                       |   Units   |
|---------------------------------------|:------:|:-----------------------------------------------------:|:---------:|
| authentication                        | global |                  requests, responses                  | packets/s |
| authentication_access_responses       | global |             accepts, rejects, challenges              | packets/s |
| bad_authentication                    | global | dropped, duplicate, invalid, malformed, unknown-types | packets/s |
| proxy_authentication                  | global |                  requests, responses                  | packets/s |
| proxy_authentication_access_responses | global |             accepts, rejects, challenges              | packets/s |
| proxy_bad_authentication              | global | dropped, duplicate, invalid, malformed, unknown-types | packets/s |
| accounting                            | global |                  requests, responses                  | packets/s |
| bad_accounting                        | global | dropped, duplicate, invalid, malformed, unknown-types | packets/s |
| proxy_accounting                      | global |                  requests, responses                  | packets/s |
| proxy_bad_accounting                  | global | dropped, duplicate, invalid, malformed, unknown-types | packets/s |

## Configuration

Edit the `go.d/freeradius.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/freeradius.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    host: 127.0.0.1

  - name: remote
    host: 203.0.113.10
    secret: secret 
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/freeradius.conf).

## Troubleshooting

To troubleshoot issues with the `freeradius` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m freeradius
  ```


