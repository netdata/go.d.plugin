# consul

[`Consul`](https://www.consul.io/) is a service networking solution to connect and secure services across any runtime platform and public or private cloud.

This module monitors `Consul` health checks.

## Charts

It produces the following charts:

-   Service Checks in `status`

-   Unbound Checks in `status`

## configuration

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:8500
      
  - name: remote
    url : http://203.0.113.10:8500
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/consul.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m consul
