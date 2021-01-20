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

| Name                                                                                              | Monitors                        |
| :------------------------------------------------------------------------------------------------ | :------------------------------ |
| [activemq](https://github.com/netdata/go.d.plugin/tree/master/modules/activemq)                   | `ActiveMQ`                      |
| [apache](https://github.com/netdata/go.d.plugin/tree/master/modules/apache)                       | `Apache`                        |
| [bind](https://github.com/netdata/go.d.plugin/tree/master/modules/bind)                           | `ISC Bind`                      |
| [cockroachdb](https://github.com/netdata/go.d.plugin/tree/master/modules/cockroachdb)             | `CockroachDB`                   |
| [consul](https://github.com/netdata/go.d.plugin/tree/master/modules/consul)                       | `Consul`                        |
| [coredns](https://github.com/netdata/go.d.plugin/tree/master/modules/coredns)                     | `CoreDNS`                       |
| [couchbase](https://github.com/netdata/go.d.plugin/tree/master/modules/couchbase)                 | `Couchbase`                     |
| [couchdb](https://github.com/netdata/go.d.plugin/tree/master/modules/couchdb)                     | `CouchDB`                       |
| [dnsdist](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsdist)                     | `Dnsdist`                       |
| [dnsmasq](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsmasq)                     | `Dnsmasq DNS Forwarder`         |
| [dnsmasq_dhcp](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsmasq_dhcp)           | `Dnsmasq DHCP`                  |
| [dns_query](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsquery)                  | `DNS Query RTT`                 |
| [docker_engine](https://github.com/netdata/go.d.plugin/tree/master/modules/docker_engine)         | `Docker Engine`                 |
| [dockerhub](https://github.com/netdata/go.d.plugin/tree/master/modules/dockerhub)                 | `Docker Hub`                    |
| [elasticsearch](https://github.com/netdata/go.d.plugin/tree/master/modules/elasticsearch)         | `Elasticsearch`                 |
| [energid](https://github.com/netdata/go.d.plugin/tree/master/modules/energid)                     | `Energi Core`                   |
| [example](https://github.com/netdata/go.d.plugin/tree/master/modules/example)                     | -                               |
| [filecheck](https://github.com/netdata/go.d.plugin/tree/master/modules/filecheck)                 | `Files and Directories`         |
| [fluentd](https://github.com/netdata/go.d.plugin/tree/master/modules/fluentd)                     | `Fluentd`                       |
| [freeradius](https://github.com/netdata/go.d.plugin/tree/master/modules/freeradius)               | `FreeRADIUS`                    |
| [hdfs](https://github.com/netdata/go.d.plugin/tree/master/modules/hdfs)                           | `HDFS`                          |
| [httpcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/httpcheck)                 | `Any HTTP Endpoint`             |
| [isc_dhcpd](https://github.com/netdata/go.d.plugin/tree/master/modules/isc_dhcpd)                 | `ISC dhcpd`                     |
| [k8s_kubelet](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubelet)             | `Kubelet`                       |
| [k8s_kubeproxy](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubeproxy)         | `Kube-proxy`                    |
| [lighttpd](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd)                   | `Lighttpd`                      |
| [lighttpd2](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd2)                 | `Lighttpd2`                     |
| [logstash](https://github.com/netdata/go.d.plugin/tree/master/modules/logstash)                   | `Logstash`                      |
| [mysql](https://github.com/netdata/go.d.plugin/tree/master/modules/mysql)                         | `MySQL`                         |
| [nginx](https://github.com/netdata/go.d.plugin/tree/master/modules/nginx)                         | `NGINX`                         |
| [nginxvts](https://github.com/netdata/go.d.plugin/tree/master/modules/nginxvts)                   | `NGINX VTS`                     |
| [openvpn](https://github.com/netdata/go.d.plugin/tree/master/modules/openvpn)                     | `OpenVPN`                       |
| [phpdaemon](https://github.com/netdata/go.d.plugin/tree/master/modules/phpdaemon)                 | `phpDaemon`                     |
| [phpfpm](https://github.com/netdata/go.d.plugin/tree/master/modules/phpfpm)                       | `PHP-FPM`                       |
| [pihole](https://github.com/netdata/go.d.plugin/tree/master/modules/pihole)                       | `Pi-hole`                       |
| [pika](https://github.com/netdata/go.d.plugin/tree/master/modules/pika)                           | `Pika`                          |
| [prometheus](https://github.com/netdata/go.d.plugin/tree/master/modules/prometheus)               | `Any Prometheus Endpoint`       |
| [portcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/portcheck)                 | `Any TCP Endpoint`              |
| [powerdns](https://github.com/netdata/go.d.plugin/tree/master/modules/powerdns)                   | `PowerDNS Authoritative Server` |
| [powerdns_recursor](https://github.com/netdata/go.d.plugin/tree/master/modules/powerdns_recursor) | `PowerDNS Recursor`             |
| [pulsar](https://github.com/netdata/go.d.plugin/tree/master/modules/portcheck)                    | `Apache Pulsar`                 |
| [rabbitmq](https://github.com/netdata/go.d.plugin/tree/master/modules/rabbitmq)                   | `RabbitMQ`                      |
| [redis](https://github.com/netdata/go.d.plugin/tree/master/modules/redis)                         | `Redis`                         |
| [scaleio](https://github.com/netdata/go.d.plugin/tree/master/modules/scaleio)                     | `Dell EMC ScaleIO`              |
| [solr](https://github.com/netdata/go.d.plugin/tree/master/modules/solr)                           | `Solr`                          |
| [squidlog](https://github.com/netdata/go.d.plugin/tree/master/modules/squidlog)                   | `Squid`                         |
| [springboot2](https://github.com/netdata/go.d.plugin/tree/master/modules/springboot2)             | `Spring Boot2`                  |
| [systemdunits](https://github.com/netdata/go.d.plugin/tree/master/modules/systemdunits)           | `Systemd unit state`            |
| [tengine](https://github.com/netdata/go.d.plugin/tree/master/modules/tengine)                     | `Tengine`                       |
| [unbound](https://github.com/netdata/go.d.plugin/tree/master/modules/unbound)                     | `Unbound`                       |
| [vcsa](https://github.com/netdata/go.d.plugin/tree/master/modules/vcsa)                           | `vCenter Server Appliance`      |
| [vernemq](https://github.com/netdata/go.d.plugin/tree/master/modules/vernemq)                     | `VerneMQ`                       |
| [vsphere](https://github.com/netdata/go.d.plugin/tree/master/modules/vsphere)                     | `VMware vCenter Server`         |
| [web_log](https://github.com/netdata/go.d.plugin/tree/master/modules/weblog)                      | `Apache/NGINX`                  |
| [whoisquery](https://github.com/netdata/go.d.plugin/tree/master/modules/whoisquery)               | `Domain Expiry`                 |
| [wmi](https://github.com/netdata/go.d.plugin/tree/master/modules/wmi)                             | `Windows Machines`              |
| [x509check](https://github.com/netdata/go.d.plugin/tree/master/modules/x509check)                 | `Digital Certificates`          |
| [zookeeper](https://github.com/netdata/go.d.plugin/tree/master/modules/zookeeper)                 | `ZooKeeper`                     |

## Configuration

Edit the `go.d.conf` configuration file using `edit-config` from the Netdata [config
directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d.conf
```

Configurations are written in [YAML](http://yaml.org/).

-   [plugin configuration](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf)
-   [specific module configuration](https://github.com/netdata/go.d.plugin/tree/master/config/go.d)

## Developing

-   Add your module to the [modules dir](https://github.com/netdata/go.d.plugin/tree/master/modules).
-   Import the module in the [main.go](https://github.com/netdata/go.d.plugin/blob/master/cmd/godplugin/main.go).
-   To build it execute `make` from the plugin root dir or `hack/go-build.sh`.
-   Run it in the debug mode `bin/godplugin -d -m <MODULE_NAME>`.
-   Use `make clean` when you are done with testing.

## Troubleshooting

Plugin CLI:

```sh
Usage:
  orchestrator [OPTIONS] [update every]

Application Options:
  -m, --modules=    module name to run (default: all)
  -c, --config-dir= config dir to read
  -w, --watch-path= config path to watch
  -d, --debug       debug mode
  -v, --version     display the version and exit

Help Options:
  -h, --help        Show this help message
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
