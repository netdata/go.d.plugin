<!--
title: "SNMP device monitoring with Netdata"
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/snmp/README.md"
sidebar_label: "SNMP"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Remotes"
-->

# SNMP device collector

Collects data from any SNMP device and uses the [gosnmp](https://github.com/gosnmp/gosnmp) package.

It supports:

- all SNMP versions: SNMPv1, SNMPv2c and SNMPv3.
- any number of SNMP devices.
- each SNMP device can be used to collect data for any number of charts.
- each chart may have any number of dimensions.
- each SNMP device may have a different update frequency.
- each SNMP device will accept one or more batches to report values (you can set `max_request_size` per SNMP server, to
  control the size of batches).

## Configuration

Edit the `go.d/snmp.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory), which is
typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/snmp.conf
```

The configuration file is a list of data collection jobs. Jobs allow you to collect values from multiple sources, each
source will have its own set of charts.

Generally the format is:

```yaml
jobs:
  - name: name1
    ...  # other configuration parameters
  - name: name2
    ...  # other configuration parameters
  - name: name3
    ...  # other configuration parameters
```

### Job configuration parameters

| Parameter                    | Default value  | Description                                                                                                      |
|------------------------------|:--------------:|------------------------------------------------------------------------------------------------------------------|
| name                         |       -        | the data collection job name                                                                                     |
| update_every                 |       10       | the update frequency for each target, in seconds                                                                 |
| hostname                     |   127.0.0.1    | the target ipv4 address                                                                                          |
| community                    |     public     | SNMPv1/2 community string                                                                                        |
| options.version              |       2        | SNMP version                                                                                                     |
| options.port                 |      161       | the target port                                                                                                  |
| options.retries              |       1        | the number of retries to attempt                                                                                 |
| options.timeout              |       10       | the timeout for one SNMP request/response                                                                        |
| options.max_request_size     |       60       | the maximum number of oids allowed in one one SNMP request                                                       |
| user.name                    |       -        | the SNMPv3 user name                                                                                             |
| user.level                   |       -        | the security level of SNMPv3 messages                                                                            |
| user.auth_proto              |       -        | the authentication protocol for SNMPv3 messages                                                                  |
| user.auth_key                |       -        | the authentication protocol pass phrase                                                                          |
| user.priv_proto              |       -        | the privacy protocol for SNMPv3 messages                                                                         |
| user.priv_key                |       -        | the privacy protocol pass phrase                                                                                 |
| charts                       |       []       | the list of charts                                                                                               |
| charts.id                    |       -        | is used to uniquely identify the chart                                                                           |
| charts.title                 | Untilted chart | the text above the chart                                                                                         |
| charts.units                 |      num       | the label of the vertical axis of the chart                                                                      |
| charts.family                |   charts.id    | the name of the dashboard submenu under which each chart will be displayed                                       |
| charts.type                  |      line      | the chart type (one of line, area or stacked)                                                                    |
| charts.priority              |     70000      | the priority of the chart as rendered on the web page                                                            |
| charts.multiply_range        |       []       | is used when you need to define many charts [using incremental OIDs](#example-using-chartsmultiply_range-option) |
| charts.dimensions            |       []       | the list of chart dimensions                                                                                     |
| charts.dimensions.oid        |       -        | the OID path to the metric you [want to collect](#finding-oids)                                                  |
| charts.dimensions.name       |       -        | the name of the dimension as it will appear at the legend of the chart                                           |
| charts.dimensions.algorithm  |    absolute    | the dimension algorithm (one of absolute, incremental)                                                           |
| charts.dimensions.multiplier |       1        | the value to multiply the collected value, applied to convert it properly to units                               |
| charts.dimensions.divisor    |       1        | the value to divide the collected value, applied to convert it properly to units                                 |

### Example: Using SNMPv1/2

In this example:

- the SNMP device is `192.0.2.1`.
- the SNMP version is `2`.
- the SNMP community is `public`.
- we will update the values every 10 seconds.
- we define 2 charts `bandwidth_port1` and `bandwidth_port2`, each having 2 dimensions: `in` and `out`.

> **SNMPv1**: just set `options.version` to 1.

> If you have multiple devices see how to [simplify the configuration](#multiple-devices-with-a-common-configuration).

```yaml
jobs:
  - name: switch
    update_every: 10
    hostname: "192.0.2.1"
    community: public
    options:
      version: 2
    charts:
      - id: "bandwidth_port1"
        title: "Switch Bandwidth for port 1"
        units: "kilobits/s"
        type: "area"
        family: "ports"
        dimensions:
          - name: "in"
            oid: "1.3.6.1.2.1.2.2.1.10.1"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: "1.3.6.1.2.1.2.2.1.16.1"
            multiplier: -8
            divisor: 1000
      - id: "bandwidth_port2"
        title: "Switch Bandwidth for port 2"
        units: "kilobits/s"
        type: "area"
        family: "ports"
        dimensions:
          - name: "in"
            oid: "1.3.6.1.2.1.2.2.1.10.2"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: "1.3.6.1.2.1.2.2.1.16.2"
            multiplier: -8
            divisor: 1000
```

Note that in this example, the algorithm chosen is `incremental`, because the collected values show the total number of bytes transferred,
which we need to transform into kbps. To chart gauges (e.g. temperature), use `absolute` instead. 

### Example: Using SNMPv3

To use SNMPv3:

- use `user` instead of `community`.
- set `options.version` to 3.

The rest of the configuration is the same as in the SNMPv1/2 [example](#example-using-snmpv12).

> If you have multiple devices see how to [simplify the configuration](#multiple-devices-with-a-common-configuration).

```yaml
jobs:
  - name: switch
    update_every: 10
    hostname: "192.0.2.1"
    options:
      version: 3
    user:
      name: "username"
      level: "authPriv"
      auth_proto: "sha256"
      auth_key: "auth_protocol_passphrase"
      priv_proto: "aes256"
      priv_key: "priv_protocol_passphrase"
```

#### SNMPv3 message authentication and privacy configuration options

The security of an SNMPv3 message as per RFC 3414 (`user.level`):

| String value | Int value | Description                              |
|:------------:|:---------:|------------------------------------------|
|     none     |     1     | no message authentication or encryption  |
|  authNoPriv  |     2     | message authentication and no encryption |
|   authPriv   |     3     | message authentication and encryption    |

The digest algorithm for SNMPv3 messages that require authentication (`user.auth_proto`):

| String value | Int value | Description                               |
|:------------:|:---------:|-------------------------------------------|
|     none     |     1     | no message authentication                 |
|     md5      |     2     | MD5 message authentication (HMAC-MD5-96)  |
|     sha      |     3     | SHA message authentication (HMAC-SHA-96)  |
|    sha224    |     4     | SHA message authentication (HMAC-SHA-224) |
|    sha256    |     5     | SHA message authentication (HMAC-SHA-256) |
|    sha384    |     6     | SHA message authentication (HMAC-SHA-384) |
|    sha512    |     7     | SHA message authentication (HMAC-SHA-512) |

The encryption algorithm for SNMPv3 messages that require privacy (`user.priv_proto`):

| String value | Int value | Description                                                             |
|:------------:|:---------:|-------------------------------------------------------------------------|
|     none     |     1     | no message encryption                                                   |
|     des      |     2     | ES encryption (CBC-DES)                                                 |
|     aes      |     3     | 128-bit AES encryption (CFB-AES-128)                                    |
|    aes192    |     4     | 192-bit AES encryption (CFB-AES-192) with "Blumenthal" key localization |
|    aes256    |     5     | 256-bit AES encryption (CFB-AES-256) with "Blumenthal" key localization |
|   aes192c    |     6     | 192-bit AES encryption (CFB-AES-192) with "Reeder" key localization     |
|   aes256c    |     7     | 256-bit AES encryption (CFB-AES-256) with "Reeder" key localization     |

### Example: Using `charts.multiply_range` option

If you need to define many charts using incremental OIDs, you can use the `charts.multiply_range` option.

This is like the SNMPv1/2 [example](#example-using-snmpv12), but the option will multiply the current chart from 1 to 24
inclusive, producing 24 charts in total for the 24 ports of the switch `192.0.2.1`.

Each of the 24 new charts will have its id (1-24) appended at:

- its chart unique `id`, i.e. `bandwidth_port_1` to `bandwidth_port_24`.
- its title, i.e. `Switch Bandwidth for port 1` to `Switch Bandwidth for port 24`.
- its `oid` (for all dimensions), i.e. dimension in will be `1.3.6.1.2.1.2.2.1.10.1` to `1.3.6.1.2.1.2.2.1.10.24`.
- its `priority` will be incremented for each chart so that the charts will appear on the dashboard in this order.

> If you have multiple devices see how to [simplify the configuration](#multiple-devices-with-a-common-configuration).

```yaml
jobs:
  - name: switch
    update_every: 10
    hostname: "192.0.2.1"
    community: public
    options:
      version: 2
    charts:
      - id: "bandwidth_port"
        title: "Switch Bandwidth for port"
        units: "kilobits/s"
        type: "area"
        family: "ports"
        multiply_range: [1, 24]
        dimensions:
          - name: "in"
            oid: "1.3.6.1.2.1.2.2.1.10"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: "1.3.6.1.2.1.2.2.1.16"
            multiplier: -8
            divisor: 1000
```

## Multiple devices with a common configuration

YAML supports [anchors](https://yaml.org/spec/1.2.2/#3222-anchors-and-aliases). The `&` defines and names an anchor, and
the `*` uses it. `<<: *anchor` means, inject the anchor, then extend. We can use anchors to share the common
configuration for multiple devices.

The following example:

- adds an `anchor` to the first job.
- injects (copies) the first job configuration to the second and updates `name` and `hostname` parameters.
- injects (copies) the first job configuration to the third and updates `name` and `hostname` parameters.

```yaml
jobs:
  - &anchor
    name: switch
    update_every: 10
    hostname: "192.0.2.1"
    community: public
    options:
      version: 2
    charts:
      - id: "bandwidth_port1"
        title: "Switch Bandwidth for port 1"
        units: "kilobits/s"
        type: "area"
        family: "ports"
        dimensions:
          - name: "in"
            oid: "1.3.6.1.2.1.2.2.1.10.1"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: "1.3.6.1.2.1.2.2.1.16.1"
            multiplier: -8
            divisor: 1000
  - <<: *anchor
    name: switch2
    hostname: "192.0.2.2"
  - <<: *anchor
    name: switch3
    hostname: "192.0.2.3"
```

## Data collection speed

Keep in mind that many SNMP switches and routers are very slow. They may not be able to report values per second.
`go.d.plugin` reports the time it took for the SNMP device to respond when executed in the [debug](#troubleshooting)
mode.

Also, if many SNMP clients are used on the same SNMP device at the same time, values may be skipped. This is a problem
of the SNMP device, not this collector. In this case, consider reducing the frequency of data collection (
increasing `update_every`).

## Finding OIDs

Use `snmpwalk`, like this:

```sh
snmpwalk -t 20 -O fn -v 2c -c public 192.0.2.1
```

- `-t 20` is the timeout in seconds.
- `-O fn` will display full OIDs in numeric format.
- `-v 2c` is the SNMP version.
- `-c public` is the SNMP community.
- `192.0.2.1` is the SNMP device.

## Troubleshooting

To troubleshoot issues with the `snmp` collector, run the `go.d.plugin` with the debug option enabled. The output should
give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m snmp
  ```
