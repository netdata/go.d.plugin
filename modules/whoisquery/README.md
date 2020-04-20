# Whois domain expiry monitoring with Netdata

This collector module checks the remaining time until a domain is expired.

## Charts
This collector produces the following chart:
-   Time until domain expiry in `seconds`

## Configuration
Edit the `go.d/whoisquery.conf` configuration file using `edit-config` from the your agent's [config
directory](/docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/whoisquery.conf
```
Needs only `source`.

Use `days_until_expiration_warning` and `days_until_expiration_critical` for each job to indicate the expiry warning and critical days. The default values are 90 for warning, and 30 days for critical.

Here is an example:

```yaml
update_every : 60

jobs:
  - name   : my_site
    source : my_site.com
    
  - name   : my_another_site
    source : my_another_site.com
    days_until_expiration_critical: 20

```

For all available options and defaults please, see the module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/whoisquery.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m whoisquery
