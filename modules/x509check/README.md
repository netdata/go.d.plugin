# x509 certificate expiry check

Checks the time until a x509 certificate expires.

It produces the following charts:

1. **Time Until Certificate Expiration** in seconds
 * expiry

### configuration

For all available options and defaults please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/x509.conf).
___

```yaml
update_every : 60

jobs:
  - name: example_org
    source: https://example.org:443
  
  - name: my_site_org
    source: https://my_site_org:443
```
