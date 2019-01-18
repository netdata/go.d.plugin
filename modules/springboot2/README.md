# springboot2

This module will monitor one or more Java Spring-boot 2 applications depending on configuration.
Netdata can be used to monitor running Java [Spring Boot 2](https://spring.io/) applications 
that expose their metrics with the use of the **Spring Boot Actuator** 
included in Spring Boot library.

## Configuration

The Spring Boot Actuator exposes these metrics over HTTP and is very easy to use:
* add `org.springframework.boot:spring-boot-starter-actuator` and `io.micrometer:micrometer-registry-prometheus` to your application dependencies
* set `management.endpoints.web.exposure.include=*` in your `application.properties`

Please refer [Spring Boot Actuator: Production-ready features](https://docs.spring.io/spring-boot/docs/current/reference/html/production-ready.html) 
and [81. Actuator - Part IX. ‘How-to’ guides](https://docs.spring.io/spring-boot/docs/current/reference/html/howto-actuator.html) 
for more information.

## Charts

1. **Response Codes** in requests/s
 * 1xx
 * 2xx
 * 3xx
 * 4xx
 * 5xx

2. **Threads**
 * daemon
 * total

3. **GC Time** in milliseconds and **GC Operations** in operations/s
 * Copy
 * MarkSweep
 * ...

4. **Heap Mmeory Usage** in KB
 * used
 * committed

## Usage

The springboot module is enabled by default. It looks up `http://localhost:8080/actuator/prometheus` 
and `http://127.0.0.1:8080/actuator/prometheus` to detect Spring Boot application by default. 
You can change it by editing `/etc/netdata/go.d/springboot2.conf` 
(to edit it on your system run `/etc/netdata/edit-config go.d/springboot2.conf`).

Please check [springboot2.conf](../../config/go.d/springboot2.conf) for more examples.
