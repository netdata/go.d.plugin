# ScaleIO

This module will monitor one or more [ScaleIO (VxFlex OS)](https://www.dellemc.com/en-us/storage/data-storage/software-defined-storage.htm) instances via VxFlex OS Gateway API.

## Notes
Module was tested on:
 - VxFlex OS REST API v2.5
 - VxFlex OS v2.6.1.1_113
 - VxFlex OS 3.0.0.1_134, REST API v3.0

## Charts

Please see [CHARTS.md](https://github.com/netdata/go.d.plugin/blob/master/modules/scaleio/CHARTS.md) 
 
## Configuration

Needs only `url` of VxFlex OS Gateway API, MDM `username` and `password`.

Here is an example for 2 instances:

```yaml
jobs:
  - name     : local
    url      : https://127.0.0.1
    username : admin
    password : password
      
  - name     : remote
    url      : https://100.64.0.1
    username : admin
    password : password
```
For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/scaleio.conf).

Without configuration module won't work.

---
