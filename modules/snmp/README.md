<!--
title: "SNMP device monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/snmp/README.md
sidebar_label: "SNMP"
-->

# SNMP device monitoring with Netdata

Collects data from any SNMP device and uses the [gosnmp](https://github.com/gosnmp/gosnmp)package.

It supports:

- all SNMP versions: SNMPv1, SNMPv2c and SNMPv3
- any number of SNMP devices
- each SNMP device can be used to collect data for any number of charts
- each chart may have any number of dimensions
- each SNMP device may have a different update frequency
- each SNMP device will accept one or more batches to report values (you can set `max_request_size` per SNMP server, to
  control the size of batches).

## Configuration

Create a config file `/etc/netdata/go.d/snmp.conf`. Following is an example of SNMPv3 config.

In this example:

- the SNMP device is `10.11.12.8`.
- the SNMP version is `3` (defined under `options.version`)
- we will update the values every 3 seconds (`update_every: 3` for the server `10.11.12.8`).
- since we are using SNMPv3, we need to define authentication mechanisms under `user`
- we define 1 chart `example` having 2 dimensions: `in` and `out`.

```yaml
jobs:
  - name: local
    update_every: 3
    hostname: "10.11.12.8"
    options:
      port: 1161
      version: 3
      max_request_size: 60
    user:
      name: "username"
      level: 3
      auth_proto: 2
      auth_key: "auth_key"
      priv_proto: 2
      priv_key: "priv_key"
    charts:
      - title: "example"
        priority: 1
        units: "kilobits/s"
        type: "area"
        family: "lan"
        dimensions:
          - name: "in"
            oid: ".1.3.6.1.2.1.2.2.1.10.2"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: ".1.3.6.1.2.1.2.2.1.16.2"
            multiplier: -8
            divisor: 1000
```

`update_every` is the update frequency for each server, in seconds.

`family` sets the name of the submenu of the dashboard each chart will appear under.

`multiplier` and `divisor` are passed by the plugin to the Netdata daemon and are applied to the metric to convert it
properly to `units`.

If many charts need to be defined using incremental OIDs, `multiply_range` can be used as in following example:

```yaml
jobs:
  - name: local
    update_every: 3
    hostname: "10.11.12.8"
    options:
      port: 1161
      version: 3
      max_request_size: 60
    user:
      name: "username"
      level: 3
      auth_proto: 2
      auth_key: "auth_key"
      priv_proto: 2
      priv_key: "priv_key"
    charts:
      - title: "example"
        priority: 1
        units: "kilobits/s"
        type: "area"
        family: "lan"
        multiply_range: [1,5]
        dimensions:
          - name: "in"
            oid: ".1.3.6.1.2.1.2.2.1.10"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: ".1.3.6.1.2.1.2.2.1.16"
            multiplier: -8
            divisor: 1000
```

This is like the previous, but with config parameter `charts.multiply_range` given. This will generate charts OID index
appended from `1` to `5` producing 5 charts in total for the 5 ports of the switch `10.11.12.8`.

Each of the 5 new charts will have its id (1-5) appended at:

1. its `charts.title`, i.e. `example_1` to `example_5`
2. its `charts.dimensions.oid` (for all dimensions), i.e. dimension `in` will be `1.3.6.1.2.1.2.2.1.10.1`
   to `1.3.6.1.2.1.2.2.1.10.5`
3. its priority (which will be incremented for each chart so that the charts will appear on the dashboard in this order)

The `options` given for each server, are:

- `port` - UDP port to send requests too. Defaults to `161`.
- `retries` - number of times to re-send a request. Defaults to `1`.
- `timeout` - number of milliseconds to wait for a response before re-trying or failing. Defaults to `5000`.
- `version` - either `1` (v1) or  `2` (v2) or `3` (v3). Defaults to `2`.
- `max_request_size` limits the maximum number of OIDs that will be requested in a single call. The default is 60.

## SNMPv1/2

To use SNMPv1 or 2:

- use `community` instead of `user`
- set `options.version` to 1 or 2

Example:

```yaml
jobs:
  - name: local
    update_every: 3
    hostname: "10.11.12.8"
    options:
      port: 1161
      version: 2
      max_request_size: 60
    community: "public"
    charts:
      - title: "example"
        priority: 1
        units: "kilobits/s"
        type: "area"
        family: "lan"
        multiply_range: [1,5]
        dimensions:
          - name: "in"
            oid: ".1.3.6.1.2.1.2.2.1.10"
            algorithm: "incremental"
            multiplier: 8
            divisor: 1000
          - name: "out"
            oid: ".1.3.6.1.2.1.2.2.1.16"
            multiplier: -8
            divisor: 1000
```

## SNMPv3

To use SNMPv3:

- use `user` instead of `community`
- set `options.version` to 3

Example: refer to the first example in this README.

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

## Testing the configuration

To test it, you can run:

```sh
/usr/libexec/netdata/plugins.d/go.d.plugin -d 1 -m snmp
```

The above will run it on your console and you will be able to see what Netdata sees, but also errors. If it works,
restart Netdata to activate the snmp collector and refresh the dashboard (if your SNMP device responds with a delay, you
may need to refresh the dashboard in a few seconds).

## Data collection speed

Keep in mind that many SNMP switches and routers are very slow. They may not be able to report values per second.
If `go.d.plugin` is executed in `debug` mode, it will report the time it took for the SNMP device to respond.

Also, if many SNMP clients are used on the same SNMP device at the same time, values may be skipped. This is a problem
of the SNMP device, not this collector.

## Finding OIDs

Use `snmpwalk`, like this:

```sh
snmpwalk -t 20 -v 1 -O fn -c public 10.11.12.8
```

- `-t 20` is the timeout in seconds
- `-v 1` is the SNMP version
- `-O fn` will display full OIDs in numeric format
- `-c public` is the SNMP community
- `10.11.12.8` is the SNMP device

Note that `snmpwalk` outputs the OIDs with a dot in front them. Remove this dot when adding OIDs to the configuration
file of this collector.
