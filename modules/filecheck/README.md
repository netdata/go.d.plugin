<!--
title: "Files and directories monitoring with Netdata"
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/fileheck/README.md
sidebar_label: "Files and Dirs"
-->

# filecheck

This module monitors files and directories.

File metrics:

-   existence
-   time since the last modification
-   size

Directory metrics:

-   existence
-   time since the last modification
-   number of files

## Charts

Files and directories have their own set of charts.

### Files

-   File Existence in `boolean`
-   File Time Since the Last Modification in `seconds`
-   File Size in `bytes`

### Directories

-   Dir Existence in `boolean`
-   Dir Time Since the Last Modification in `seconds`
-   Dir Number of Files in `files`

## Troubleshooting

To troubleshoot issues with the Portcheck collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m portcheck
```
