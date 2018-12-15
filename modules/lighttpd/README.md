# lighttpd

This module will monitor one or more Lighttpd servers depending on configuration.

**Requirements:**
 * lighttpd with enabled `mod_status`

It produces the following charts:

1. **Requests** in requests/s
 * requests

2. **Bandwidth** in kilobytes/s
 * sent

3. **Servers** in servers
 * idle
 * busy

4. **Scoreboard** events
 * waiting
 * open
 * close
 * hard error
 * keepalive
 * read
 * read post
 * write
 * handle request
 * requests start
 * requests end
 * response start
 * requests end

5. **Uptime** in seconds
 * uptime


### configuration

Needs only `url` to server's `server-status?auto`

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://localhost/server-status?auto
      
  - name: remote
    url : http://100.64.0.1/server-status?auto
```

Without configuration, module attempts to connect to `http://localhost/server-status?auto`

---
