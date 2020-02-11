# Windows Machines monitoring with Netdata

This module will monitor one or more Windows machines, using the [wmi_exporter](https://github.com/martinlindhe/wmi_exporter).

Module collects metrics from following collectors:

-   cpu
-   memory
-   net
-   logical_disk
-   os
-   system
-   logon

Run `wmi_exporter` with these collectors:     
    
 > wmi-exporter-0.9.0-386.exe --collectors.enabled="cpu,memory,net,logical_disk,os,system,logon"
 

Installation: please follow [official guide](https://github.com/martinlindhe/wmi_exporter#installation).
 
## Charts

#### cpu 

-   Total CPU Utilization (all cores) in `percentage`
-   Received and Serviced Deferred Procedure Calls (DPC) in `dpc/c`
-   Received and Serviced Hardware Interrupts in `interrupts/s`
-   CPU Utilization Per Core in `percentage`
-   Time Spent in Low-Power Idle State Per Core in `percentage`

#### memory
 
-   Memory Utilization in `KiB`
-   Memory Page Faults in `events/s`
-   Swap Utilization in `KiB`
-   Swap Operations in `operations/s`
-   Swap Pages in `pages/s`
-   Cached Data in `KiB`
-   Cache Faults in `events/s`
-   System Memory Pool in `KiB`

#### network
 
-   Bandwidth Per NIC in `kilobits/s`
-   Packets Per NIC in `packets/s`
-   Errors Per NIC in `errors/s`
-   Discards Per NIC in `discards/s`

#### disk
 
-   Utilization Per Disk in `KiB`
-   Bandwidth Per Disk in `KiB/s`
-   Operations Per Disk in `operations/s`
-   Average Read/Write Latency Disk in `milliseconds`
  
#### system
 
-   Processes in `number`
-   Threads in `number`
-   Uptime in `seconds`

#### logon
 
-   Active User Logon Sessions By Type in `sessions`
  
## Configuration

Edit the `go.d/wmi.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/wmi.conf
```

Needs only `url` to `wmi_exporter` metrics endpoint. Here is an example for 2 instances:

```yaml
jobs:
  - name : win_server1
    url  : http://203.0.113.10:9182/metrics

  - name : win_server2
    url  : http://203.0.113.11:9182/metrics
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/wmi.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m wmi
