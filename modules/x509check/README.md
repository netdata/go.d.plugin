<!--
title: "x509 certificate monitoring with Netdata"
description: "Monitor the health and performance of x509 certificates with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/x509check/README.md"
sidebar_label: "x509 certificates"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# x509 certificate collector

This module checks the time until a x509 certificate expiration and its revocation status.

## Metrics

All metrics have "x509." prefix.

Labels per scope:

- global: source.

| Metric                | Scope  | Dimensions |  Units  |
|-----------------------|:------:|:----------:|:-------:|
| time_until_expiration | global |   expiry   | seconds |
| revocation_status     | global |  revoked   | boolean |

## Configuration

Edit the `go.d/x509check.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/x509check.conf
```

Needs only `source`.

Use `smtp` scheme for smtp servers, `file` for files and `https` or `tcp` for others. Port is mandatory for all non-file
schemes.

Here is an example for 3 sources:

```yaml
update_every: 60

jobs:
  - name: my_site_cert
    source: https://my_site.org:443

  - name: my_file_cert
    source: file:///home/me/cert.pem

  - name: my_smtp_cert
    source: smtp://smtp.my_mail.org:587
```

For all available options and defaults please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/x509check.conf).

## Revocation status

Revocation status check is disabled by default. To enable it set `check_revocation_status` to yes.

```yaml
jobs:
  - name: my_site_cert
    source: https://my_site.org:443
    check_revocation_status: yes
```

## Troubleshooting

To troubleshoot issues with the `x509check` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m x509check
  ```
