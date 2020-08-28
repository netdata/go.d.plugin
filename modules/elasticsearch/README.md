<!--
title: "Elasticsearch monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/elasticsearch/README.md
sidebar_label: "Elasticsearch"
-->

# Elasticsearch endpoint monitoring with Netdata

[`Elasticsearch`](https://www.elastic.co/elasticsearch/) is a search engine based on the Lucene library.

This module monitors one or more `Elasticsearch` instances, depending on your configuration.

## Charts


## Configuration

Edit the `go.d/elasticsearch.conf` configuration file using `edit-config` from the Agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/elasticsearch.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: node_exporter_local
    url: http://127.0.0.1:9200

  - name: win10
    url: http://203.0.113.0:9200
```

For all available options, see the Elasticsearch collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/elasticsearch.conf).


## Troubleshooting

To troubleshoot issues with the Elasticsearch collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m elasticsearch
```
