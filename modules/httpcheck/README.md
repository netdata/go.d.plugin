# Any HTTP endpoints monitoring with Netdata

This module monitors one or more http servers availability and response time.

## Charts

It produces the following charts:

-   HTTP Response Time in `ms`
-   HTTP Check Status in `boolean`
-   HTTP Current State Duration in `seconds`
-   HTTP Response Body Length in `characters`

## Check statuses

| Status        | Description|
| ------------- |-------------|
| success      |No error on HTTP request, body reading and body content checking |
| timeout      |Timeout error on HTTP request|
| bad content |The body of the response didn't match the regex (only if `response_match` option is set)|
| bad status |Response status code not in `status_accepted`|
| no connection |Any other network error not specifically handled by the module|


## Configuration

Edit the `go.d/httpcheck.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/httpcheck.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: cool_website1
    url: http://cool.website1:8080/home
      
  - name: cool_website2
    url: http://cool.website2:8080/home
    status_accepted:
      - 200
      - 201
      - 202
    response_match: <title>My cool website!<\/title>
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/httpcheck.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m httpcheck

