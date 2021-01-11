<!--
title: "Solr monitoring with Netdata"
description: "Monitor the health and performance of Solr search servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/solr/README.md
sidebar_label: "Solr"
-->

# Solr monitoring with Netdata

[`Solr`](https://lucene.apache.org/solr/) is an open-source enterprise-search platform, written in Java, from the Apache
Lucene project.

This module monitors `Solr` request handler statistics.

## Requirement

- `Solr` version 6.4+

## Charts

It produces the following charts per core:

- Search Requests in `requests/s`
- Search Errors in `errors/s`
- Search Errors By Type in `errors/s`
- Search Requests Processing Time in `milliseconds`
- Search Requests Timings in `milliseconds`
- Search Requests Processing Time Percentile in `milliseconds`
- Update Requests in `requests/s`
- Update Errors in `errors/s`
- Update Errors By Type in `errors/s`
- Update Requests Processing Time in `milliseconds`
- Update Requests Timings in `milliseconds`
- Update Requests Processing Time Percentile in `milliseconds`

## Configuration

Edit the `go.d/solr.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/solr.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://localhost:8983

  - name: remote
    url: http://203.0.113.10:8983

```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/solr.conf).

## Troubleshooting

To troubleshoot issues with the `solr` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m solr
```
