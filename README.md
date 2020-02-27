# go.d.plugin

[![CircleCI](https://circleci.com/gh/netdata/go.d.plugin.svg?style=svg)](https://circleci.com/gh/netdata/go.d.plugin)

`go.d.plugin` is a `netdata` external plugin. It is an **orchestrator** for data collection modules written in `go`.

1. It runs as an independent process `ps fax` shows it.
2. It is started and stopped automatically by `netdata`.
3. It communicates with `netdata` via a unidirectional pipe (sending data to the `netdata` daemon).
4. Supports any number of data collection [modules](https://github.com/netdata/go.d.plugin/tree/master/modules).
5. Allows each [module](https://github.com/netdata/go.d.plugin/tree/master/modules) to have any number of data collection **jobs**.

## Install

Shipped with `netdata`.

## Contributing
If you have time and willing to help, there are a lof of ways to contribute:

-   Fix and [report bugs](https://github.com/netdata/go.d.plugin/issues/new)
-   [Review code and feature proposals](https://github.com/netdata/go.d.plugin/pulls)
-   [Contribute modules](https://github.com/netdata/go.d.plugin/blob/master/CONTRIBUTING.md) (wip, module interface may be changed soon)

## Available modules

-   [activemq](https://github.com/netdata/go.d.plugin/tree/master/modules/activemq)
-   [apache](https://github.com/netdata/go.d.plugin/tree/master/modules/apache) *
-   [bind](https://github.com/netdata/go.d.plugin/tree/master/modules/bind) *
-   [cockroachdb](https://github.com/netdata/go.d.plugin/tree/master/modules/cockroachdb)
-   [consul](https://github.com/netdata/go.d.plugin/tree/master/modules/consul)
-   [coredns](https://github.com/netdata/go.d.plugin/tree/master/modules/coredns)
-   [dnsmasq_dhcp](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsmasq_dhcp)
-   [dns_query](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsquery)
-   [docker_engine](https://github.com/netdata/go.d.plugin/tree/master/modules/docker_engine)
-   [dockerhub](https://github.com/netdata/go.d.plugin/tree/master/modules/dockerhub)
-   [example](https://github.com/netdata/go.d.plugin/tree/master/modules/example) *
-   [fluentd](https://github.com/netdata/go.d.plugin/tree/master/modules/fluentd)
-   [freeradius](https://github.com/netdata/go.d.plugin/tree/master/modules/freeradius) *
-   [hdfs](https://github.com/netdata/go.d.plugin/tree/master/modules/hdfs)
-   [httpcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/httpcheck)
-   [k8s_kubelet](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubelet)
-   [k8s_kubeproxy](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubeproxy)
-   [lighttpd](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd) *
-   [lighttpd2](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd2)
-   [logstash](https://github.com/netdata/go.d.plugin/tree/master/modules/logstash)
-   [mysql](https://github.com/netdata/go.d.plugin/tree/master/modules/mysql) *
-   [nginx](https://github.com/netdata/go.d.plugin/tree/master/modules/nginx) *
-   [openvpn](https://github.com/netdata/go.d.plugin/tree/master/modules/openvpn)
-   [phpdaemon](https://github.com/netdata/go.d.plugin/tree/master/modules/phpdaemon)
-   [phpfpm](https://github.com/netdata/go.d.plugin/tree/master/modules/phpfpm) *
-   [pihole](https://github.com/netdata/go.d.plugin/tree/master/modules/pihole)
-   [portcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/portcheck)
-   [rabbitmq](https://github.com/netdata/go.d.plugin/tree/master/modules/rabbitmq) *
-   [scaleio](https://github.com/netdata/go.d.plugin/tree/master/modules/scaleio)
-   [solr](https://github.com/netdata/go.d.plugin/tree/master/modules/solr)
-   [squidlog](https://github.com/netdata/go.d.plugin/tree/master/modules/squidlog) *
-   [springboot2](https://github.com/netdata/go.d.plugin/tree/master/modules/springboot2)
-   [tengine](https://github.com/netdata/go.d.plugin/tree/master/modules/tengine)
-   [unbound](https://github.com/netdata/go.d.plugin/tree/master/modules/unbound)
-   [vcsa](https://github.com/netdata/go.d.plugin/tree/master/modules/vcsa)
-   [vernemq](https://github.com/netdata/go.d.plugin/tree/master/modules/vernemq)
-   [vsphere](https://github.com/netdata/go.d.plugin/tree/master/modules/vsphere)
-   [web_log](https://github.com/netdata/go.d.plugin/tree/master/modules/weblog) *
-   [wmi](https://github.com/netdata/go.d.plugin/tree/master/modules/wmi)
-   [x509check](https://github.com/netdata/go.d.plugin/tree/master/modules/x509check)
-   [zookeeper](https://github.com/netdata/go.d.plugin/tree/master/modules/zookeeper)

`*` - disabled by default.

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
```
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
```
# become user netdata
sudo su -s /bin/bash netdata

# run plugin in debug mode
./go.d.plugin -d -m <module name>
```

Change `<module name>` to the module name you want to debug.
See the [whole list](#available-modules) of available modules.
