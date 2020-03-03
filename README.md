# go.d.plugin

[![CircleCI](https://circleci.com/gh/netdata/go.d.plugin.svg?style=svg)](https://circleci.com/gh/netdata/go.d.plugin)

`go.d.plugin` is a `Netdata` external plugin. It is an **orchestrator** for data collection modules written in `go`.

1. It runs as an independent process `ps fax` shows it.
2. It is started and stopped automatically by `Netdata`.
3. It communicates with `Netdata` via a unidirectional pipe (sending data to the `Netdata` daemon).
4. Supports any number of data collection [modules](https://github.com/netdata/go.d.plugin/tree/master/modules).
5. Allows each [module](https://github.com/netdata/go.d.plugin/tree/master/modules) to have any number of data collection **jobs**.

## Install

Shipped with [`Netdata`](https://github.com/netdata/netdata).

## Contributing

If you have time and willing to help, there are a lof of ways to contribute:

-   Fix and [report bugs](https://github.com/netdata/go.d.plugin/issues/new)
-   [Review code and feature proposals](https://github.com/netdata/go.d.plugin/pulls)
-   [Contribute modules](https://github.com/netdata/go.d.plugin/blob/master/CONTRIBUTING.md) (wip, module interface may be changed soon)

## Available modules

| Name                                                                                      | Monitors                   | Disabled|
| :---------------------------------------------------------------------------------------- | :------------------------- | :-------|
| [activemq](https://github.com/netdata/go.d.plugin/tree/master/modules/activemq)           | `ActiveMQ`                 |         |
| [apache](https://github.com/netdata/go.d.plugin/tree/master/modules/apache)               | `Apache`                   | yes     |
| [bind](https://github.com/netdata/go.d.plugin/tree/master/modules/bind)                   | `ISC Bind`                 | yes     |
| [cockroachdb](https://github.com/netdata/go.d.plugin/tree/master/modules/cockroachdb)     | `CockroachDB`              |         | 
| [consul](https://github.com/netdata/go.d.plugin/tree/master/modules/consul)               | `Consul`                   |         |
| [coredns](https://github.com/netdata/go.d.plugin/tree/master/modules/coredns)             | `CoreDNS`                  |         |
| [dnsmasq_dhcp](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsmasq_dhcp)   | `Dnsmasq`                  |         |
| [dns_query](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsquery)          | `DNS Query RTT`            |         |
| [docker_engine](https://github.com/netdata/go.d.plugin/tree/master/modules/docker_engine) | `Docker Engine`            |         |
| [dockerhub](https://github.com/netdata/go.d.plugin/tree/master/modules/dockerhub)         | `Docker Hub`               |         |
| [example](https://github.com/netdata/go.d.plugin/tree/master/modules/example)             | -                          | yes     | 
| [fluentd](https://github.com/netdata/go.d.plugin/tree/master/modules/fluentd)             | `Fluentd`                  |         |
| [freeradius](https://github.com/netdata/go.d.plugin/tree/master/modules/freeradius)       | `FreeRADIUS`               | yes     |
| [hdfs](https://github.com/netdata/go.d.plugin/tree/master/modules/hdfs)                   | `HDFS`                     |         |
| [httpcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/httpcheck)         | `Any HTTP Endpoint`        |         |
| [k8s_kubelet](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubelet)     | `Kubelet`                  |         |
| [k8s_kubeproxy](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubeproxy) | `Kube-proxy`               |         |
| [lighttpd](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd)           | `Lighttpd`                 | yes     |
| [lighttpd2](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd2)         | `Lighttpd2`                |         |
| [logstash](https://github.com/netdata/go.d.plugin/tree/master/modules/logstash)           | `Logstash`                 |         |
| [mysql](https://github.com/netdata/go.d.plugin/tree/master/modules/mysql)                 | `MySQL`                    | yes     |
| [nginx](https://github.com/netdata/go.d.plugin/tree/master/modules/nginx)                 | `NGINX`                    | yes     |
| [openvpn](https://github.com/netdata/go.d.plugin/tree/master/modules/openvpn)             | `OpenVPN`                  | yes     |
| [phpdaemon](https://github.com/netdata/go.d.plugin/tree/master/modules/phpdaemon)         | `phpDaemon`                |         |
| [phpfpm](https://github.com/netdata/go.d.plugin/tree/master/modules/phpfpm)               | `PHP-FPM`                  | yes     |
| [pihole](https://github.com/netdata/go.d.plugin/tree/master/modules/pihole)               | `Pi-hole`                  |         |
| [portcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/portcheck)         | `Any TCP Endpoint`         |         |
| [rabbitmq](https://github.com/netdata/go.d.plugin/tree/master/modules/rabbitmq)           | `RabbitMQ`                 | yes     |
| [scaleio](https://github.com/netdata/go.d.plugin/tree/master/modules/scaleio)             | `Dell EMC ScaleIO`         |         |
| [solr](https://github.com/netdata/go.d.plugin/tree/master/modules/solr)                   | `Solr`                     |         |
| [squidlog](https://github.com/netdata/go.d.plugin/tree/master/modules/squidlog)           | `Squid`                    | yes     |
| [springboot2](https://github.com/netdata/go.d.plugin/tree/master/modules/springboot2)     | `Spring Boot2`             |         |
| [tengine](https://github.com/netdata/go.d.plugin/tree/master/modules/tengine)             | `Tengine`                  |         |
| [unbound](https://github.com/netdata/go.d.plugin/tree/master/modules/unbound)             | `Unbound`                  |         |
| [vcsa](https://github.com/netdata/go.d.plugin/tree/master/modules/vcsa)                   | `vCenter Server Appliance` |         |
| [vernemq](https://github.com/netdata/go.d.plugin/tree/master/modules/vernemq)             | `VerneMQ`                  |         | 
| [vsphere](https://github.com/netdata/go.d.plugin/tree/master/modules/vsphere)             | `VMware vCenter Server`    |         |
| [web_log](https://github.com/netdata/go.d.plugin/tree/master/modules/weblog)              | `Apache/NGINX`             | yes     |
| [wmi](https://github.com/netdata/go.d.plugin/tree/master/modules/wmi)                     | `Windows Machines`         |         |
| [x509check](https://github.com/netdata/go.d.plugin/tree/master/modules/x509check)         | `Digital Certificates`     |         |
| [zookeeper](https://github.com/netdata/go.d.plugin/tree/master/modules/zookeeper)         | `ZooKeeper`                |         |

## Why disabled? How to enable?

We are in process of migrating collectors from `python` to `go`.

Configurations are incompatible. All rewritten in `go` modules are disabled by default.
This is a temporary solution, we are working on it.

To enable module please do the following:

-   explicitly disable python module in `python.d.conf`
-   explicitly enable go module in `go.d.conf`
-   move python module jobs to go module configuration file (change syntax, see go module configuration file for details).
-   restart `netdata.service`

If case of problems:

-   check `error.log` for module related errors (`grep <module name> error.log`)
-   run plugin in [debug mode](#troubleshooting)

## Configuration

`go.d.plugin` itself can be configured using the configuration file `/etc/netdata/go.d.conf`
(to edit it on your system run `/etc/netdata/edit-config go.d.conf`). This file is a BASH script.

Configurations are written in [YAML](http://yaml.org/).

-   [plugin configuration](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf)
-   [specific module configuration](https://github.com/netdata/go.d.plugin/tree/master/config/go.d)

## Troubleshooting

Plugin CLI:

```sh
Usage:
  go.d.plugin [OPTIONS] [update every]

Application Options:
  -d, --debug    debug mode
  -m, --modules= modules name (default: all)
  -c, --config=  config dir

Help Options:
  -h, --help     Show this help message
```

To debug specific module:

```sh
# become user netdata
sudo su -s /bin/bash netdata

# run plugin in debug mode
./go.d.plugin -d -m <module name>
```

Change `<module name>` to the module name you want to debug.
See the [whole list](#available-modules) of available modules.
