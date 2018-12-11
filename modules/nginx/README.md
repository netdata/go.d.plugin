# nginx

This module will monitor one or more nginx servers depending on configuration. Servers can be either local or remote.

**Requirements:**
 * nginx with configured 'ngx_http_stub_status_module'
 * 'location /stub_status'

Example nginx configuration can be found in 'go.d/nginx.conf'

It produces following charts:

1. **Active Connections**
 * active

2. **Requests** in requests/s
 * requests

3. **Active Connections by Status**
 * reading
 * writing
 * waiting

4. **Connections Rate** in connections/s
 * accepts
 * handled

### configuration

Needs only `url` to server's `stub_status`

Here is an example for local server:

```yaml
jobs:
  - name: local
    url : http://localhost/stub_status
      
  - name: remote
    url : http://100.64.0.1/stub_status
```

Without configuration, module attempts to connect to `http://localhost/stub_status`

---
