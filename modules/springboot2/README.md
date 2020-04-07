# Any Spring Boot 2 application monitoring with Netdata

This module monitors one or more Java Spring-boot 2 applications depending on configuration.
Netdata can be used to monitor running Java [Spring Boot 2](https://spring.io/) applications that expose their metrics with the use of the **Spring Boot Actuator** included in Spring Boot library.

Springboot2 module looks up `http://localhost:8080/actuator/prometheus` and `http://127.0.0.1:8080/actuator/prometheus` to detect Spring Boot application by default.

## Charts

-   Response Codes in `requests/s`
-   Threads in `threads`
-   Heap Memory Usage Overview in `bytes`
-   Heap Memory Usage Eden Space in `bytes`
-   Heap Memory Usage Survivor Space in `bytes`
-   Heap Memory Usage Old Space in `bytes`
-   Uptime in `seconds`

## Configuration

Edit the `go.d/springboot2.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/springboot2.conf
```

The Spring Boot Actuator exposes these metrics over HTTP and is very easy to use:

-   add `org.springframework.boot:spring-boot-starter-actuator` and `io.micrometer:micrometer-registry-prometheus` to your application dependencies
-   set `management.endpoints.web.exposure.include=*` in your `application.properties`

Please refer to the [Spring Boot Actuator: Production-ready features](https://docs.spring.io/spring-boot/docs/current/reference/html/production-ready.html) 
and [81. Actuator - Part IX. ‘How-to’ guides](https://docs.spring.io/spring-boot/docs/current/reference/html/howto-actuator.html) for more information.

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://localhost:8080/actuator/prometheus

  - name: remote
    url: http://203.0.113.10:8080/actuator/prometheus
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/springboot2.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m springboot2
