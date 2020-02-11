# Solr monitoring with Netdata

[`Solr`](https://lucene.apache.org/solr/) is an open-source enterprise-search platform, written in Java, from the Apache Lucene project.

This module monitors `Solr` request handler statistics.

## Requirement

-   `Solr` version 6.4+

## Charts

It produces the following charts per core:

-   Search Requests in `requests/s`
-   Search Errors in `errors/s`
-   Search Errors By Type in `errors/s`
-   Search Requests Processing Time in `milliseconds`
-   Search Requests Timings in `milliseconds`
-   Search Requests Processing Time Percentile in `milliseconds` 
-   Update Requests in `requests/s`
-   Update Errors in `errors/s`
-   Update Errors By Type in `errors/s` 
-   Update Requests Processing Time in `milliseconds`
-   Update Requests Timings in `milliseconds` 
-   Update Requests Processing Time Percentile in `milliseconds`

## Configuration

Edit the `go.d/solr.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/solr.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://localhost:8983
      
  - name: remote
    url : http://203.0.113.10:8983

```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/solr.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m solr
