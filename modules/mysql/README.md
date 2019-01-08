# mysql

Module monitor MySQL DB performance and health metrics.

Following charts are drawn:

_TODO:_

### configuration

```yaml
jobs:
  - name: local
    host: ::1
    port: 3306
    user: user
    pass: pass
    # socket: /var/run/mysqld/mysql.sock
```

```yaml
jobs:
  - name: local
    dsn: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
    # username:password@protocol(address)/dbname?param=value
    # user:password@/dbname
```
---
