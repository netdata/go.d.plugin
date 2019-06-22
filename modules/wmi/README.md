# wmi

This module will monitor one or more Windows machines, using the [wmi_exporter](https://github.com/martinlindhe/wmi_exporter).

**wmi_exporter configuration**

Module collects metrics from following collectors:
   - cpu
   - memory
   - net
   - logical_disk
   - os
   - system

Run `wmi_exporter` with these collectors:     
    
 > wmi-exporter-0.7.0-386.exe --collectors.enabled="cpu,memory,net,logical_disk,os,system"
 

Installation: please follow [official guide](https://github.com/martinlindhe/wmi_exporter#installation).
 
### charts

#### cpu 

1. **Total CPU Utilization (all cores)** in percentage
  * dpc
  * user
  * privileged
  * interrupt

2. **Received and Serviced Deferred Procedure Calls (DPC)** in dpc/c

3. **Received and Serviced Hardware Interrupts** in interrupts/s

4. **CPU Utilization** Per Core in percentage
  * dpc
  * user
  * privileged
  * interrupt

5. **Time Spent in Low-Power Idle State** Per Core in percentage
  * c1
  * c2
  * c3

#### memory
 
1. **Memory Utilization** in KiB
  * available
  * used

2. **Memory Page Faults** in events/s
  * page faults

3. **Swap Utilization** in KiB
  * available
  * used

4. **Swap Operations** in operations/s
  * read
  * write

5. **Swap Pages** pages/s
  * read
  * written

6. **Cached Data** in KiB
  * cached

7. **Cache Faults** in events/s
  * cache faults

8. **System Memory Pool** in KiB
  * paged
  * non-pages

#### network
 
1. **Bandwidth** Per NIC in kilobits/s
  * received
  * sent

2. **Packets** Per NIC in packets/s
  * received
  * sent

3. **Errors** Per NIC in errors/s
  * inbound
  * outbound

4. **Discards** Per NIC in discards/s
  * inbound
  * outbound

#### disk
 
1. **Utilization** Per Disk in KiB
  * free
  * used

2. **Bandwidth** Per Disk in KiB/s
  * read
  * write

3. **Operations** Per Disk in operations/s
  * reads
  * writes
  
#### system
 
1. **Processes** in number
  * processes

2. **Threads** in number
  * threads

3. **Uptime** in seconds
  * time

 
 
### configuration

Needs only `url` to `wmi_exporter` metrics endpoint.

Here is an example for 2 instances:

```yaml
jobs:
  - name : win_server1
    url  : http://10.0.0.1:9182/metrics

  - name : win_server2
    url  : http://10.0.0.2:9182/metrics
```
For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/wmi.conf).

---
