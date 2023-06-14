# x509 certificate collector

## Overview

An X.509 certificate is a digital certificate based on the widely accepted International Telecommunications Union (ITU)
X.509 standard.

This collector monitors the time until a x509 certificate expires and its revocation status.

Information about X509 certificates can be collected through a local file, TCP, UDP, HTTPS or SMTP protocols.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### source

These metrics refer to the configured source.

Labels:

| Label  | Description        |
|--------|--------------------|
| source | Configured source. |

Metrics:

| Metric                          | Dimensions |  Unit   |
|---------------------------------|:----------:|:-------:|
| x509check.time_until_expiration |   expiry   | seconds |
| x509check.revocation_status     |  revoked   | boolean |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/x509check.conf`.

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
sudo ./edit-config go.d/x509check.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|              Name              | Description                                                                                               | Default | Required |
|:------------------------------:|-----------------------------------------------------------------------------------------------------------|:-------:|:--------:|
|          update_every          | Data collection frequency.                                                                                |    1    |          |
|      autodetection_retry       | Re-check interval in seconds. Zero means not to schedule re-check.                                        |    0    |          |
|             source             | Certificate source. Allowed schemes: https, tcp, tcp4, tcp6, udp, udp4, udp6, file.                       |         |          |
| days_until_expiration_warning  | Number of days before the alarm status is warning.                                                        |   30    |          |
| days_until_expiration_critical | Number of days before the alarm status is critical.                                                       |   15    |          |
|    check_revocation_status     | Whether to check the revocation status of the certificate.                                                |   no    |          |
|            timeout             | SSL connection timeout.                                                                                   |    2    |          |
|        tls_skip_verify         | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |   no    |          |
|             tls_ca             | Certification authority that the client uses when verifying the server's certificates.                    |         |          |
|            tls_cert            | Client TLS certificate.                                                                                   |         |          |
|            tls_key             | Client TLS key.                                                                                           |         |          |

</details>

#### Examples

##### Website certificate

Website certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my_site_cert
    source: https://my_site.org:443
```

</details>

##### Local file certificate

Local file certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my_file_cert
    source: file:///home/me/cert.pem
```

</details>

##### SMTP certificate

SMTP certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my_smtp_cert
    source: smtp://smtp.my_mail.org:587
```

</details>

##### Multi-instance

> **Note**: When you define more than one job, their names must be unique.

Check the expiration status of the multiple websites' certificates.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: my_site_cert1
    source: https://my_site1.org:443

  - name: my_site_cert2
    source: https://my_site1.org:443

  - name: my_site_cert3
    source: https://my_site3.org:443
```

</details>

## Troubleshooting

### Debug mode

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
