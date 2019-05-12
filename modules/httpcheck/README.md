# httpcheck

This module will monitor one or more http servers for availability and response time.

It produces the following charts:

1. **HTTP Response Time** in ms
 * ms

2. **HTTP Check Status** in boolean
 * success
 * no connection
 * timeout
 * dns lookup error
 * address parse error
 * redirect error
 * body read error
 * bad content
 * bad status

3. **HTTP Response Body Length** in characters
 * length

### check statuses

| Status        | Description           |
| ------------- |:-------------:|
| success      | no error on HTTP request, body reading and its content checking |
| timeout      | timeout error on HTTP request |
| dns lookup error | dns lookup error on HTTP request|
| address parse error | address parse error on HTTP request |
| redirect error | redirect is disabled, but server returned 3xx|
| body read error | `response_match` option is set. Error during response body reading|
| bad content | `response_match` option is set. The body of the response didn't match the regex|
| bad status | response status code not in `accepted_statuses`|


### configuration
 
Here is an example for 2 servers:

```yaml
jobs:
  - name: cool_website1
    url: http://cool.website1:8080/home
      
  - name: cool_website2
    url: http://cool.website2:8080/home
    accepted_statuses: [200, 201, 202]
    response_match: <title>My cool website!<\/title>
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/httpcheck.conf).

---
