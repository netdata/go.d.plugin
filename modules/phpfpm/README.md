# phpfpm

This module will monitor one or more php-fpm instances depending on configuration.

**Requirements:**

-   php-fpm with enabled `status` page
-   access to `status` page via web server

It produces following charts:

1.  **Active Connections**

    -   active
    -   maxActive
    -   idle

2.  **Requests** in requests/s

    -   requests

3.  **Performance**

    -   reached
    -   slow

## configuration

Needs only `url` to server's `status`

Here is an example for local server:

```yaml
jobs:
  - name: local
    url: http://localhost/status?full&json

  - name: local
    url: http://127.0.0.1/status?full&json

  - name: local
    url: http://[::1]/status?full&json
```

Without configuration, module attempts to connect to `http://localhost/status?full&json`

---