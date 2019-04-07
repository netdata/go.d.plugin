# coredns

This module will monitor one or more coredns instances.


It produces the following charts:

 - Summary:

1. **Number Of DNS Requests** in requests/s

2. **Number Of DNS Responses** in responses/s

3. **Number Of Processed And Dropped DNS Requests** in requests/s

4. **Number Of Dropped DNS Requests Because Of No Matching Zone** in requests/s

5. **Number Of Panics** in panics/s

6. **Number Of DNS Requests Per Transport Protocol** in requests/s
 * TCP, UDP

7. **Number Of DNS Requests Per IP Family** in requests/s
 * IPv4, IPv6

8. **Number Of DNS Requests Per Type** in requests/s
 * A, AAAA, MX, SOA, CNAME, PTR, TXT, NS, DS, DNSKEY, RRSIG, NSEC, NSEC3, IXFR, ANY, other
 
9. **Number Of DNS Responses Per Rcode** in responses/s
 * NOERROR, FORMERR, SERVFAIL, NXDOMAIN, NOTIMP, REFUSED, YXDOMAIN, YXRESET, NSRESET
 * NOTAUTH, NOTZONE, BADSIG, BADKEY, BADTIME, BADMODE, BADNAME, BADLAG,
 * BADTRUNC, BADCOOKIE, other 
 
 - Per Server (if configured):

1. **Number Of DNS Requests** in requests/s

2. **Number Of DNS Responses** in responses/s

3. **Number Of Processed And Dropped DNS Requests** in requests/s

4. **Number Of DNS Requests Per Transport Protocol** in requests/s
 * TCP, UDP

5. **Number Of DNS Requests Per IP Family** in requests/s
 * IPv4, IPv6

6. **Number Of DNS Requests Per Type** in requests/s
 * A, AAAA, MX, SOA, CNAME, PTR, TXT, NS, DS, DNSKEY, RRSIG, NSEC, NSEC3, IXFR, ANY, other
 
7. **Number Of DNS Responses Per Rcode** in responses/s
 * NOERROR, FORMERR, SERVFAIL, NXDOMAIN, NOTIMP, REFUSED, YXDOMAIN, YXRESET, NSRESET
 * NOTAUTH, NOTZONE, BADSIG, BADKEY, BADTIME, BADMODE, BADNAME, BADLAG,
 * BADTRUNC, BADCOOKIE, other
 
 - Per Zone (if configured):

1. **Number Of DNS Requests** in requests/s

2. **Number Of DNS Responses** in responses/s

3. **Number Of DNS Requests Per Transport Protocol** in requests/s
 * TCP, UDP

4. **Number Of DNS Requests Per IP Family** in requests/s
 * IPv4, IPv6

5. **Number Of DNS Requests Per Type** in requests/s
 * A, AAAA, MX, SOA, CNAME, PTR, TXT, NS, DS, DNSKEY, RRSIG, NSEC, NSEC3, IXFR, ANY, other
 
6. **Number Of DNS Responses Per Rcode** in responses/s
 * NOERROR, FORMERR, SERVFAIL, NXDOMAIN, NOTIMP, REFUSED, YXDOMAIN, YXRESET, NSRESET
 * NOTAUTH, NOTZONE, BADSIG, BADKEY, BADTIME, BADMODE, BADNAME, BADLAG,
 * BADTRUNC, BADCOOKIE, other


### configuration

Needs only `url` to coredns metric-address.

Here is an example:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1:9153/metrics
      
  - name: remote
    url : http://100.64.0.1:9153/metrics
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/coredns.conf).

Without configuration, module attempts to connect to `http://127.0.0.1:9153/metrics`

---
