# mysql

Module mysql monitors one or more MySQL servers.

It will produce following charts (if data is available):

1. **Bandwidth** in kbps
 * in
 * out

2. **Queries** in queries/sec
 * queries
 * questions
 * slow queries

3. **Operations** in operations/sec
 * opened tables
 * flush
 * commit
 * delete
 * prepare
 * read first
 * read key
 * read next
 * read prev
 * read random
 * read random next
 * rollback
 * save point
 * update
 * write

4. **Table Locks** in locks/sec
 * immediate
 * waited

5. **Select Issues** in issues/sec
 * full join
 * full range join
 * range
 * range check
 * scan

6. **Sort Issues** in issues/sec
 * merge passes
 * range
 * scan

### configuration

```yaml
jobs:
  - name: local
    dsn: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
    # username:password@protocol(address)/dbname?param=value
    # user:password@/dbname
    # Examples:
    # - name: local
    #   dsn: user:pass@unix(/usr/local/var/mysql/mysql.sock)/
    # - name: remote
    #   dsn: user:pass5@localhost/mydb?charset=utf8
```

If no configuration is given, module will attempt to connect to mysql server via unix socket at:
1. `/var/run/mysqld/mysqld.sock` without password and with username `root`;
2. `/usr/local/var/mysql/mysql.sock` without password and with username `root`;
3. `localhost:3306` without password and with username `root`.
---
