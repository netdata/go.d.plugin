# solr

Module monitors solr request handler statistics.

Following charts are drawn (per core):

1. **Search Requests** in requests/s
 * search
 
2. **Search Errors** in errors/s
 * errors
 
3. **Search Errors By Type** in errors/s
 * client
 * server
 * timeout
 
4. **Search Requests Processing Time** in milliseconds
 * time

5. **Search Requests Timings** in milliseconds
 * min
 * median
 * mean
 * max
 
6. **Search Requests Processing Time Percentile** in milliseconds
 * p75
 * p95
 * p99
 * p999
 
7. **Update Requests** in requests/s
 * update
 
8. **Update Errors** in errors/s
 * errors
 
9. **Update Errors By Type** in errors/s
 * client
 * server
 * timeout
 
10. **Update Requests Processing Time** in milliseconds
 * time

11. **Update Requests Timings** in milliseconds
 * min
 * median
 * mean
 * max
 
12. **Update Requests Processing Time Percentile** in milliseconds
 * p75
 * p95
 * p99
 * p999

### configuration

```yaml
jobs:
  - name: local
    url : http://localhost:8983
      
  - name: remote
    url : http://100.64.0.1:8983

```

When no configuration file is found, module tries to connect to: `localhost:8983`.

---