# Elasticsearch and OpenSearch collector

This module monitors one or more Elasticsearch or OpenSearch instances, depending on your configuration.

[`Elasticsearch`](https://www.elastic.co/elasticsearch/) is a search engine based on the Lucene library. The original Elasticsearch project 
was continued as an open-source project called [OpenSearch](https://opensearch.org/) by Amazon. 

Used endpoints:

- Local node metrics: `/_nodes/_local/stats`
- Local node indices' metrics: `/_cat/indices?local=true`
- Cluster health metrics: `/_cluster/health`
- Cluster metrics: `/_cluster/stats`

Each endpoint can be enabled/disabled in the module configuration file.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/elasticsearch/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/elasticsearch.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/elasticsearch.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9200

  - name: remote
    url: http://203.0.113.0:9200
```

For all available options, see the Elasticsearch
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/elasticsearch.conf).

OpenSearch by default has plugin security enabled (controlled by `plugins.security.disabled` in `/etc/opensearch/opensearch.yml`).
You can provide the required username and password in the OpenSearch-specific job, seen in `go.d/elasticsearch.conf`:
```
# opensearch
  - name: local
    url: https://127.0.0.1:9200
    tls_skip_verify: yes
    username: admin
    password: admin
```

## Troubleshooting

To troubleshoot issues with the `elasticsearch` collector, run the `go.d.plugin` with the debug option enabled. The
output should give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m elasticsearch
  ```
