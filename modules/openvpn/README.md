# openvpn

This module will monitor one or more OpenVPN instances via Management Interface.

**Requirements:**
 * OpenVPN with enabled [`management-interface`](https://openvpn.net/community-resources/management-interface/)


It produces the following charts:

1. **Total Number Of Active Clients** in clients
 * clients

2. **Total Traffic** in KiB/s
 * in
 * out
 
 
### configuration

`openvpn` collector is disabled by default. Should be explicitly enabled in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

Reason:
 >  Currently,the OpenVPN daemon can at most support a single management client any one time.

So to not break other tools that uses Management Interface we decided to disable it by default. 
___

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/openvpn.conf).
___

Needs only `address` of OpenVPN Management Interface.

Here is an example for 2 OpenVPN instances:

```yaml
jobs:
  - name: local
    url : /dev/openvpn
      
  - name: remote
    url : http://100.64.0.1:7505
```

Without configuration, module attempts to connect to `127.0.0.1:7505`

---
