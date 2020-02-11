# ActiveMQ monitoring with Netdata

[`ActiveMQ`](https://activemq.apache.org/) is an open source message broker written in Java together with a full Java Message Service client.

This plugin collects queues and topics metrics using ActiveMQ Console API.

## Charts

It produces following charts per queue and per topic:

-   Messages in `messages/s`
-   Unprocessed Messages in `messages`
-   Consumers in `consumers`

## Configuration

Edit the `go.d/activemq.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at
`/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/activemq.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8161
    webadmin: admin
    max_queues: 100
    max_topics: 100
    queues_filter: '!sandr* *'
    topics_filter: '!sandr* *'
    
  - name: remote
    url: http://203.0.113.10:8161
    webadmin: admin
```

For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/activemq.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m activemq
