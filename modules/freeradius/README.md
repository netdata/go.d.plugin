# FreeRADIUS monitoring with Netdata

[`FreeRADIUS`](https://freeradius.org/) is a modular, high performance free RADIUS suite.

This module will monitor one or more `FreeRADIUS` servers, depending on your configuration.

## Requirements

-   `FreeRADIUS` with enabled status feature.

The configuration for the status server is automatically created in the sites-available directory.
By default, server is enabled and can be queried from every client.

To enable status feature do the following:

-   `cd sites-enabled`
-   `ln -s ../sites-available/status status`
-   restart FreeRADIUS server


## Charts

It produces following charts:

-   Authentication in `pps`
-   Authentication Responses in `pps`
-   Bad Authentication Requests in `pps`
-   Proxy Authentication in `pps`
-   Proxy Authentication Responses in `pps`
-   Proxy Bad Authentication Requests in `pps`
-   Accounting in `pps`
-   Bad Accounting Requests in `pps` 
-   Proxy Accounting in `pps`
-   Proxy Bad Accounting Requests in `pps` 

## Configuration
 
Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    host: 127.0.0.1

  - name: remote
    host: 203.0.113.10
    secret: secret 
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/freeradius.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m freeradius
