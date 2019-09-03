# vcsa

This module will collect health statistics from one or more [`vCenter Server Appliance`](https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.vcsa.doc/GUID-223C2821-BD98-4C7A-936B-7DBE96291BA4.html) using [`Health API`](https://code.vmware.com/apis/60/vcenter-server-appliance-management).

Requirement:
  - `vSphere` 6.5+

It produces the following chart:
  - Health Status in `status`
    - appl_mgmt
    - database_storage
    - load
    - mem
    - software packages
    - storage
    - swap
    - system
    
Health statuses:

| Numeric | Text | Description |
| :---: | :---: | :--- |
| `-1`  | `unknown`  | Module failed to decode status.|
| `0`   | `green`  | The component is healthy.|
| `1`   | `yellow`  | The component is healthy, but may have some problems.|
| `2`   | `orange`  | The component is degraded, and may have serious problems.|
| `3`   | `red`  | The component is unavailable, or will stop functioning soon.|
| `4`   | `gray`  | No health data is available.|


### configuration

Needs only `url`, `username` and `password`.

Here is an example for 2 servers:

```yaml
jobs:
  - name         : vcsa1
    url          : https://203.0.113.0
    username     : admin@vsphere.local
    password     : somepassword

  - name         : vcsa2
    url          : https://203.0.113.10
    username     : admin@vsphere.local
    password     : somepassword
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/vcenter.conf).

---
