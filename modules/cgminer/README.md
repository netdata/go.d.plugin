# cgminer collector

A data collection module for cgminer.

## Charts

This module collects metrics from cgminer and creates charts to visualize the data. The metrics collected include information on overall performance, individual GPU performance, and more.

## Configuration

Edit the `go.d/cgminer.conf` configuration file using `edit-config` from the Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/cgminer.conf
```

Disabled by default. Should be explicitly enabled in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

```yaml
# go.d.conf
modules:
  cgminer: yes
```

Disabled by default. Should be explicitly enabled in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

```yaml
# go.d.conf
modules:
  cgminer: yes
```
For all available options, see the cgminer collector's configuration file.

## Troubleshooting

To troubleshoot issues with the cgminer collector, run the go.d.plugin with the debug option enabled. The output should give you clues as to why the collector isn't working.

- Navigate to the plugins.d directory, usually at /usr/libexec/netdata/plugins.d/. If that's not the case on your system, open netdata.conf and look for the plugins setting under [directories].
```
cd /usr/libexec/netdata/plugins.d/
```
- Switch to the netdata user.
```
sudo -u netdata -s
```
- Run the go.d.plugin to debug the collector:
```
./go.d.plugin -d -m cgminer
```
