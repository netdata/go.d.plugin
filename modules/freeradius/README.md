# freeradius

It produces following charts:

1. **Authentication** in pps
 * requests
 * responses
 
2. **Authentication Responses** in pps
 * accepts
 * rejects
 * challenges
 
 3. **Bad Authentication Requests** in pps
 * dropped
 * duplicate
 * invalid
 * malformed
 * unknown-types
 
 4. **Proxy Authentication** in pps
  * requests
  * responses
  
 5. **Proxy Authentication Responses** in pps
  * accepts
  * rejects
  * challenges
  
 6. **Proxy Bad Authentication Requests** in pps
  * dropped
  * duplicate
  * invalid
  * malformed
  * unknown-types

7. **Accounting** in pps
 * requests
 * responses

8. **Bad Accounting Requests** in pps 
  * dropped
  * duplicate
  * invalid
  * malformed
  * unknown-types

9. **Proxy Accounting** in pps
 * requests
 * responses

10. **Proxy Bad Accounting Requests** in pps 
  * dropped
  * duplicate
  * invalid
  * malformed
  * unknown-types


### configuration

Module specific options:
 * `host`    - server address. Default is 127.0.0.1.
 * `port`    - server port. Default is 18121.
 * `secret`  - secret. Default is `adminsecret`.
 * `timeout` - request timeout. Default is 1 seconds.
 
Without configuration, module will try to connect to 127.0.0.1:18121 and will use `adminsecret` as a secret.
 
Configuration sample:

```yaml
jobs:
  - name: local
    host : 127.0.0.1
```

### prerequisite

Module query FreeRADIUS using `StatusServer` packet. FreeRADIUS status feature is disabled by default.
Should be enabled.

The configuration for the status server is automatically created in the sites-available directory.
By default, server is enabled and can be queried from every client.
FreeRADIUS will only respond to status-server messages if the status-server virtual server has been enabled.

To enable FreeRADIUS status do the following:
 * cd sites-enabled
 * ln -s ../sites-available/status status
 * restart FreeRADIUS server
---
