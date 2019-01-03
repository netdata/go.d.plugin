# rabbitmq

Module monitor Oracle DB performance and health metrics.

Following charts are drawn:

1. **Processes**
  * total

2. **Sessions**
  * total
  * active
  * inactive

2. **Activity**
 * parse count (total)
 * execute count
 * user commits
 * user rollbacks

3. **Wait Times**
 * configuration
 * administrative
 * system I/O
 * application
 * concurrency
 * commit
 * network
 * user I/O
 * other

4. **Tablespace Size**
 * max_bytes_system
 * free_bytes_system
 * bytes_system
 * max_bytes_sysaux
 * free_bytes_sysaux
 * bytes_sysaux
 * max_bytes_users
 * free_bytes_users
 * bytes_users
 * max_bytes_temp
 * free_bytes_temp
 * bytes_temp

### configuration

```yaml
jobs:
  - name: local
    dsn: SYSTEM/Oracle12345@ORCL

```
---
