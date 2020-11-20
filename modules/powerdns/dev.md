`PowerDNS Authoritative Server` docker setup:

-   https://github.com/pschiffe/docker-pdns#pdns-mysql

Run `MariaDB` backend:

```cmd
docker run -d \
  --name pqdn-mariadb-backend \
  -e MYSQL_ROOT_PASSWORD=supersecret \
  mariadb:10.1
```

Run `PowerDNS Authoritative Server` with enabled webserver and api:

```cmd
docker run -d -p 8081:8081 --name pdns-master \
  --hostname ns1.example.com --link pqdn-mariadb-backend:mysql \
  -e PDNS_master=yes \
  -e PDNS_api=yes \
  -e PDNS_api_key=secret \
  -e PDNS_webserver=yes \
  -e PDNS_webserver-port=8081 \
  -e PDNS_webserver_address=0.0.0.0 \
  -e PDNS_webserver-allow-from=0.0.0.0/0 \
  -e PDNS_gmysql_password=supersecret \
  pschiffe/pdns-mysql
```

Gather metrics (seems doesn't support unauthenticated requests to the API):

```cmd
curl -H 'X-Api-Key: secret' http://127.0.0.1:8082/api/v1/servers/localhost/statistics
```

Useful links:
-   [Authoritative Server documentation](https://doc.powerdns.com/authoritative/).
-   [Webserver/API/URL Endpoints docs](https://doc.powerdns.com/authoritative/http-api/index.html).
-   [Statistics endpoint metric description](https://doc.powerdns.com/authoritative/http-api/statistics.html).
