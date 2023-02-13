<!--
title: go.d.plugin
description: "go.d.plugin is an external plugin for Netdata, responsible for running individual data collectors written in Go."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/README.md
sidebar_label: "go.d.plugin"
learn_status: "Published"
learn_topic_type: "Tasks"
learn_rel_path: "Developers/External plugins"
sidebar_position: 1
-->

# go.d.plugin

`go.d.plugin` is a [Netdata](https://github.com/netdata/netdata) external plugin. It is an **orchestrator** for data
collection modules written in `go`.

1. It runs as an independent process (`ps fax` shows it).
2. It is started and stopped automatically by Netdata.
3. It communicates with Netdata via a unidirectional pipe (sending data to the Netdata daemon).
4. Supports any number of data collection [modules](https://github.com/netdata/go.d.plugin/tree/master/modules).
5. Allows each [module](https://github.com/netdata/go.d.plugin/tree/master/modules) to have any number of data
   collection jobs.

## Bug reports, feature requests, and questions

Are welcome! We are using [netdata/netdata](https://github.com/netdata/netdata/) repository for bugs, feature requests,
and questions.

- [GitHub Issues](https://github.com/netdata/netdata/issues/new/choose): report bugs or open a new feature request.
- [GitHub Discussions](https://github.com/netdata/netdata/discussions): ask a question or suggest a new idea.

## Install

Go.d.plugin is shipped with Netdata.

## Available modules

| Name                                                                                                |           Monitors            |
|:----------------------------------------------------------------------------------------------------|:-----------------------------:|
| [activemq](https://github.com/netdata/go.d.plugin/tree/master/modules/activemq)                     |           ActiveMQ            |
| [apache](https://github.com/netdata/go.d.plugin/tree/master/modules/apache)                         |            Apache             |
| [bind](https://github.com/netdata/go.d.plugin/tree/master/modules/bind)                             |           ISC Bind            |
| [cassandra](https://github.com/netdata/go.d.plugin/tree/master/modules/cassandra)                   |           Cassandra           |
| [chrony](https://github.com/netdata/go.d.plugin/tree/master/modules/chrony)                         |            Chrony             |
| [cockroachdb](https://github.com/netdata/go.d.plugin/tree/master/modules/cockroachdb)               |          CockroachDB          |
| [consul](https://github.com/netdata/go.d.plugin/tree/master/modules/consul)                         |            Consul             |
| [coredns](https://github.com/netdata/go.d.plugin/tree/master/modules/coredns)                       |            CoreDNS            |
| [couchbase](https://github.com/netdata/go.d.plugin/tree/master/modules/couchbase)                   |           Couchbase           |
| [couchdb](https://github.com/netdata/go.d.plugin/tree/master/modules/couchdb)                       |            CouchDB            |
| [dnsdist](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsdist)                       |            Dnsdist            |
| [dnsmasq](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsmasq)                       |     Dnsmasq DNS Forwarder     |
| [dnsmasq_dhcp](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsmasq_dhcp)             |         Dnsmasq DHCP          |
| [dns_query](https://github.com/netdata/go.d.plugin/tree/master/modules/dnsquery)                    |         DNS Query RTT         |
| [docker](https://github.com/netdata/go.d.plugin/tree/master/modules/docker)                         |         Docker Engine         |
| [docker_engine](https://github.com/netdata/go.d.plugin/tree/master/modules/docker_engine)           |         Docker Engine         |
| [dockerhub](https://github.com/netdata/go.d.plugin/tree/master/modules/dockerhub)                   |          Docker Hub           |
| [elasticsearch](https://github.com/netdata/go.d.plugin/tree/master/modules/elasticsearch)           |         Elasticsearch         |
| [energid](https://github.com/netdata/go.d.plugin/tree/master/modules/energid)                       |          Energi Core          |
| [example](https://github.com/netdata/go.d.plugin/tree/master/modules/example)                       |               -               |
| [filecheck](https://github.com/netdata/go.d.plugin/tree/master/modules/filecheck)                   |     Files and Directories     |
| [fluentd](https://github.com/netdata/go.d.plugin/tree/master/modules/fluentd)                       |            Fluentd            |
| [freeradius](https://github.com/netdata/go.d.plugin/tree/master/modules/freeradius)                 |          FreeRADIUS           |
| [haproxy](https://github.com/netdata/go.d.plugin/tree/master/modules/haproxy)                       |            HAProxy            |
| [hdfs](https://github.com/netdata/go.d.plugin/tree/master/modules/hdfs)                             |             HDFS              |
| [httpcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/httpcheck)                   |       Any HTTP Endpoint       |
| [isc_dhcpd](https://github.com/netdata/go.d.plugin/tree/master/modules/isc_dhcpd)                   |           ISC DHCP            |
| [k8s_kubelet](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubelet)               |            Kubelet            |
| [k8s_kubeproxy](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_kubeproxy)           |          Kube-proxy           |
| [k8s_state](https://github.com/netdata/go.d.plugin/tree/master/modules/k8s_state)                   |   Kubernetes cluster state    |
| [lighttpd](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd)                     |           Lighttpd            |
| [lighttpd2](https://github.com/netdata/go.d.plugin/tree/master/modules/lighttpd2)                   |           Lighttpd2           |
| [logind](https://github.com/netdata/go.d.plugin/tree/master/modules/logind)                         |        systemd-logind         |
| [logstash](https://github.com/netdata/go.d.plugin/tree/master/modules/logstash)                     |           Logstash            |
| [mongoDB](https://github.com/netdata/go.d.plugin/tree/master/modules/mongodb)                       |            MongoDB            |
| [mysql](https://github.com/netdata/go.d.plugin/tree/master/modules/mysql)                           |             MySQL             |
| [nginx](https://github.com/netdata/go.d.plugin/tree/master/modules/nginx)                           |             NGINX             |
| [nginxplus](https://github.com/netdata/go.d.plugin/tree/master/modules/nginxplus)                   |          NGINX Plus           |
| [nginxvts](https://github.com/netdata/go.d.plugin/tree/master/modules/nginxvts)                     |           NGINX VTS           |
| [ntpd](https://github.com/netdata/go.d.plugin/tree/master/modules/ntpd)                             |          NTP daemon           |
| [nvme](https://github.com/netdata/go.d.plugin/tree/master/modules/nvme)                             |         NVMe devices          |
| [openvpn](https://github.com/netdata/go.d.plugin/tree/master/modules/openvpn)                       |            OpenVPN            |
| [openvpn_status_log](https://github.com/netdata/go.d.plugin/tree/master/modules/openvpn_status_log) |            OpenVPN            |
| [pgbouncer](https://github.com/netdata/go.d.plugin/tree/master/modules/pgbouncer)                   |           PgBouncer           |
| [phpdaemon](https://github.com/netdata/go.d.plugin/tree/master/modules/phpdaemon)                   |           phpDaemon           |
| [phpfpm](https://github.com/netdata/go.d.plugin/tree/master/modules/phpfpm)                         |            PHP-FPM            |
| [pihole](https://github.com/netdata/go.d.plugin/tree/master/modules/pihole)                         |            Pi-hole            |
| [pika](https://github.com/netdata/go.d.plugin/tree/master/modules/pika)                             |             Pika              |
| [ping](https://github.com/netdata/go.d.plugin/tree/master/modules/ping)                             |       Any network host        |
| [prometheus](https://github.com/netdata/go.d.plugin/tree/master/modules/prometheus)                 |    Any Prometheus Endpoint    |
| [portcheck](https://github.com/netdata/go.d.plugin/tree/master/modules/portcheck)                   |       Any TCP Endpoint        |
| [postgres](https://github.com/netdata/go.d.plugin/tree/master/modules/postgres)                     |          PostgreSQL           |
| [powerdns](https://github.com/netdata/go.d.plugin/tree/master/modules/powerdns)                     | PowerDNS Authoritative Server |
| [powerdns_recursor](https://github.com/netdata/go.d.plugin/tree/master/modules/powerdns_recursor)   |       PowerDNS Recursor       |
| [proxysql](https://github.com/netdata/go.d.plugin/tree/master/modules/proxysql)                     |           ProxySQL            |
| [pulsar](https://github.com/netdata/go.d.plugin/tree/master/modules/portcheck)                      |         Apache Pulsar         |
| [rabbitmq](https://github.com/netdata/go.d.plugin/tree/master/modules/rabbitmq)                     |           RabbitMQ            |
| [redis](https://github.com/netdata/go.d.plugin/tree/master/modules/redis)                           |             Redis             |
| [scaleio](https://github.com/netdata/go.d.plugin/tree/master/modules/scaleio)                       |       Dell EMC ScaleIO        |
| [SNMP](https://github.com/netdata/go.d.plugin/blob/master/modules/snmp)                             |             SNMP              |
| [solr](https://github.com/netdata/go.d.plugin/tree/master/modules/solr)                             |             Solr              |
| [squidlog](https://github.com/netdata/go.d.plugin/tree/master/modules/squidlog)                     |             Squid             |
| [springboot2](https://github.com/netdata/go.d.plugin/tree/master/modules/springboot2)               |         Spring Boot2          |
| [supervisord](https://github.com/netdata/go.d.plugin/tree/master/modules/supervisord)               |          Supervisor           |
| [systemdunits](https://github.com/netdata/go.d.plugin/tree/master/modules/systemdunits)             |      Systemd unit state       |
| [tengine](https://github.com/netdata/go.d.plugin/tree/master/modules/tengine)                       |            Tengine            |
| [traefik](https://github.com/netdata/go.d.plugin/tree/master/modules/traefik)                       |            Traefik            |
| [unbound](https://github.com/netdata/go.d.plugin/tree/master/modules/unbound)                       |            Unbound            |
| [vcsa](https://github.com/netdata/go.d.plugin/tree/master/modules/vcsa)                             |   vCenter Server Appliance    |
| [vernemq](https://github.com/netdata/go.d.plugin/tree/master/modules/vernemq)                       |            VerneMQ            |
| [vsphere](https://github.com/netdata/go.d.plugin/tree/master/modules/vsphere)                       |     VMware vCenter Server     |
| [web_log](https://github.com/netdata/go.d.plugin/tree/master/modules/weblog)                        |         Apache/NGINX          |
| [wireguard](https://github.com/netdata/go.d.plugin/tree/master/modules/wireguard)                   |           WireGuard           |
| [whoisquery](https://github.com/netdata/go.d.plugin/tree/master/modules/whoisquery)                 |         Domain Expiry         |
| [windows](https://github.com/netdata/go.d.plugin/tree/master/modules/windows)                       |            Windows            |
| [x509check](https://github.com/netdata/go.d.plugin/tree/master/modules/x509check)                   |     Digital Certificates      |
| [zookeeper](https://github.com/netdata/go.d.plugin/tree/master/modules/zookeeper)                   |           ZooKeeper           |

## Configuration

Edit the `go.d.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d.conf
```

Configurations are written in [YAML](http://yaml.org/).

- [plugin configuration](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf)
- [specific module configuration](https://github.com/netdata/go.d.plugin/tree/master/config/go.d)

### Enable a collector

To enable a collector you should edit `go.d.conf` to uncomment the collector in question and change it from `no` to `yes`. 

For example, to enable the `example` plugin you would need to update `go.d.conf` from something like:

```yaml
modules:
#  example: no 
```

to

```yaml
modules:
  example: yes
```

Then [restart netdata](https://github.com/netdata/netdata/blob/master/docs/configure/start-stop-restart.md) for the change to take effect.

## Contributing

If you want to contribute to this project, we are humbled. Please take a look at
our [contributing guidelines](https://learn.netdata.cloud/contribute/handbook) and don't hesitate to contact us in our
forums. We have a whole [category](https://community.netdata.cloud/c/agent-development/9) just for this purpose!

### How to develop a collector

- Take a look at our [contributing guidelines](https://learn.netdata.cloud/contribute/handbook).
- [Fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo) this repository to your personal
  GitHub account.
- [Clone](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository#:~:text=to%20GitHub%20Desktop-,On%20GitHub%2C%20navigate%20to%20the%20main%20page%20of%20the%20repository,Desktop%20to%20complete%20the%20clone.)
  locally the **forked** repository (e.g `git clone https://github.com/odyslam/go.d.plugin`).
- Using a terminal, `cd` into the directory (e.g `cd go.d.plugin`)
- Add your module to the [modules](https://github.com/netdata/go.d.plugin/tree/master/modules) directory.
- Import the module in [`modules/init.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/init.go).
- To build it, run `make` from the plugin root dir. This will create a new `go.d.plugin` binary that includes your newly
  developed collector. It will be placed into the `bin` directory (e.g `go.d.plugin/bin`)
- Run it in the debug mode `bin/godplugin -d -m <MODULE_NAME>`. This will output the `STDOUT` of the collector, the same
  output that is sent to the Netdata Agent and is transformed into charts. You can read more about this collector API in
  our [documentation](https://learn.netdata.cloud/docs/agent/collectors/plugins.d#external-plugins-api).
- If you want to test the collector with the actual Netdata Agent, you need to replace the `go.d.plugin` binary that
  exists in the Netdata Agent installation directory with the one you just compiled. Once
  you [restart](https://learn.netdata.cloud/docs/configure/start-stop-restart) the Netdata Agent, it will detect and run
  it, creating all the charts. It is advised not to remove the default `go.d.plugin` binary, but simply rename it
  to `go.d.plugin.old` so that the Agent doesn't run it, but you can easily rename it back once you are done.
- Run `make clean` when you are done with testing.

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

Change `<module name>` to the module name you want to debug. See the [whole list](#available-modules) of available
modules.

## Netdata Community

This repository follows the Netdata Code of Conduct and is part of the Netdata Community.

- [Community Forums](https://community.netdata.cloud)
- [Netdata Code of Conduct](https://learn.netdata.cloud/contribute/code-of-conduct)
