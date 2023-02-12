<!--
title: "Energi Core Wallet monitoring with Netdata"
description: "Monitor the health and performance of Energi Core Wallets with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/energid/README.md"
sidebar_label: "Energi Core Wallet"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/Remotes/Devices"
-->

# Energi Core Wallet monitoring with Netdata

This module monitors one or more `Energi Core Wallet` instances, depending on your configuration.

## Requirements

Works only with [Generation 2 wallets](https://docs.energi.software/en/downloads/gen2-core-wallet).

## Metrics

All metrics have "energid." prefix.

| Metric            | Scope  |            Dimensions             |    Units     |
|-------------------|:------:|:---------------------------------:|:------------:|
| blockindex        | global |          blocks, headers          |    count     |
| difficulty        | global |            difficulty             |  difficulty  |
| mempool           | global |        max, usage, tx_size        |    bytes     |
| secmem            | global |     total, used, free, locked     |    bytes     |
| network           | global |            connections            | connections  |
| timeoffset        | global |            timeoffset             |   seconds    |
| utxo_transactions | global | transactions, output_transactions | transactions |

## Configuration

Edit the `go.d/energid.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/energid.conf
```

Needs `url`, `username` and `password`. Here is an example with two jobs:

```yaml
jobs:
  - name: local
    url: 'http://127.0.0.1:9796'
    username: 'netdata'
    password: 'netdata'

  - name: remote
    url: 'http://203.0.113.0:9796'
    username: 'netdata'
    password: 'netdata'
```

For all available options, see the `energid`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/energid.conf).

## Troubleshooting

To troubleshoot issues with the `energid` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m energid
  ```
