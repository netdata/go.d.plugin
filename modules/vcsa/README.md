# vCenter Server Appliance collector

## Overview

The [vCenter Server Appliance](https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.vcsa.doc/GUID-223C2821-BD98-4C7A-936B-7DBE96291BA4.html)
is a preconfigured Linux virtual machine, which is optimized for running VMware vCenter Server® and the associated
services on Linux.

This collector
monitors [health statistics](https://developer.vmware.com/apis/vsphere-automation/latest/appliance/health/) from one or
more vCenter Server Appliance servers, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.
<details>
<summary>See health statuses</summary>
Overall System Health:

| Status  | Description                                                                                                              |
|:-------:|:-------------------------------------------------------------------------------------------------------------------------|
|  green  | All components in the appliance are healthy.                                                                             |
| yellow  | One or more components in the appliance might become overloaded soon.                                                    |
| orange  | One or more components in the appliance might be degraded.                                                               |
|   red   | One or more components in the appliance might be in an unusable status and the appliance might become unresponsive soon. |
|  gray   | No health data is available.                                                                                             |
| unknown | Collector failed to decode status.                                                                                       |

Components Health:

| Status  | Description                                                  |
|:-------:|:-------------------------------------------------------------|
|  green  | The component is healthy.                                    |
| yellow  | The component is healthy, but may have some problems.        |
| orange  | The component is degraded, and may have serious problems.    |
|   red   | The component is unavailable, or will stop functioning soon. |
|  gray   | No health data is available.                                 |
| unknown | Collector failed to decode status.                           |

Software Updates Health:

| Status  | Description                                          |
|:-------:|:-----------------------------------------------------|
|  green  | No updates available.                                |
| orange  | Non-security patches might be available.             |
|   red   | Security patches might be available.                 |
|  gray   | An error retrieving information on software updates. |
| unknown | Collector failed to decode status.                   |

</details>


This scope has no labels.

Metrics:

| Metric                               |                Dimensions                 |  Unit  |
|--------------------------------------|:-----------------------------------------:|:------:|
| vcsa.system_health_status            | green, red, yellow, orange, gray, unknown | status |
| vcsa.applmgmt_health_status          | green, red, yellow, orange, gray, unknown | status |
| vcsa.load_health_status              | green, red, yellow, orange, gray, unknown | status |
| vcsa.mem_health_status               | green, red, yellow, orange, gray, unknown | status |
| vcsa.swap_health_status              | green, red, yellow, orange, gray, unknown | status |
| vcsa.database_storage_health_status  | green, red, yellow, orange, gray, unknown | status |
| vcsa.storage_health_status           | green, red, yellow, orange, gray, unknown | status |
| vcsa.software_packages_health_status |     green, red, orange, gray, unknown     | status |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/vcsa.conf`.

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
sudo ./edit-config go.d/vcsa.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               | Default | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:-------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |    5    |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |    0    |          |
|         url          | Server URL.                                                                                               |         |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |    1    |          |
|       username       | Username for basic HTTP authentication.                                                                   |         |   yes    |
|       password       | Password for basic HTTP authentication.                                                                   |         |   yes    |
|      proxy_url       | Proxy URL.                                                                                                |         |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |         |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |         |          |
|        method        | HTTP request method.                                                                                      |   GET   |          |
|         body         | HTTP request body.                                                                                        |         |          |
|       headers        | HTTP request headers.                                                                                     |         |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |   no    |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |   no    |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |         |          |
|       tls_cert       | Client TLS certificate.                                                                                   |         |          |
|       tls_key        | Client TLS key.                                                                                           |         |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: vcsa1
    url: https://203.0.113.1
    username: admin@vsphere.local
    password: password
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Two instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: vcsa1
    url: https://203.0.113.1
    username: admin@vsphere.local
    password: password

  - name: vcsa2
    url: https://203.0.113.10
    username: admin@vsphere.local
    password: password
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `vcsa` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m vcsa
  ```

