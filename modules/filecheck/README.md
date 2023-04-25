<!--
title: "Files and directories monitoring with Netdata"
description: "Monitor the health and performance of files and directories with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/filecheck/README.md"
sidebar_label: "Files and Dirs"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/System metrics"
-->

# Files and directories collector

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

> **Warning**: For security reasons, this should not be applied recursively, but only to the exact set of directories
> that lead to the file/dir you want to monitor.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/filecheck/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/filecheck.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

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
  ./go.d.plugin -d -m filecheck
  ```
