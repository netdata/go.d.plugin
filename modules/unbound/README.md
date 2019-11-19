# unbound

This module monitors one or more [`Unbound`](https://nlnetlabs.nl/projects/unbound/about/) servers, depending on your configuration.

#### Requirements

-   `Unbound` with enabled `remote-control` interface (see [unbound.conf](https://nlnetlabs.nl/documentation/unbound/unbound.conf))

If using unix socket:

-  socket should be readable and writeable by `netdata` user

If using ip socket and TLS is disabled:

-  socket should be accessible via network

If TLS is enabled, in addition:

-  `control-key-file` should be readable by `netdata` user
-  `control-cert-file` should be readable by `netdata` user

For auto detection parameters from `unbound.conf`:

-  `unbound.conf` should be readable by `netdata` user  


### Configuration

Needs only `address` to server's `remote-control` interface if TLS is disabled or `address` is unix socket.
Otherwise you need to set path to the `control-key-file` and `control-cert-file` files.

Module tries to auto detect following parameters reading `unbound.conf`:
  - address (`control-interface` and `control-port`)
  - cumulative (`statistics-cumulative`)
  - use_tls (`control-use-cert`)
  - tls_cert (`control-cert-file`)
  - tls_key (`control-key-file`)

Module supports both cumulative and non cumulative modes. Default is non cumulative. If your server has enabled 
`statistics-cumulative` but module fails to auto detect it (`unbound.conf` is not readable or it is a remote server) 
you need to set it manually in the configuration file. 

Here is an example for several servers:

```yaml
jobs:
  - name: local
    address: 127.0.0.1:8953
    use_tls: yes
    tls_skip_verify: yes
    tls_cert: /etc/unbound/unbound_control.pem
    tls_key: /etc/unbound/unbound_control.key

  - name: remote
    address: 203.0.113.10:8953
    use_tls: no

  - name: remote_cumulative
    address: 203.0.113.11:8953
    use_tls: no
    cumulative: yes
      
  - name: socket
    address: /var/run/unbound.sock
```
 
For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/unbound.conf).


### Troubleshooting 
Ensure that the control protocol is actually configured correctly.
Run following command as `root` user:
> unbound-control stats_noreset

It should print out a bunch of info about the internal statistics of the server.
If this returns an error, you don't have the control protocol set up correctly.

Check the module debug output.
Run following command as `netdata` user:

> ./go.d.plugin -d -m unbound

---
