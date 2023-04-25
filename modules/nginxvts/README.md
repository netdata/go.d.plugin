<!--
title: "NGINX VTS monitoring"
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/nginxvts/README.md"
sidebar_label: "NGINX VTS"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# NGINX VTS collector

`nginxvts` can monitor statistics of NGINX which configured
with [`nginx-module-vts`](https://github.com/vozlt/nginx-module-vts).

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/nginxvts/metrics.csv) for a list of
metrics.

Refer [`nginx-module-vts`](https://github.com/vozlt/nginx-module-vts#json) for more information.

## Configuration

Edit the `go.d/nginxvts.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/nginxvts.conf
```

Needs only `url` to server's `stub_status`. Here is an example for local and remote servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/status/format/json
  - name: remote
    url: http://203.0.113.0/status/format/json
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginxvts.conf).

## Troubleshooting

To troubleshoot issues with the `nginxvts` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m nginxvts
  ```
