# Docker Engine monitoring with Netdata

[`Docker Engine`](https://docs.docker.com/engine/) is the industry’s de facto container runtime that runs on various Linux (CentOS, Debian, Fedora, Oracle Linux, RHEL, SUSE, and Ubuntu) and Windows Server operating systems.

This module will monitor one or more `Docker Engines` applications, depending on your configuration.

## Requirements

-   Docker with enabled [`metric-address`](https://docs.docker.com/config/thirdparty/prometheus/)

## Charts

It produces the following charts:

-   Container Actions in `actions/s`
-   Container States in `containers`
-   Builder Builds Fails By Reason in `fails/s`
-   Health Checks in `events/s`

If Docker is running in in [Swarm mode](https://docs.docker.com/engine/swarm/) and the instance is a Swarm manager:

-   Swarm Manager Leader in `bool`
-   Swarm Manager Object Store in `count`
-   Swarm Manager Nodes Per State in `count`
-   Swarm Manager Tasks Per State in `count`

## Configuration

Edit the `go.d/docker_engine.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/docker_engine.conf
```

Needs only `url` to docker `metric-address`. Here is an example for 2 docker instances:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:9323/metrics
      
  - name: remote
    url : http://203.0.113.10:9323/metrics
```

For all available options, please see the module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/docker_engine.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m docker_engine
