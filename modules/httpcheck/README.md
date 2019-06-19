# httpcheck

This module will monitor one or more http servers for availability and response time.

It produces the following charts:

1. **HTTP Response Time** in ms
 * ms

2. **HTTP Check Status** in boolean
 * success
 * no connection
 * timeout
 * bad content
 * bad status

3. **HTTP Response Body Length** in characters
 * length

### check statuses

| Status        | Description|
| ------------- |-------------|
| success      |No error on HTTP request, body reading and body content checking |
| timeout      |Timeout error on HTTP request|
| bad content |The body of the response didn't match the regex (only if `response_match` option is set)|
| bad status |Response status code not in `accepted_statuses`|
| no connection |Any other network error not specifically handled by the module|


### configuration
 
Here is an example for 2 servers:

```yaml
jobs:
  - name: cool_website1
    url: http://cool.website1:8080/home
      
  - name: cool_website2
    url: http://cool.website2:8080/home
    status_accepted: [200, 201, 202]
    response_match: <title>My cool website!<\/title>
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/httpcheck.conf).

---
