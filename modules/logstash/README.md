# Logstash monitoring with Netdata

[`Logstash`](https://www.elastic.co/products/logstash) is an open-source data processing pipeline that allows you to collect, process, and load data into `Elasticsearch`.

This module will monitor one or more `Logstash` instances, depending on your configuration.

## Charts

It produces following charts:

-   JVM Threads in `count`
-   JVM Heap Memory Percentage in `percent`
-   JVM Heap Memory in `KiB`
-   JVM Pool Survivor Memory in `KiB`
-   JVM Pool Old Memory in `KiB`
-   JVM Pool Eden Memory in `KiB`
-   Garbage Collection Count in `counts/s`
-   Time Spent On Garbage Collection in `ms`
-   Uptime in `time`

## Configuration

Edit the `go.d/logstash.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/logstash.conf
```

Here is a simple example for local and remote server:

```yaml
jobs:
  - name: local
    url : http://localhost:9600

  - name: remote
    url : http://203.0.113.10:9600
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/logstash.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m logstash

