# vCenter Server Appliance monitoring with Netdata

The [`vCenter Server Appliance`](https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.vcsa.doc/GUID-223C2821-BD98-4C7A-936B-7DBE96291BA4.html) using [`Health API`](https://code.vmware.com/apis/60/vcenter-server-appliance-management) is a preconfigured Linux virtual machine, which is optimized for running VMware vCenter Server® and the associated services on Linux.

This module collects health statistics from one or more `vCenter Server Appliance` servers, depending on your configuration.

## Requirements

-  `vSphere` 6.5+

## Charts

-   Overall System Health in `status`
-   Components Health in `status`
-   Software Updates Health in `status`

## Health statuses
    
Overall System Health:

| Numeric | Text | Description |
| :---: | :---: | :--- |
| `-1`  | `unknown`  | Module failed to decode status.|
| `0`   | `green`  | All components in the appliance are healthy.|
| `1`   | `yellow`  | One or more components in the appliance might become overloaded soon.|
| `2`   | `orange`  | One or more components in the appliance might be degraded.|
| `3`   | `red`  | One or more components in the appliance might be in an unusable status and the appliance might become unresponsive soon.|
| `4`   | `gray`  | No health data is available.|

Components Health:

| Numeric | Text | Description |
| :---: | :---: | :--- |
| `-1`  | `unknown`  | Module failed to decode status.|
| `0`   | `green`  | The component is healthy.|
| `1`   | `yellow`  | The component is healthy, but may have some problems.|
| `2`   | `orange`  | The component is degraded, and may have serious problems.|
| `3`   | `red`  | The component is unavailable, or will stop functioning soon.|
| `4`   | `gray`  | No health data is available.|

Software Updates Health:

| Numeric | Text | Description |
| :---: | :---: | :--- |
| `-1`  | `unknown`  | Module failed to decode status.|
| `0`   | `green`  | No updates available.|
| `2`   | `orange`  | Non-security patches might be available.|
| `3`   | `red`  | Security patches might be available.|
| `4`   | `gray`  | An error retrieving information on software updates.|


## Configuration

Edit the `go.d/vsca.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/vsca.conf
```

Needs only `url`, `username` and `password`. Here is an example for 2 servers:

```yaml
jobs:
  - name: vcsa1
    url: https://203.0.113.0
    username: admin@vsphere.local
    password: somepassword

  - name: vcsa2
    url: https://203.0.113.10
    username: admin@vsphere.local
    password: somepassword
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/vcenter.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m vcsa
