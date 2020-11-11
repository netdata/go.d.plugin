<!--
title: "File/Folder size monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/du/README.md
sidebar_label: "Du"
-->

# File/Folder size monitoring with Netdata

This module provides file/folder size monitoring like Linux `du` command. 

## Charts

It produces only one chart:

-   File or Folder size in `bytes`

## Configuration

Edit the `go.d/du.conf` configuration file using `edit-config` from the your agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/du.conf
```

Here is an example:

```yaml
jobs:
  - name: job1
    paths:
      - /var/log/nginx
      - /tmp/tmp.log
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/du.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m du
