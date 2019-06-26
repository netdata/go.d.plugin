# phpdaemon

This module collects [`phpdaemon`](https://github.com/kakserpom/phpdaemon) workers statistics via http.

**Requirements:**
 * `phpdaemon` with enabled [`http`] server.
 * statistics should be reported in `json` format.

It produces the following charts:

1. **Workers** in workers
 * alive
 * shutdown

2. **Alive Workers State** in workers
 * idle
 * busy
 * reloading

3. **Idle Workers State** in workers
 * preinit
 * init
 * initialized
 
4. **Uptime** in seconds
 * time
 * starting
 * reading
 * sending
 * keepalive
 * dns lookup
 * closing
 * logging
 * finishing
 * idle cleanup
 * open
 
### configuration

Here is an example for 2 instances:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:8509/FullStatus
      
  - name: remote
    url : http://10.0.0.1:8509/FullStatus
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/phpdaemon.conf).

Without configuration, module attempts to connect to `http://127.0.0.1:8509/FullStatus`

### phpdaemon configuration

Instruction from [@METAJIJI](https://github.com/METAJIJI)

For enable `phpd` statistics on http, you must enable the http server and write an application.

Application is important, because standalone application [ServerStatus.php](https://github.com/kakserpom/phpdaemon/blob/master/PHPDaemon/Applications/ServerStatus.php) provide statistics in html format and unusable for `netdata`. 

> /opt/phpdaemon/conf/phpd.conf

```php
path /opt/phpdaemon/conf/AppResolver.php;
Pool:HTTPServer {
    privileged;
    listen '127.0.0.1';
    port 8509;
}
```

> /opt/phpdaemon/conf/AppResolver.php

```php
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

> /opt/phpdaemon/conf/PHPDaemon/Applications/FullStatus.php

```php
<?php
namespace PHPDaemon\Applications;

class FullStatus extends \PHPDaemon\Core\AppInstance {
    public function beginRequest($req, $upstream) {
        return new FullStatusRequest($this, $upstream, $req);
    }
}
```

> /opt/phpdaemon/conf/PHPDaemon/Applications/FullStatusRequest.php

```php
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

---
