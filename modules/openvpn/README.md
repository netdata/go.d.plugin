# openvpn

This module will monitor one or more OpenVPN instances via Management Interface.

**Requirements:**
 * OpenVPN with enabled [`management-interface`](https://openvpn.net/community-resources/management-interface/)


It produces the following charts:

1. **Total Number Of Active Clients** in clients
 * clients

2. **Total Traffic** in kilobits/s
 * in
 * out
 
Per user charts (disabled by default, see `per_user_stats` in the module config file):

1. **User Traffic** in kilobits/s
 * received
 * sent

2. **User Connection Time** in seconds
 * time
 
 
### configuration

`openvpn` collector is disabled by default. Should be explicitly enabled in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

Reason:
 >  Currently,the OpenVPN daemon can at most support a single management client any one time.

So to not break other tools that uses Management Interface we decided to disable it by default.

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/openvpn.conf).
___

Needs only `address` of OpenVPN Management Interface.

Here is an example for 2 OpenVPN instances:

```yaml
jobs:
  - name: local
    address : /dev/openvpn
      
  - name: remote
    address : 100.64.0.1:7505
```

Without configuration, module attempts to connect to `127.0.0.1:7505`

---
