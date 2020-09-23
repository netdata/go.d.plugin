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

## Configuration

Edit the `go.d/filecheck.conf` configuration file using `edit-config` from the Agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/filecheck.conf
```

Needs only a path to a file or a directory. **The path doesn't support any wildcards**.

Here is an example:

```yaml
jobs:
 - name: files_dirs_example
   files:
     include:
       - '/path/to/file1'
       - '/path/to/file2'
   dirs:
     include:
       - '/path/to/dir1'
       - '/path/to/dir2'

 - name: files_example
   files:
     include:
       - '/path/to/file1'
       - '/path/to/file2'

 - name: dirs_example
   dirs:
     include:
       - '/path/to/dir1'
       - '/path/to/dir2'
```

For all available options, see the Filecheck collector's [configuration
file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/filecheck.conf).

## Limitations

-   file/dir path pattern doesn't support any wildcards
-   filecheck uses `stat` call to collect metrics, which is not very efficient.

## Troubleshooting

To troubleshoot issues with the Filecheck collector, run the `go.d.plugin` with the debug option enabled.
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
./go.d.plugin -d -m filecheck
```
