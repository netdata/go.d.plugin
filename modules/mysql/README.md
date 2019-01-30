# mysql

Module mysql monitors one or more MySQL servers.

It will produce following charts (if data is available):

1. **Bandwidth** in kilobits/s
 * in
 * out

2. **Queries** in queries/sec
 * queries
 * questions
 * slow queries

3. **Queries By Type** in queries/s
 * select
 * delete
 * update
 * insert
 * cache hits
 * replace

4. **Handlerse** in handlers/s
 * commit
 * delete
 * prepare
 * read first
 * read key
 * read next
 * read prev
 * read rnd
 * read rnd next
 * rollback
 * savepoint
 * savepoint rollback
 * update
 * write

4. **Table Locks** in locks/s
 * immediate
 * waited

5. **Table Select Join Issuess** in joins/s
 * full join
 * full range join
 * range
 * range check
 * scan

6. **Table Sort Issuess** in joins/s
 * merge passes
 * range
 * scan

7. **Tmp Operations** in created/s
 * disk tables
 * files
 * tables

8. **Connections** in connections/s
 * all
 * aborted

9. **Connections Active** in connections/s
 * active
 * limit
 * max active

10. **Binlog Cache** in threads
 * disk
 * all

11. **Threads** in transactions/s
 * connected
 * cached
 * running

12. **Threads Creation Rate** in threads/s
 * created

13. **Threads Cache Misses** in misses
 * misses

14. **InnoDB I/O Bandwidth** in KiB/s
 * read
 * write

15. **InnoDB I/O Operations** in operations/s
 * reads
 * writes
 * fsyncs

16. **InnoDB Pending I/O Operations** in operations/s
 * reads
 * writes
 * fsyncs

17. **InnoDB Log Operations** in operations/s
 * waits
 * write requests
 * writes

18. **InnoDB OS Log Pending Operations** in operations
 * fsyncs
 * writes

19. **InnoDB OS Log Operations** in operations/s
 * fsyncs

20. **InnoDB OS Log Bandwidth** in KiB/s
 * write

21. **InnoDB Current Row Locks** in operations
 * current waits

22. **InnoDB Row Operations** in operations/s
 * inserted
 * read
 * updated
 * deleted

23. **InnoDB Buffer Pool Pagess** in pages
 * data
 * dirty
 * free
 * misc
 * total

24. **InnoDB Buffer Pool Flush Pages Requests** in requests/s
 * flush pages

25. **InnoDB Buffer Pool Bytes** in MiB
 * data
 * dirty

26. **InnoDB Buffer Pool Operations** in operations/s
 * disk reads
 * wait free

27. **QCache Operations** in queries/s
 * hits
 * lowmem prunes
 * inserts
 * no caches

28. **QCache Queries in Cache** in queries
 * queries

29. **QCache Free Memory** in MiB
 * free

30. **QCache Memory Blocks** in blocks
 * free
 * total

31. **MyISAM Key Cache Blocks** in blocks
 * unused
 * used
 * not flushed

32. **MyISAM Key Cache Requests** in requests/s
 * reads
 * writes

33. **MyISAM Key Cache Requests** in requests/s
 * reads
 * writes

34. **MyISAM Key Cache Disk Operations** in operations/s
 * reads
 * writes

35. **Open Files** in files
 * files

36. **Opened Files Rate** in files/s
 * files

37. **Binlog Statement Cache** in statements/s
 * disk
 * all

38. **Connection Errors** in errors/s
 * accept
 * internal
 * max
 * peer addr
 * select
 * tcpwrap

39. **Slave Behind Seconds** in seconds
 * time

40. **I/O / SQL Thread Running State** in bool
 * sql
 * io

41. **Replicated Writesets** in writesets/s
 * rx
 * tx

42. **Replicated Bytes** in KiB/s
 * rx
 * tx

43. **Galera Queue** in writesets
 * rx
 * tx

44. **Replication Conflicts** in transactions
 * bf aborts
 * cert fails

45. **Flow Control** in ms
 * paused


### configuration
[DSN syntax in details](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

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
