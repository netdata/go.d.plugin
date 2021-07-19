<!--
title: "Go-ethereum monitoring with Netdata"
description: "Monitor the health and performance of your go-ethereum Nodes (Geth) with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/geth/README.md
sidebar_label: "Go-ethereum"
-->

# Geth Monitoring with Netdata

Go Ethereum, written in Google’s Go language, is one of the three original implementations of the Ethereum protocol,
alongside C++ and Python.

Go-Ethereum, and subsequently Geth, are built and maintained by the Ethereum community. It’s open source which means
anyone can contribute to Geth through its [Github](https://github.com/ethereum/go-ethereum).

With Netdata, you can effortlessly monitor your Geth node

## Requirements

Run `geth` with the flag `--metrics`. That will enable the metric server, with default port `6060` and
path `/debug/metrics/prometheus`.

## Charts

This is an initial number of metrics that we chose to collect and organize. It is **very easy** to add more charts based
on the available metrics in the prometheus endpoint. Head over to [Contribute](#contribute) to learn how you can help to
expand this collector.

- Chaindata:
    - total read/write for the session
    - read/write per second
- Transaction Pool
    - Pending
    - Queued
- Peer-to-Peer
    - bandwidth per second (ingress/egress)
    - number of peers
    - serves/dials calls per second
- rpc calls
    - successful/failed per second
- reorgs
    - Total number of executed reorgs
    - Total number of added/removed blocks due to reorg
- number of active goroutines
- chainhead
    - block, receipt and header. If block = header, then Geth node is fully synced.

## Contribute

We have started
a [topic](https://community.netdata.cloud/t/lets-build-a-golang-collector-for-monitoring-ethereum-full-nodes/1426) on
our community forums about this collector.

**The best contribution you can make is to tell us what metrics you want to see and how they should be organized (e.g
what charts to make).**

As you can read in the topic, it's trivial to add more metrics from the prometheus endpoint and create the relevant
charts. The hard part is the domain expertise that we don't have, but you, as a user, have.

The second best contribution you can make is to tell us what alerts we should be shipping as defaults for this
collector. For example, we are shipping an alert about the node being in sync (or not). We simply compare the
chainhead `block` and `header` values.

If you are proficient in Golang, visit the topic and make a PR yourself to the collector. We will happily help you to
merge it and have your code being shipped with **every** Netdata agent.

## Configuration

Edit the `go.d/geth.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/geth.conf
```

Needs only `url` to `geth` metrics endpoint. Here is an example for 2 instances:

```yaml
jobs:
  - name: geth_node_1
    url: http://203.0.113.10:6060/debug/metrics/prometheus

  - name: geth_node_2
    url: http://203.0.113.11:6060/debug/metrics/prometheus
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/geth.conf).

## Troubleshooting

To troubleshoot issues with the `geth` collector, run the `go.d.plugin` with the debug option enabled. The output should
give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m geth
```
