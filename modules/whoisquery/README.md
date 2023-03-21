<!--
title: "Whois domain expiry monitoring with Netdata"
description: "Monitor the health and performance of domain expiry with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/whoisquery/README.md"
sidebar_label: "Whois domain expiry"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# Whois domain expiry collector

This collector module checks the remaining time until a domain is expired.

## Metrics

All metrics have "whoisquery." prefix.

Labels per scope:

- global: domain.

| Metric                | Scope  | Dimensions |  Units  |
|-----------------------|:------:|:----------:|:-------:|
| time_until_expiration | global |   expiry   | seconds |

## Configuration

Edit the `go.d/whoisquery.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/whoisquery.conf
```

Needs only `source`.

Use `days_until_expiration_warning` and `days_until_expiration_critical` for each job to indicate the expiry warning and
critical days. The default values are 90 for warning, and 30 days for critical.

Here is an example:

```yaml
update_every: 60

jobs:
  - name: my_site
    source: my_site.com

  - name: my_another_site
    source: my_another_site.com
    days_until_expiration_critical: 20

```

For all available options and defaults please, see the
module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/whoisquery.conf).

## Troubleshooting

To troubleshoot issues with the `whoisquery` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m whoisquery
  ```
