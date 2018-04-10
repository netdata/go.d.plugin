# go.d.plugin (WIP)

[Netdata](https://github.com/firehol/netdata) `python.d.plugin` rewritten in `go`

Plugin configuration directory: `/etc/netdata/`

Modules configuration directory: `/etc/netdata/go.d/`

Configurations are written in [TOML](https://github.com/toml-lang/toml/blob/master/README.md).

#### Plugin configuration file:

```toml
# enable/disable the whole go.d.plugin (all its modules)
enabled = "yes"

# default_run determines whether the module is enabled by default
default_run = "yes"

# sets the maximum number of CPUs that can be executing
# if not setted go.d.plugin will use all logical CPUs (default)
#max_procs = 1

# enable/disable specific module
[modules]
  httpcheck = "yes"
  portcheck = "yes"
 
```

#### Module configuration file:

All jobs support a set of predefined parameters.
Each job may define its own, overriding the defaults.
All of them are optional.

```toml
  update_every = 1         # data collection frequency in seconds
  retries = 60             # number of restoration attempts
  autodetection_retry = 0  # re-check interval in seconds
  chart_cleanup = 10       # chart cleanup interval in iterations
```

The default jobs share the same *name*. Jobs with the same name
are mutually exclusive. Only one of them will be allowed running at
any time. This allows autodetection to try several alternatives and
pick the one that works.

Every configuration file must have one of these formats:

 - Configuration for job that uses only base parameters:

```toml
[global]
  update_every = 1
  retries = 60
  autodetection_retry = 0
  chart_cleanup = 10
```

- Configuration for job that uses specific parameters or for many jobs:

```toml
[global] # default for all jobs
  update_every = 1
  retries = 60
  autodetection_retry = 0
  chart_cleanup = 10
  
[job1.base] # overrides global per job
  update_every = 2
  name = "job1"
[job1.specific]
  var1 = 33
  var2 = "var2"
  var3 = true
  
[job2.base] # overrides global per job
  update_every = 3
  name = "job1"
[job2.specific]
  var1 = 44
  var2 = "var2"
  var3 = false
```

---

The following modules are supported:

# httpcheck
Monitors remote http server for availability and response time.

Module *specific* variables and their default values:
```toml
  status_accepted  = [200]   # Optional. List of accepted statuses
  response_match   = ""      # Optional. Regex match in body of response (ex.: "REGULAR_EXPRESSION")
  url              = ""      # Required. URL. 
  body             = ""      # Optional. HTTP request body (ex.: '''{'some':'data'}''')
  header           = {""=""} # Optional. HTTP request headers (ex.: {"X-API-Key" = "key"}
  method           = "GET"   # Optional. HTTP request method
  username         = ""      # Optional. Username for basic auth
  password         = ""      # Optional. Password for basic auth
  proxy_username   = ""      # Optional. Proxy username for proxy basic auth
  proxy_password   = ""      # Optional. Proxy password for proxy basic auth
  proxy_url        = ""      # Optional. Proxy URL (default is to use a proxy URL from enviroment)
  tls_verify       = false   # Optional. Controls whether a client verifies the server's certificate chain and host name
  follow_redirects = false   # Optional. Whether to follow redirects from the server
  timeout          = ""      # Optional. A time limit for requests (default = UpdateEvery)
```

#### sample configuration
```toml
[global]
  update_every = 5
  
  [job1.base]
    name = "server1"
  [job1.specific]
    url = "http://1.2.3.4:80/_check"
    status_accepted = [401, 402, 403]
    timeout = "750ms"

  [job2.base]
    name = "server2"
  [job2.specific]
    url = "http://4.3.2.1:80/_check"
    status_accepted = [200, 300, 400]
    timeout = "1s"

```

---

# portcheck
Monitors a remote TCP service.

Module  *specific* variables and their default values:
```toml
  host     = ""   # Required. DNS name or ip.
  ports    = []   # Required. List of pors to monitor (ex.: [22, 80, 3028, 8080])
  timeout  = ""   # Optional. A time limit for requests (default = UpdateEvery)
```

#### sample configuration

```toml
[global]
  update_every = 5
  
  [job1.base]
    name = "localhost"
  [job1.specific]
    host = "127.0.0.1"
    ports = [80, 3128, 8080]
    timeout = "500ms"

  [job2.base]
    name = "remote"
  [job2.specific]
    host = "1.2.3.4"
    ports = [81, 3129, 8081]
    timeout = "2s"

```
