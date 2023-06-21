# phpDaemon collector

## Overview

[phpDaemon](https://github.com/kakserpom/phpdaemon) is an asynchronous server-side framework for Web and network
applications implemented in PHP using libevent.

This collector monitors metrics from one or more phpDaemon instances, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                  |         Dimensions         |  Unit   |
|-------------------------|:--------------------------:|:-------:|
| phpdaemon.workers       |      alive, shutdown       | workers |
| phpdaemon.alive_workers |   idle, busy, reloading    | workers |
| phpdaemon.idle_workers  | preinit, init, initialized | workers |
| phpdaemon.uptime        |            time            | seconds |

## Setup

### Prerequisites

#### Enable phpDaemon's HTTP server

Statistics expected to be in JSON format.

<details>
<summary>phpDaemon configuration</summary>

Instruction from [@METAJIJI](https://github.com/METAJIJI).

For enable `phpd` statistics on http, you must enable the http server and write an application.

Application is important, because standalone
application [ServerStatus.php](https://github.com/kakserpom/phpdaemon/blob/master/PHPDaemon/Applications/ServerStatus.php)
provides statistics in html format and unusable for `netdata`.

```php
// /opt/phpdaemon/conf/phpd.conf

path /opt/phpdaemon/conf/AppResolver.php;
Pool:HTTPServer {
    privileged;
    listen '127.0.0.1';
    port 8509;
}
```

```php
// /opt/phpdaemon/conf/AppResolver.php

<?php

class MyAppResolver extends \PHPDaemon\Core\AppResolver {
    public function getRequestRoute($req, $upstream) {
        if (preg_match('~^/(ServerStatus|FullStatus)/~', $req->attrs->server['DOCUMENT_URI'], $m)) {
            return $m[1];
        }
    }
}

return new MyAppResolver;
```

```php
/opt/phpdaemon/conf/PHPDaemon/Applications/FullStatus.php

<?php
namespace PHPDaemon\Applications;

class FullStatus extends \PHPDaemon\Core\AppInstance {
    public function beginRequest($req, $upstream) {
        return new FullStatusRequest($this, $upstream, $req);
    }
}
```

```php
// /opt/phpdaemon/conf/PHPDaemon/Applications/FullStatusRequest.php

<?php
namespace PHPDaemon\Applications;

use PHPDaemon\Core\Daemon;
use PHPDaemon\HTTPRequest\Generic;

class FullStatusRequest extends Generic {
    public function run() {
        $stime = microtime(true);
        $this->header('Content-Type: application/javascript; charset=utf-8');

        $stat = Daemon::getStateOfWorkers();
        $stat['uptime'] = time() - Daemon::$startTime;
        echo json_encode($stat);
    }
}
```

</details>

### Configuration

#### File

The configuration file name is `go.d/phpdaemon.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/phpdaemon.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|         Name         | Description                                                                                               |             Default              | Required |
|:--------------------:|-----------------------------------------------------------------------------------------------------------|:--------------------------------:|:--------:|
|     update_every     | Data collection frequency.                                                                                |                1                 |          |
| autodetection_retry  | Re-check interval in seconds. Zero means not to schedule re-check.                                        |                0                 |          |
|         url          | Server URL.                                                                                               | http://127.0.0.1:8509/FullStatus |   yes    |
|       timeout        | HTTP request timeout.                                                                                     |                2                 |          |
|       username       | Username for basic HTTP authentication.                                                                   |                                  |          |
|       password       | Password for basic HTTP authentication.                                                                   |                                  |          |
|      proxy_url       | Proxy URL.                                                                                                |                                  |          |
|    proxy_username    | Username for proxy basic HTTP authentication.                                                             |                                  |          |
|    proxy_password    | Password for proxy basic HTTP authentication.                                                             |                                  |          |
|        method        | HTTP request method.                                                                                      |               GET                |          |
|         body         | HTTP request body.                                                                                        |                                  |          |
|       headers        | HTTP request headers.                                                                                     |                                  |          |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects.                                  |                no                |          |
|   tls_skip_verify    | Server certificate chain and hostname validation policy. Controls whether the client performs this check. |                no                |          |
|        tls_ca        | Certification authority that the client uses when verifying the server's certificates.                    |                                  |          |
|       tls_cert       | Client TLS certificate.                                                                                   |                                  |          |
|       tls_key        | Client TLS key.                                                                                           |                                  |          |

</details>

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8509/FullStatus
```

</details>

##### HTTP authentication

HTTP authentication.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8509/FullStatus
    username: username
    password: password
```

</details>

##### HTTPS with self-signed certificate

HTTPS with self-signed certificate.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8509/FullStatus
    tls_skip_verify: yes
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collecting metrics from local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8509/FullStatus

  - name: remote
    url: http://192.0.2.1:8509/FullStatus
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `phpdaemon` collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m phpdaemon
  ```
