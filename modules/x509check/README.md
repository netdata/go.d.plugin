# x509 certificates monitoring with Netdata

This module checks the time until a x509 certificate expiration.

## Charts

It produces only one chart:

-   Time Until Certificate Expiration in `seconds`
 
## Configuration

Edit the `go.d/x509check.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/x509check.conf
```

Needs only `source`.

Use `smtp` scheme for smtp servers, `file` for files and `https` or `tcp` for others. Port is mandatory for all non-file schemes.

Here is an example for 3 sources:

```yaml
update_every : 60

jobs:
  - name   : my_site_cert
    source : https://my_site.org:443
    
  - name   : my_file_cert
    source : file:///home/me/cert.pem

  - name   : my_smtp_cert
    source : smtp://smtp.my_mail.org:587
```

For all available options and defaults please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/x509check.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m x509check
