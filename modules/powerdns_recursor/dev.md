`PowerDNS Recursor` docker setup:

- https://github.com/pschiffe/docker-pdns#pdns-recursor

Run `PowerDNS Recursor` with enabled webserver and api:

```cmd
docker run -d -p 8081:8081 --name pdns-recursor \
  -e PDNS_api_key=secret \
  -e PDNS_webserver=yes \
  -e PDNS_webserver-port=8081 \
  -e PDNS_webserver_address=0.0.0.0 \
  -e PDNS_webserver-allow-from=0.0.0.0/0 \
  pschiffe/pdns-recursor
```

Gather metrics:

```cmd
curl http://127.0.0.1:8081/metrics
curl http://127.0.0.1:8081/api/v1/servers/localhost/statistics
```

Useful links:

- [Recursor documentation](https://doc.powerdns.com/recursor/).
- [Webserver/API/URL Endpoints docs](https://doc.powerdns.com/recursor/http-api/index.html).
- [Statistics endpoint metric description](https://doc.powerdns.com/recursor/metrics.html#metricnames).
