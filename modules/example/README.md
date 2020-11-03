<!--
title: "Example collector"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/example/README.md
sidebar_label: "Example collector"
-->

# Example collector

## Charts

## Configuration

Edit the `go.d/example.conf` configuration file using `edit-config` from the Agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/example.conf
```

To add a new endpoint to collect metrics from, or change the URL that Netdata looks for, add or configure the `name` and
`url` values. Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example configuration with several jobs:

```yaml
jobs:
  - name: example_1
    num_of_charts: 1
    num_of_dimensions: 3

  - name: example_2
    num_of_charts: 2
    num_of_dimensions: 6
```

---

For all available options, see the Example collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/example.conf).


## Troubleshooting

To troubleshoot issues with the Example collector, run the `go.d.plugin` orchestrator with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugins directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m example
```
