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

5. **System Metrics from gv$sysmetric**
 * buffer_cachehit_ratio
 * cursor_cachehit_ratio
 * library_cachehit_ratio
 * shared_pool_free
 * physical_reads
 * physical_writes
 * enqueue_timeouts
 * gc_cr_block_received
 * cache_blocks_corrupt
 * cache_blocks_lost
 * logons
 * active_sessions
 * long_table_scans
 * service_response_time
 * user_rollbacks
 * sorts_per_user_call
 * rows_per_sort
 * disk_sorts
 * memory_sorts_ratio
 * database_wait_time_ratio
 * session_limit_usage
 * session_count
 * temp_space_used

### configuration

```yaml
jobs:
  - name: local
    dsn: SYSTEM/Oracle12345@ORCL

```
---
