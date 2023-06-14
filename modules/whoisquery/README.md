# Domain name expiry collector

## Overview

A domain name (often simply called a domain) is an easy-to-remember name
that's associated with a physical IP address on the Internet.

This collector monitors the remaining time before the domain expires.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### domain

These metrics refer to the configured source.

Labels:

| Label  | Description       |
|--------|-------------------|
| domain | Configured source |

Metrics:

| Metric                           | Dimensions |  Unit   |
|----------------------------------|:----------:|:-------:|
| whoisquery.time_until_expiration |   expiry   | seconds |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/whoisquery.conf`.

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
sudo ./edit-config go.d/whoisquery.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|              Name              | Description                                                        | Default | Required |
|:------------------------------:|--------------------------------------------------------------------|:-------:|:--------:|
|          update_every          | Data collection frequency.                                         |    1    |          |
|      autodetection_retry       | Re-check interval in seconds. Zero means not to schedule re-check. |    0    |          |
|             source             | Domain address.                                                    |         |   yes    |
| days_until_expiration_warning  | Number of days before the alarm status is warning.                 |   30    |          |
| days_until_expiration_critical | Number of days before the alarm status is critical.                |   15    |          |
|            timeout             | The query timeout in seconds.                                      |    5    |          |

</details>

#### Examples

##### Basic

Basic configuration example
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my_site
    source: my_site.com
```

</details>

##### Multi-instance

> **Note**: When you define more than one job, their names must be unique.

Check the expiration status of the multiple domains.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my_site1
    source: my_site1.com

  - name: my_site2
    source: my_site2.com
```

</details>

## Troubleshooting

### Debug mode

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
