<!--
title: "Files and directories monitoring with Netdata"
description: "Monitor the health and performance of files and directories with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/fileheck/README.md
sidebar_label: "Files and Dirs"
-->

# filecheck

This module monitors files and directories.

File metrics:

- existence
- time since the last modification
- size

Directory metrics:

- existence
- time since the last modification
- number of files
- size

## Permissions

`netdata` user needs the following permissions on all the directories in pathname that lead to the file/dir:

- files monitoring: `execute`.
- directories monitoring: `read` and `execute`.

If you need to modify the permissions we
suggest [to use file access control lists](https://linux.die.net/man/1/setfacl):

```cmd
setfacl -m u:netdata:rx file ...
``` 

> :warning: For security reasons, this should not be applied recursively, but only to the exact set of directories that lead to the file/dir you want to monitor.

## Charts

Files and directories have their own set of charts.

### Files

- File Existence in `boolean`
- File Time Since the Last Modification in `seconds`
- File Size in `bytes`

### Directories

- Dir Existence in `boolean`
- Dir Time Since the Last Modification in `seconds`
- Dir Number of Files in `files`
- Dir Size in `bytes`

## Configuration

Edit the `go.d/filecheck.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/filecheck.conf
```

Needs only a path to a file or a directory. The path supports `*` wildcard.

Here is an example:

```yaml
jobs:
  - name: files_dirs_example
    discovery_every: 30s
    files:
      include:
        - '/path/to/file1'
        - '/path/to/file2'
        - '/path/to/*.log'
    dirs:
      collect_dir_size: no
      include:
        - '/path/to/dir1'
        - '/path/to/dir2'
        - '/path/to/dir3*'

  - name: files_example
    discovery_every: 30s
    files:
      include:
        - '/path/to/file1'
        - '/path/to/file2'
        - '/path/to/*.log'

  - name: dirs_example
    discovery_every: 30s
    dirs:
      collect_dir_size: yes
      include:
        - '/path/to/dir1'
        - '/path/to/dir2'
        - '/path/to/dir3*'
```

For all available options, see the Filecheck
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/filecheck.conf).

## Limitations

- filecheck uses `stat` call to collect metrics, which is not very efficient.

## Troubleshooting

To troubleshoot issues with the `filecheck` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m filecheck
```
